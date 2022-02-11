package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"

	. "github.com/sufficit/sufficit-quepasa-fork/models"
)

const APIVersion1 string = "v1"

var ControllerPrefixV1 string = fmt.Sprintf("/%s/bot/{token}", APIVersion1)

func RegisterAPIV1Controllers(r chi.Router) {
	r.Get(ControllerPrefixV1, InformationControllerV1)
	r.Post(ControllerPrefixV1+"/send", SendAPIHandlerV1)
	r.Get(ControllerPrefixV1+"/receive", ReceiveAPIHandlerV1)
	r.Post(ControllerPrefixV1+"/attachment", AttachmentAPIHandlerV2)
	r.Post(ControllerPrefixV1+"/webhook", WebhookControllerV1)
}

//region CONTROLLER - INFORMATION

// Renders route GET "/{version}/bot/{token}"
func InformationControllerV1(w http.ResponseWriter, r *http.Request) {

	// setting default reponse type as json
	w.Header().Set("Content-Type", "application/json")

	token := chi.URLParam(r, "token")
	server, err := GetServerFromToken(token)
	if err != nil {
		RespondNotFound(w, fmt.Errorf("Token '%s' not found", token))
		return
	}

	wid := server.GetWid()

	var ep QPEndPoint
	if !strings.Contains(wid, "@") {
		ep.ID = wid + "@s.whatsapp.net"
	} else {
		ep.ID = wid
	}

	ep.Phone = server.Bot.GetNumber()
	if server.Bot.Verified {
		ep.Status = "verified"
	} else {
		ep.Status = "unverified"
	}

	RespondSuccess(w, ep)
}

//endregion
//region CONTROLLER - WEBHOOK

func WebhookControllerV1(w http.ResponseWriter, r *http.Request) {

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

	server.Bot.WebHook = p.Url
	// Atualizando banco de dados
	if err := server.Bot.WebHookUpdate(); err != nil {
		return
	}

	RespondSuccess(w, ToQPBotV1(server.Bot))
}

//endregion
