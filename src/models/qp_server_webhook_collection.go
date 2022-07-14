package models

import (
	"strings"

	log "github.com/sirupsen/logrus"
)

// Webhook model
type QpServerWebhookCollection struct {
	Webhooks []*QpWebhook `json:"webhooks,omitempty"`
	context  string
	db       QpDataWebhookInterface
}

// Fill start memory cache
func (source *QpServerWebhookCollection) WebhookFill(context string, db QpDataWebhookInterface) (err error) {
	source.Webhooks = []*QpWebhook{}
	source.context = context
	source.db = db

	whooks, err := source.db.FindAll(source.context)
	if err != nil {
		log.Errorf("erro", err)
		return
	}

	for _, element := range whooks {
		source.Webhooks = append(source.Webhooks, element.QpWebhook)
	}

	return
}

func (source *QpServerWebhookCollection) WebhookAdd(url string, forwardinternal bool, trackid string) (affected uint, err error) {
	botWHook, err := source.db.Find(source.context, url)
	if err != nil {
		return
	}

	var wHook *QpWebhook
	if botWHook != nil {
		botWHook.ForwardInternal = forwardinternal
		botWHook.TrackId = trackid
		err = source.db.Update(*botWHook)
		if err != nil {
			return
		}
		wHook = botWHook.QpWebhook
	} else {
		err = source.db.Add(source.context, url, forwardinternal, trackid)
		if err == nil {
			wHook = &QpWebhook{
				Url:             url,
				ForwardInternal: forwardinternal,
				TrackId:         trackid,
			}
		}
	}

	if wHook != nil {
		exists := false
		for index, element := range source.Webhooks {
			if element.Url == url {
				source.Webhooks = append(source.Webhooks[:index], source.Webhooks[index+1:]...) // remove
				source.Webhooks = append(source.Webhooks, wHook)                                // append a clean one
				exists = true
				affected++
			}
		}

		if !exists {
			source.Webhooks = append(source.Webhooks, wHook)
			affected++
		}
	}

	return
}

func (source *QpServerWebhookCollection) WebhookRemove(url string) (affected uint, err error) {
	i := 0 // output index
	for _, element := range source.Webhooks {
		if strings.Contains(element.Url, url) {
			err = source.db.Remove(source.context, element.Url)
			if err == nil {
				affected++
			} else {
				source.Webhooks[i] = element
				i++
				break
			}
		} else {
			source.Webhooks[i] = element
			i++
		}
	}

	for j := i; j < len(source.Webhooks); j++ {
		source.Webhooks[j] = nil
	}
	source.Webhooks = source.Webhooks[:i]
	return
}

func (source *QpServerWebhookCollection) WebhookClear() (err error) {
	return source.db.Clear(source.context)
}

/*
func (source *QpServerWebhookCollection) WebhookFailure(url string) {
	log.Infof("failure on webhook from: %s", url)
	for index, element := range source.Webhooks {
		if element.Url == url {
			//source.Webhooks[index].Failure = &time.Time{}
		}
	}
}
*/
