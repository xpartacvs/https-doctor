package worker

import (
	"https-doctor/packages/alert"
	"https-doctor/packages/client"
	"https-doctor/packages/config"
	"https-doctor/packages/logger"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog"
)

func Start() error {
	scheduler := gocron.NewScheduler(config.Get().Location())
	if _, err := scheduler.Cron(config.Get().Schedule()).Do(job); err != nil {
		return err
	}
	scheduler.StartBlocking()
	return nil
}

func job() {
	if len(config.Get().Hosts()) <= 0 {
		logger.Log().Warn().Msg("Worker: No host to check on")
		return
	}

	var wg sync.WaitGroup
	notif := alert.New(config.Get().DishookBotMessage())
	notif.SetBotName(config.Get().DishookBotName()).SetBotAvatar(config.Get().DishookBotAvatar())

	for _, h := range config.Get().Hosts() {
		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			var title, content string = strings.ToUpper(host), ""
			httpsClient := client.New(host, logger.Log())
			expiry, err := httpsClient.GetExpiry()
			switch err {
			case client.ErrTimeout:
				content = "Connection timeout."
				if config.Get().ZerologLevel() != zerolog.Disabled {
					notif.AddField(title, content, false)
				}
			case client.ErrConnection:
				content = "Connection error."
				if config.Get().ZerologLevel() != zerolog.Disabled {
					notif.AddField(title, content, false)
				}
			case client.ErrCertInvalid:
				content = "Wrong SSL certificate."
				notif.AddField(title, content, false)
			case client.ErrCertExpired:
				days := int(time.Until(*expiry).Hours()/24) * -1
				content = "SSL expired for " + strconv.Itoa(days) + " days."
				notif.AddField(title, content, false)
			default:
				graceTime := expiry.AddDate(0, 0, config.Get().Graceperiod())
				if time.Now().After(graceTime) {
					days := int(time.Until(*expiry).Hours() / 24)
					content = "SSL expired in " + strconv.Itoa(days) + " days."
					notif.AddField(title, content, false)
				}
			}
		}(h)
	}

	wg.Wait()
	if err := notif.Send(config.Get().DishookURL(), true); err != nil {
		logger.Log().Err(err)
	}
}
