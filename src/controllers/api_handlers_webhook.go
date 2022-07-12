package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi"

	models "github.com/sufficit/sufficit-quepasa-fork/models"
)

//region CONTROLLER - WEBHOOK

func WebhookController(w http.ResponseWriter, r *http.Request) {

	// setting default reponse type as json
	w.Header().Set("Content-Type", "application/json")

	response := &models.QpWebhookResponse{}

	token := chi.URLParam(r, "token")
	server, err := models.GetServerFromToken(token)
	if err != nil {
		response.ParseError(err)
		RespondNotFound(w, response)
		return
	}

	// Declare a new Person struct.
	var p models.QPWebhookRequest

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	_ = json.NewDecoder(r.Body).Decode(&p)

	switch os := r.Method; os {
	case http.MethodPost:
		affected, err := server.WebhookAdd(p.Url)
		if err != nil {
			response.ParseError(err)
			RespondServerError(server, w, response)
		} else {
			response.Affected = affected
			response.ParseSuccess("updated with success")
			RespondSuccess(w, response)
			if affected > 0 {
				server.Log.Infof("updating webhook url: %s, items affected: %v", p.Url, affected)
			}
		}
		return
	case http.MethodDelete:
		affected, err := server.WebhookRemove(p.Url)
		if err != nil {
			response.ParseError(err)
			RespondServerError(server, w, response)
		} else {
			response.Affected = affected
			response.ParseSuccess("deleted with success")
			RespondSuccess(w, response)
			if affected > 0 {
				server.Log.Infof("removing webhook url: %s, items affected: %v", p.Url, affected)
			}
		}
		return
	default:
		response.Webhooks = filterByUrl(server.Webhooks, p.Url)
		if len(p.Url) > 0 {
			response.ParseSuccess(fmt.Sprintf("getting with filter: %s", p.Url))
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
