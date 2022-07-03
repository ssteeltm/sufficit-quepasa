package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"

	. "github.com/sufficit/sufficit-quepasa-fork/models"
	. "github.com/sufficit/sufficit-quepasa-fork/whatsapp"
)

const CurrentAPIVersion string = "v3"

var CurrentControllerPrefix string = "/bot/{token}"

func RegisterAPIControllers(r chi.Router) {
	aliases := []string{"/current", "", "/" + CurrentAPIVersion}
	for _, endpoint := range aliases {
		r.Get(endpoint+CurrentControllerPrefix, InformationController)
		r.Post(endpoint+CurrentControllerPrefix+"/sendtext", SendTextAPIHandlerV2)
		r.Post(endpoint+CurrentControllerPrefix+"/senddocument", SendDocumentAPIHandlerV2)
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
	token := chi.URLParam(r, "token")
	server, err := GetServerFromToken(token)
	if err != nil {
		RespondNotFound(w, fmt.Errorf("Token '%s' not found", token))
		return
	}

	RespondSuccess(w, server)
}

//endregion
//region CONTROLLER - DOWNLOAD MESSAGE ATTACHMENT

var DownloadControllerEnpoint string = CurrentControllerPrefix + "/download/{id}"

func GetDownloadPrefix(token string) (path string) {
	path = DownloadControllerEnpoint
	path = strings.Replace(path, "{token}", token, -1)
	path = strings.Replace(path, "{id}", "", -1)
	return
}

// Url Params token & id
func DownloadController(w http.ResponseWriter, r *http.Request) {

	token := chi.URLParam(r, "token")
	server, err := GetServerFromToken(token)
	if err != nil {
		RespondNotFound(w, fmt.Errorf("Token '%s' not found", token))
		return
	}

	// Evitando tentativa de download de anexos sem o bot estar devidamente sincronizado
	status := server.GetStatus()
	if status != Ready {
		RespondNotReady(w, &ApiServerNotReadyException{Wid: server.GetWid(), Status: status})
		return
	}

	id := chi.URLParam(r, "id")
	if strings.HasPrefix(id, "message") {
		id = r.URL.Query().Get("id")
	}

	att, err := server.Download(id)
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
