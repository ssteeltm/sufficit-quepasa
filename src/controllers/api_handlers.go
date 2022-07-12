package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"

	metrics "github.com/sufficit/sufficit-quepasa-fork/metrics"
	models "github.com/sufficit/sufficit-quepasa-fork/models"
	whatsapp "github.com/sufficit/sufficit-quepasa-fork/whatsapp"
)

const CurrentAPIVersion string = "v3"

var CurrentControllerPrefix string = "/bot/{token}"

func RegisterAPIControllers(r chi.Router) {
	aliases := []string{"/current", "", "/" + CurrentAPIVersion}
	for _, endpoint := range aliases {
		r.Get(endpoint+CurrentControllerPrefix, InformationController)
		r.Post(endpoint+CurrentControllerPrefix+"/sendtext", SendTextAPIHandlerV2)
		r.Post(endpoint+CurrentControllerPrefix+"/senddocument", SendDocumentAPIHandlerV2)

		r.Post(endpoint+CurrentControllerPrefix+"/sendurl", SendDocumentFromUrl)
		r.Post(endpoint+CurrentControllerPrefix+"/sendbinary/{chatId}/{fileName}/{textLabel}", SendDocumentFromBinary)
		r.Post(endpoint+CurrentControllerPrefix+"/sendbinary/{chatId}/{fileName}", SendDocumentFromBinary)
		r.Post(endpoint+CurrentControllerPrefix+"/sendbinary/{chatId}", SendDocumentFromBinary)
		r.Post(endpoint+CurrentControllerPrefix+"/sendbinary", SendDocumentFromBinary)
		r.Post(endpoint+CurrentControllerPrefix+"/sendencoded", SendDocumentFromEncoded)

		r.Get(endpoint+CurrentControllerPrefix+"/receive", ReceiveAPIHandler)
		r.Post(endpoint+CurrentControllerPrefix+"/attachment", AttachmentAPIHandlerV2)
		r.Get(endpoint+DownloadControllerEnpoint, DownloadController)

		r.Post(endpoint+CurrentControllerPrefix+"/webhook", WebhookController)
		r.Get(endpoint+CurrentControllerPrefix+"/webhook", WebhookController)
		r.Delete(endpoint+CurrentControllerPrefix+"/webhook", WebhookController)
	}
}

//region CONTROLLER - INFORMATION

// InformationController renders route GET "/{version}/bot/{token}"
func InformationController(w http.ResponseWriter, r *http.Request) {
	// setting default reponse type as json
	w.Header().Set("Content-Type", "application/json")

	response := &models.QpInfoResponse{}

	server, err := GetServer(w, r)
	if err != nil {
		metrics.MessageSendErrors.Inc()
		response.ParseError(err)
		RespondServerError(server, w, response)
		return
	}

	response.ParseSuccess(*server)
	RespondSuccess(w, response)
}

//endregion
//region CONTROLLER - DOWNLOAD MESSAGE ATTACHMENT

var DownloadControllerEnpoint string = CurrentControllerPrefix + "/download/{messageId}"

func GetDownloadPrefix(token string) (path string) {
	path = DownloadControllerEnpoint
	path = strings.Replace(path, "{token}", token, -1)
	path = strings.Replace(path, "{messageId}", "", -1)
	return
}

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

	server, err := GetServer(w, r)
	if err != nil {
		return
	}

	// Evitando tentativa de download de anexos sem o bot estar devidamente sincronizado
	status := server.GetStatus()
	if status != whatsapp.Ready {
		RespondNotReady(w, &ApiServerNotReadyException{Wid: server.GetWid(), Status: status})
		return
	}

	// Default parameters
	messageId := chi.URLParam(r, "id")
	if strings.HasPrefix(messageId, "message") {
		messageId = r.URL.Query().Get("id")
	} else if strings.HasPrefix(messageId, "messageId") {
		messageId = r.URL.Query().Get("messageId")
	} else if len(messageId) == 0 {
		messageId = r.Header.Get("X-QUEPASA-MESSAGEID")
	}

	att, err := server.Download(messageId)
	if err != nil {
		log.Error(err)
		RespondServerError(server, w, fmt.Errorf("cannot download data: %s", err))
		return
	}

	if len(att.FileName) > 0 {
		w.Header().Set("Content-Disposition", "attachment; filename="+att.FileName)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(*att.GetContent())
}

//endregion
