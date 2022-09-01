package controllers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	models "github.com/sufficit/sufficit-quepasa/models"
)

//<summary>Find a whatsapp server by token passed on Url Path parameters</summary>
func GetServer(w http.ResponseWriter, r *http.Request) (server *models.QPWhatsappServer, err error) {
	token := chi.URLParam(r, "token")
	server, err = models.GetServerFromToken(token)
	if err != nil {
		RespondNotFound(w, fmt.Errorf("token '%s' not found", token))
	}
	return
}

//<summary>Find a system track identifier to follow the message</summary>
func GetTrackId(r *http.Request) string {
	return r.Header.Get("X-QUEPASA-TRACKID")
}
