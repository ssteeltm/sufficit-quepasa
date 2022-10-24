package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	models "github.com/sufficit/sufficit-quepasa/models"
)

//region CONTROLLER - WEBHOOK

func WebhookController(w http.ResponseWriter, r *http.Request) {

	// setting default reponse type as json
	w.Header().Set("Content-Type", "application/json")

	response := &models.QpWebhookResponse{}

	server, err := GetServer(r)
	if err != nil {
		response.ParseError(err)
		RespondInterface(w, response)
		return
	}

	// Declare a new Person struct.
	var webhook *models.QpWebhook

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err = json.NewDecoder(r.Body).Decode(&webhook)
	if err != nil {
		response.ParseError(err)
		RespondInterface(w, response)
		return
	}

	switch os := r.Method; os {
	case http.MethodPost:
		affected, err := server.WebhookAdd(webhook)
		if err != nil {
			response.ParseError(err)
			RespondInterface(w, response)
		} else {
			response.Affected = affected
			response.ParseSuccess("updated with success")
			RespondSuccess(w, response)
			if affected > 0 {
				server.Log.Infof("updating webhook url: %s, items affected: %v", webhook.Url, affected)
			}
		}
		return
	case http.MethodDelete:
		affected, err := server.WebhookRemove(webhook.Url)
		if err != nil {
			response.ParseError(err)
			RespondInterface(w, response)
		} else {
			response.Affected = affected
			response.ParseSuccess("deleted with success")
			RespondSuccess(w, response)
			if affected > 0 {
				server.Log.Infof("removing webhook url: %s, items affected: %v", webhook.Url, affected)
			}
		}
		return
	default:
		url := r.Header.Get("X-QUEPASA-WHURL")
		response.Webhooks = filterByUrl(server.Webhooks, url)
		if len(url) > 0 {
			response.ParseSuccess(fmt.Sprintf("getting with filter: %s", url))
		} else {
			response.ParseSuccess("getting without filter")
		}

		RespondSuccess(w, response)
		return
	}
}

func filterByUrl(source []*models.QpWebhook, filter string) (out []models.QpWebhook) {
	for _, element := range source {
		if strings.Contains(element.Url, filter) {
			out = append(out, *element)
		}
	}
	return
}

//endregion
