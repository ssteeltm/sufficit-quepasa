package controllers

import (
	"fmt"
	//"encoding/json"
	"net/http"

	//log "github.com/sirupsen/logrus"
	"github.com/go-chi/chi"

	. "github.com/sufficit/sufficit-quepasa-fork/models"
)

const APIVersion1 string = "v1"

var ControllerPrefixV1 string = fmt.Sprintf("/%s/bot/{token}", APIVersion1)

func RegisterAPIV1Controllers(r chi.Router) {
	r.Get(ControllerPrefixV1, InformationControllerV1)
	r.Post(ControllerPrefixV1+"/send", SendAPIHandlerV1)
	r.Get(ControllerPrefixV1+"/receive", ReceiveAPIHandlerV1)
	r.Post(ControllerPrefixV1+"/attachment", AttachmentAPIHandlerV2)
	r.Post(ControllerPrefixV1+"/webhook", WebhookController)
}

//region CONTROLLER - INFORMATION

// Renders route GET "/{version}/bot/{token}"
func InformationControllerV1(w http.ResponseWriter, r *http.Request) {

	// setting default reponse type as json
	w.Header().Set("Content-Type", "application/json")

	token := chi.URLParam(r, "token")
	bot, err := WhatsAppService.DB.Bot.FindByToken(token)
	if err != nil {
		RespondNotFound(w, fmt.Errorf("Token '%s' not found", token))
		return
	}

	var ep QPEndPoint
	ep.ID = bot.ID
	ep.Phone = bot.GetNumber()
	if bot.Verified {
		ep.Status = "verified"
	} else {
		ep.Status = "unverified"
	}

	RespondSuccess(w, ep)
}

//endregion
