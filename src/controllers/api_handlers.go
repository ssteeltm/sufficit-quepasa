package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	metrics "github.com/sufficit/sufficit-quepasa/metrics"
	models "github.com/sufficit/sufficit-quepasa/models"
	whatsapp "github.com/sufficit/sufficit-quepasa/whatsapp"
)

const CurrentAPIVersion string = "v4"

func RegisterAPIControllers(r chi.Router) {
	aliases := []string{"/current", "", "/" + CurrentAPIVersion}
	for _, endpoint := range aliases {

		// CONTROL METHODS ************************
		// ----------------------------------------
		r.Get(endpoint+"/info", InformationController)
		r.Get(endpoint+"/scan", ScannerController)

		// ----------------------------------------
		// CONTROL METHODS ************************

		// SENDING MSG ----------------------------
		// ----------------------------------------

		// used to dispatch alert msgs via url, triggers on monitor systems like zabbix
		r.Get(endpoint+"/send", SendAny)

		r.Post(endpoint+"/send", SendAny)
		r.Post(endpoint+"/send/{chatid}", SendAny)
		r.Post(endpoint+"/sendtext", SendText)
		r.Post(endpoint+"/sendtext/{chatid}", SendText)

		// SENDING MSG ATTACH ---------------------

		// deprecated, discard/remove on next version
		r.Post(endpoint+"/senddocument", SendDocumentAPIHandlerV2)

		r.Post(endpoint+"/sendurl", SendDocumentFromUrl)
		r.Post(endpoint+"/sendbinary/{chatid}/{fileName}/{text}", SendDocumentFromBinary)
		r.Post(endpoint+"/sendbinary/{chatid}/{fileName}", SendDocumentFromBinary)
		r.Post(endpoint+"/sendbinary/{chatid}", SendDocumentFromBinary)
		r.Post(endpoint+"/sendbinary", SendDocumentFromBinary)
		r.Post(endpoint+"/sendencoded", SendDocumentFromEncoded)

		// ----------------------------------------
		// SENDING MSG ----------------------------

		r.Get(endpoint+"/receive", ReceiveAPIHandler)
		r.Post(endpoint+"/attachment", AttachmentAPIHandlerV2)

		r.Get(endpoint+"/download/{messageId}", DownloadController)
		r.Get(endpoint+"/download", DownloadController)

		// PICTURE INFO | DATA --------------------
		// ----------------------------------------

		r.Post(endpoint+"/picinfo", PictureController)
		r.Get(endpoint+"/picinfo/{chatid}/{pictureid}", PictureController)
		r.Get(endpoint+"/picinfo/{chatid}", PictureController)
		r.Get(endpoint+"/picinfo", PictureController)

		r.Post(endpoint+"/picdata", PictureController)
		r.Get(endpoint+"/picdata/{chatid}/{pictureid}", PictureController)
		r.Get(endpoint+"/picdata/{chatid}", PictureController)
		r.Get(endpoint+"/picdata", PictureController)

		// ----------------------------------------
		// PICTURE INFO | DATA --------------------

		r.Post(endpoint+"/webhook", WebhookController)
		r.Get(endpoint+"/webhook", WebhookController)
		r.Delete(endpoint+"/webhook", WebhookController)

		// INVITE METHODS ************************
		// ----------------------------------------

		r.Get(endpoint+"/invite", InviteController)
		r.Get(endpoint+"/invite/{chatid}", InviteController)

		// ----------------------------------------
		// INVITE METHODS ************************
	}
}

//region CONTROLLER - INFORMATION

// InformationController renders route GET "/{version}/info"
func InformationController(w http.ResponseWriter, r *http.Request) {
	// setting default reponse type as json
	w.Header().Set("Content-Type", "application/json")

	response := &models.QpInfoResponse{}

	server, err := GetServer(r)
	if err != nil {
		metrics.MessageSendErrors.Inc()
		response.ParseError(err)
		RespondInterface(w, response)
		return
	}

	response.ParseSuccess(*server)
	RespondSuccess(w, response)
}

//endregion
//region CONTROLLER - DOWNLOAD MESSAGE ATTACHMENT

/*
<summary>
	Renders route GET "/{{version}}/bot/{{token}}/download/{messageId}"

	Any of then, at this order of priority
	Path parameters: {messageId}
	Url parameters: ?messageId={messageId} || ?id={messageId}
	Header parameters: X-QUEPASA-MESSAGEID = {messageId}
</summary>
*/
func DownloadController(w http.ResponseWriter, r *http.Request) {

	response := &models.QpResponse{}

	server, err := GetServer(r)
	if err != nil {
		response.ParseError(err)
		RespondInterface(w, response)
		return
	}

	// Evitando tentativa de download de anexos sem o bot estar devidamente sincronizado
	status := server.GetStatus()
	if status != whatsapp.Ready {
		err := &ApiServerNotReadyException{Wid: server.GetWid(), Status: status}
		response.ParseError(err)
		RespondInterface(w, response)
		return
	}

	// Default parameters
	messageId := chi.URLParam(r, "messageId")
	if strings.Contains(messageId, "message") || (len(messageId) == 0 && r.URL.Query().Has("id")) {
		messageId = r.URL.Query().Get("id")
	} else if len(messageId) == 0 && r.URL.Query().Has("messageId") {
		messageId = r.URL.Query().Get("messageId")
	} else if len(messageId) == 0 {
		messageId = r.Header.Get("X-QUEPASA-MESSAGEID")
	}

	if len(messageId) == 0 {
		metrics.MessageSendErrors.Inc()
		err := fmt.Errorf("empty message id")
		response.ParseError(err)
		RespondInterface(w, response)
		return
	}

	att, err := server.Download(messageId)
	if err != nil {
		response.ParseError(err)
		RespondInterface(w, response)
		return
	}

	if len(att.FileName) > 0 {
		w.Header().Set("Content-Disposition", "attachment; filename="+att.FileName)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(*att.GetContent())
}

//endregion

func ScannerController(w http.ResponseWriter, r *http.Request) {
	// setting default reponse type as json
	w.Header().Set("Content-Type", "application/json")

	response := &models.QpInfoResponse{}

	server, err := GetServer(r)
	if err != nil {
		metrics.MessageSendErrors.Inc()
		response.ParseError(err)
		RespondServerError(server, w, response)
		return
	}

	response.ParseSuccess(*server)
	RespondSuccess(w, response)
}
