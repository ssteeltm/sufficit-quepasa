package controllers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	models "github.com/sufficit/sufficit-quepasa/models"
)

/*
<summary>
	Find a whatsapp server by token passed on Url Path parameters
</summary>
*/
func GetServer(r *http.Request) (server *models.QPWhatsappServer, err error) {
	token := GetToken(r)
	return models.GetServerFromToken(token)
}

/*
<summary>
	Get Token From Http Request
	1º Url Param (/:token/)
	2º Url Query (?token=)
	3º Header (X-QUEPASA-TOKEN)
</summary>
*/
func GetToken(r *http.Request) (result string) {

	// retrieve from url path parameter
	result = chi.URLParam(r, "token")
	if len(result) == 0 {

		// retrieve from url query parameter
		if r.URL.Query().Has("token") {

			result = r.URL.Query().Get("token")
		} else {

			// retrieve from header parameter
			result = r.Header.Get("X-QUEPASA-TOKEN")
		}
	}
	return
}

//<summary>Find a whatsapp server by token passed on Url Path parameters</summary>
func GetServerRespondOnError(w http.ResponseWriter, r *http.Request) (server *models.QPWhatsappServer, err error) {
	token := chi.URLParam(r, "token")
	server, err = models.GetServerFromToken(token)
	if err != nil {
		RespondNotFound(w, fmt.Errorf("token '%s' not found", token))
	}
	return
}
