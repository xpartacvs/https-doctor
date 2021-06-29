package alert

import (
	"errors"

	"github.com/rs/zerolog"
	"github.com/xpartacvs/go-dishook"
)

type alert struct {
	payload dishook.Payload
	logger  *zerolog.Logger
}

type Alert interface {
	AddField(title, content string, inline bool) Alert
	FlushFields() Alert
	Send(url string, flushAfter bool) error
	SetBotName(name string) Alert
	SetBotAvatar(url string) Alert
	SetZerolog(logger *zerolog.Logger) Alert
}

func (a *alert) SetZerolog(logger *zerolog.Logger) Alert {
	a.logger = logger
	return a
}

func (a *alert) Send(url string, flushAfter bool) error {
	for _, e := range a.payload.Embeds {
		if len(e.Fields) <= 0 {
			return errors.New("alert has nothing to send")
		}
	}
	_, err := dishook.Send(url, a.payload)
	if flushAfter {
		a.FlushFields()
	}
	return err
}

func (a *alert) FlushFields() Alert {
	for i := range a.payload.Embeds {
		a.payload.Embeds[i].Fields = nil
	}
	return a
}

func (a *alert) AddField(title, content string, inline bool) Alert {
	f := dishook.Field{
		Name:   title,
		Value:  content,
		Inline: inline,
	}
	for i := range a.payload.Embeds {
		a.payload.Embeds[i].Fields = append(a.payload.Embeds[i].Fields, f)
	}
	return a
}

func (a *alert) SetBotAvatar(url string) Alert {
	a.payload.AvatarUrl = dishook.Url(url)
	return a
}

func (a *alert) SetBotName(name string) Alert {
	a.payload.Username = name
	return a
}

func New(message string) Alert {
	return &alert{
		payload: dishook.Payload{
			Content: message,
			Embeds: []dishook.Embed{
				{
					Color:       dishook.ColorWarn,
					Title:       "Some Hostnames Need Your Attention",
					Description: "Please consider the following",
					Fields:      nil,
				},
			},
		},
		logger: nil,
	}
}
