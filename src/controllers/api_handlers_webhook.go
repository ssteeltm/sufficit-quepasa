package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"

	. "github.com/sufficit/sufficit-quepasa-fork/models"
)

//region CONTROLLER - WEBHOOK

func WebhookController(w http.ResponseWriter, r *http.Request) {

	// setting default reponse type as json
	w.Header().Set("Content-Type", "application/json")

	token := chi.URLParam(r, "token")
	server, err := GetServerFromToken(token)
	if err != nil {
		RespondNotFound(w, fmt.Errorf("Token '%s' not found on WebHookHandler", token))
		return
	}

	// Declare a new Person struct.
	var p QPWebhookRequest

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err = json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		RespondServerError(server, w, err)
	}

	// Já tratei os parametros
	if ENV.IsDevelopment() {
		log.Printf("(%s) Updating Webhook: %s", server.Bot.GetNumber(), p.Url)
	}
	switch os := r.Method; os {
	case http.MethodPost:
		err = server.WebhookAdd(p.Url)
		if err != nil {
			RespondServerError(server, w, err)
			return
		}
	case http.MethodDelete:
		err = server.WebhookRemove(p.Url)
		if err != nil {
			RespondServerError(server, w, err)
			return
		}
	}

	RespondSuccess(w, server)
}

//endregion
