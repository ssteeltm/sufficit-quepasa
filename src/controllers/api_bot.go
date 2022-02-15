package controllers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	. "github.com/sufficit/sufficit-quepasa-fork/metrics"
	. "github.com/sufficit/sufficit-quepasa-fork/models"
	. "github.com/sufficit/sufficit-quepasa-fork/whatsapp"
)

// ReceiveAPIHandler renders route GET "/v1/bot/{token}/receive"
func ReceiveAPIHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("ReceiveAPIHandlerV3: %+v\n", r)

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

	messages, err := GetMessages(server, timestamp)
	if err != nil {
		MessageReceiveErrors.Inc()
		RespondServerError(server, w, err)
		return
	}

	MessagesReceived.Add(float64(len(messages)))

	out := QPFormReceiveResponse{
		Bot:      *server.Bot,
		Messages: messages,
	}

	RespondSuccess(w, out)
}
