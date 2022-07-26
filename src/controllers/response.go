package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	models "github.com/sufficit/sufficit-quepasa-fork/models"
)

type errorResponse struct {
	Result string `json:"result"`
}

func RespondSuccess(w http.ResponseWriter, res interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func RespondBadRequest(w http.ResponseWriter, err error) {
	log.Println("!Request Bad Format: ", err)

	RespondErrorCode(w, err, http.StatusBadRequest)
}

func RespondUnauthorized(w http.ResponseWriter, err error) {
	log.Println("!Request Unauthorized: ", err)

	RespondErrorCode(w, err, http.StatusUnauthorized)
}

func RespondNotFound(w http.ResponseWriter, err error) {
	log.Println("!Request Not found: ", err)

	RespondErrorCode(w, err, http.StatusNotFound)
}

/// Usado para avisar que o bot ainda não esta pronto
func RespondNotReady(w http.ResponseWriter, err error) {
	RespondErrorCode(w, err, http.StatusServiceUnavailable)
}

func RespondServerError(server *models.QPWhatsappServer, w http.ResponseWriter, err error) {
	if strings.Contains(err.Error(), "invalid websocket") {

		// Desconexão forçado é algum evento iniciado pelo whatsapp
		log.Printf("(%s) Desconexão forçada por motivo de websocket inválido ou sem resposta", server.GetWid())
		go server.Restart()

	} else {
		if models.ENV.DEBUGRequests() {
			log.Printf("(%s) !Request Server error: %s", server.GetWid(), err)
		}
	}
	RespondErrorCode(w, err, http.StatusInternalServerError)
}

func RespondErrorCode(w http.ResponseWriter, err error, code int) {
	res := &errorResponse{
		Result: err.Error(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(res)
}

/*
<summary>
	Default response method
	Used in v3 *models.QpResponse
	Returns OK | Bad Request
</summary>
*/
func RespondInterfaceCode(w http.ResponseWriter, response interface{}, code int) {

	// setting default reponse type as json
	w.Header().Set("Content-Type", "application/json")

	// Writing header code
	if code != 0 {
		w.WriteHeader(code)
	} else {
		if qpresponse, ok := response.(models.QpResponseInterface); ok {
			if qpresponse.IsSuccess() {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		}
	}

	json.NewEncoder(w).Encode(response)
}

func RespondInterface(w http.ResponseWriter, response interface{}) {
	RespondInterfaceCode(w, response, 0)
}
