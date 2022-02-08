package controllers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	. "github.com/sufficit/sufficit-quepasa-fork/models"
)

const APIVersion2 string = "v2"

var ControllerPrefixV2 string = fmt.Sprintf("/%s/bot/{token}", APIVersion2)

func RegisterAPIV2Controllers(r chi.Router) {
	r.Get(ControllerPrefixV2, InformationControllerV2)
	r.Post(ControllerPrefixV2+"/sendtext", SendTextAPIHandlerV2)
	r.Post(ControllerPrefixV2+"/senddocument", SendDocumentAPIHandlerV2)
	r.Get(ControllerPrefixV2+"/receive", ReceiveAPIHandlerV2)
	r.Post(ControllerPrefixV2+"/attachment", AttachmentAPIHandlerV2)
	r.Post(ControllerPrefixV2+"/webhook", WebhookController)
}

//region CONTROLLER - INFORMATION

// InformationController renders route GET "/{version}/bot/{token}"
func InformationControllerV2(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	server, err := GetServerFromToken(token)
	if err != nil {
		RespondNotFound(w, fmt.Errorf("Token '%s' not found", token))
		return
	}

	RespondSuccess(w, server.Bot)
}

//endregion
