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
		r.Post(endpoint+CurrentControllerPrefix+"/webhook", WebhookController)
		r.Get(endpoint+DownloadControllerEnpoint, DownloadController)
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
	if server.Status != "ready" {
		RespondNotReady(w, fmt.Errorf("bot not ready yet ! try later."))
		return
	}

	id := chi.URLParam(r, "id")
	if strings.HasPrefix(id, "message") {
		id = r.URL.Query().Get("id")
	}

	data, err := server.Download(id)
	if err != nil {
		RespondServerError(server, w, fmt.Errorf("cannot download data: %s", err))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//endregion
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

	server.Bot.WebHook = p.Url
	// Atualizando banco de dados
	if err := server.Bot.WebHookUpdate(); err != nil {
		return
	}

	RespondSuccess(w, server.Bot)
}

//endregion
