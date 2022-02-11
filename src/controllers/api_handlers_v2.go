package controllers

import (
	"fmt"

	"github.com/go-chi/chi"
)

const APIVersion2 string = "v2"

var ControllerPrefixV2 string = fmt.Sprintf("/%s/bot/{token}", APIVersion2)

func RegisterAPIV2Controllers(r chi.Router) {
	r.Get(ControllerPrefixV2, InformationControllerV1)
	r.Post(ControllerPrefixV2+"/sendtext", SendTextAPIHandlerV2)
	r.Post(ControllerPrefixV2+"/senddocument", SendDocumentAPIHandlerV2)
	r.Get(ControllerPrefixV2+"/receive", ReceiveAPIHandlerV1)
	r.Post(ControllerPrefixV2+"/attachment", AttachmentAPIHandlerV2)
	r.Post(ControllerPrefixV2+"/webhook", WebhookControllerV1)
}
