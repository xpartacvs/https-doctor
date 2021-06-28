package worker

import (
	"fmt"
	"https-doctor/packages/alert"
	"https-doctor/packages/client"
	"https-doctor/packages/config"
	"https-doctor/packages/logger"
	"strings"
	"sync"
	"time"
)

func Start() error {

	return execute()
}

func execute() error {
	var wg sync.WaitGroup
	notif := alert.New(config.Get().DishookBotMessage())
	notif.SetBotName(config.Get().DishookBotName()).SetBotAvatar(config.Get().DishookBotAvatar())

	for _, h := range config.Get().Hosts() {
		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			var title, content string = strings.ToUpper(host), ""
			httpsClient := client.New(host, logger.Log())
			result, expiry := httpsClient.GetExpiry()
			switch result {
			case client.ErrTimeout:
				content = "Connection timeout."
			case client.ErrConnection:
				content = "Connection error."
			case client.ErrCertInvalid:
				content = "Wrong SSL certificate."
			case client.ErrCertExpired:
				days := int(time.Until(expiry).Hours()/24) * -1
				content = fmt.Sprintf("Expired for %d days.", days)
			default:
				days := int(time.Until(expiry).Hours() / 24)
				content = fmt.Sprintf("Expired in %d days.", days)
			}
			notif.AddField(title, content, true)
		}(h)
	}

	wg.Wait()
	if err := notif.Send(config.Get().DishookURL(), true); err != nil {
		logger.Log().Err(err)
		return err
	}

	return nil
}
