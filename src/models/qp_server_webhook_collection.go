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

func (source *QpServerWebhookCollection) WebhookAdd(webhook *QpWebhook) (affected uint, err error) {
	botWHook, err := source.db.Find(source.context, webhook.Url)
	if err != nil {
		return
	}

	if botWHook != nil {
		botWHook.ForwardInternal = webhook.ForwardInternal
		botWHook.TrackId = webhook.TrackId
		botWHook.Extra = webhook.Extra
		err = source.db.Update(*botWHook)
		if err != nil {
			return
		}
	} else {
		dbWebhook := &QpBotWebhook{
			Context:   source.context,
			QpWebhook: webhook,
		}
		err = source.db.Add(*dbWebhook)
		if err != nil {
			return
		}
	}

	exists := false
	for index, element := range source.Webhooks {
		if element.Url == webhook.Url {
			source.Webhooks = append(source.Webhooks[:index], source.Webhooks[index+1:]...) // remove
			source.Webhooks = append(source.Webhooks, webhook)                              // append a clean one
			exists = true
			affected++
		}
	}

	if !exists {
		source.Webhooks = append(source.Webhooks, webhook)
		affected++
	}

	return
}

func (source *QpServerWebhookCollection) WebhookRemove(url string) (affected uint, err error) {
	i := 0 // output index
	for _, element := range source.Webhooks {
		if len(url) == 0 || strings.Contains(element.Url, url) {
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
