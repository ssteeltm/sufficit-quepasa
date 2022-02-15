package controllers

import (
	"encoding/json"
	"fmt"

	"net/http"

	"github.com/go-chi/chi"
	. "github.com/sufficit/sufficit-quepasa-fork/metrics"
	. "github.com/sufficit/sufficit-quepasa-fork/models"
	. "github.com/sufficit/sufficit-quepasa-fork/whatsapp"
)

// SendAPIHandler renders route "/v1/bot/{token}/send"
func SendAPIHandlerV1(w http.ResponseWriter, r *http.Request) {

	// setting default reponse type as json
	w.Header().Set("Content-Type", "application/json")

	token := chi.URLParam(r, "token")
	server, err := GetServerFromToken(token)
	if err != nil {
		RespondNotFound(w, fmt.Errorf("Token '%s' not found", token))
		return
	}

	// Declare a new Person struct.
	var request QPSendRequestV1

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		RespondServerError(server, w, err)
		return
	}

	msg, err := ToWhatsappMessageV1(&request)
	if err != nil {
		RespondServerError(server, w, err)
		return
	}

	err = server.SendMessage(msg)
	if err != nil {
		RespondServerError(server, w, err)
		return
	}

	MessagesSent.Inc()
	res := &QPSendResponse{
		Result: &QPSendResult{
			Source:    server.GetWid(),
			Recipient: request.Recipient,
			MessageId: msg.GetID(),
		},
	}

	RespondSuccess(w, res)
}

// Renders route GET "/{version}/bot/{token}/receive"
func ReceiveAPIHandlerV1(w http.ResponseWriter, r *http.Request) {

	// setting default reponse type as json
	w.Header().Set("Content-Type", "application/json")

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

	queryValues := r.URL.Query()
	timestamp := queryValues.Get("timestamp")

	messages, err := GetMessagesV1(server.Bot, timestamp)
	if err != nil {
		MessageReceiveErrors.Inc()
		RespondServerError(server, w, err)
		return
	}

	MessagesReceived.Add(float64(len(messages)))

	out := QPFormReceiveResponseV1{
		Bot:      *ToQPBotV1(server.Bot),
		Messages: messages,
	}

	RespondSuccess(w, out)
}
