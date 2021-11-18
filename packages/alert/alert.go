package alert

import (
	"errors"

	"github.com/xpartacvs/go-dishook"
)

type Alert struct {
	payload dishook.Payload
}

func (a *Alert) Send(url string, flushAfter bool) error {
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

func (a *Alert) FlushFields() *Alert {
	for i := range a.payload.Embeds {
		a.payload.Embeds[i].Fields = nil
	}
	return a
}

func (a *Alert) AddField(title, content string, inline bool) *Alert {
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

func (a *Alert) SetBotAvatar(url string) *Alert {
	a.payload.AvatarUrl = dishook.Url(url)
	return a
}

func (a *Alert) SetBotName(name string) *Alert {
	a.payload.Username = name
	return a
}

func New(message string) *Alert {
	return &Alert{
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
	}
}
