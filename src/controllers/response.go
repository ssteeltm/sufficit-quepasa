package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	. "github.com/sufficit/sufficit-quepasa-fork/models"
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

	RespondError(w, err, http.StatusBadRequest)
}

func RespondUnauthorized(w http.ResponseWriter, err error) {
	log.Println("!Request Unauthorized: ", err)

	RespondError(w, err, http.StatusUnauthorized)
}

func RespondNotFound(w http.ResponseWriter, err error) {
	log.Println("!Request Not found: ", err)

	RespondError(w, err, http.StatusNotFound)
}

/// Usado para avisar que o bot ainda não esta pronto
func RespondNotReady(w http.ResponseWriter, err error) {
	RespondError(w, err, http.StatusServiceUnavailable)
}

func RespondServerError(server *QPWhatsappServer, w http.ResponseWriter, err error) {
	if strings.Contains(err.Error(), "invalid websocket") {

		// Desconexão forçado é algum evento iniciado pelo whatsapp
		log.Printf("(%s) Desconexão forçada por motivo de websocket inválido ou sem resposta", server.GetWid())
		go server.Restart()

	} else {
		if ENV.DEBUGRequests() {
			log.Printf("(%s) !Request Server error: %s", server.GetWid(), err)
		}
	}
	RespondError(w, err, http.StatusInternalServerError)
}

func RespondError(w http.ResponseWriter, err error, code int) {
	res := &errorResponse{
		Result: err.Error(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(res)
}
