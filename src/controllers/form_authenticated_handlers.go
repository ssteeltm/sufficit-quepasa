package controllers

import(
	"os"
	"html/template"
	"net/http"
	
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"

	. "github.com/sufficit/sufficit-quepasa-fork/models"
)

// Token of authentication / encryption
var TokenAuth = jwtauth.New("HS256", []byte(os.Getenv("SIGNING_SECRET")), nil)

// Prefix on forms endpoints to avoid conflict with api
const FormEndpointPrefix string = "/form"
var FormWebsocketEndpoint string = FormEndpointPrefix + "/verify/ws"
var FormAccountEndpoint string = FormEndpointPrefix + "/account"
var FormVerifyEndpoint string = FormEndpointPrefix + "/verify"
var FormDeleteEndpoint string = FormEndpointPrefix + "/delete"

func RegisterFormAuthenticatedControllers(r chi.Router) {
	r.Use(jwtauth.Verifier(TokenAuth))
	r.Use(HttpAuthenticatorHandler)

	r.HandleFunc(FormWebsocketEndpoint, VerifyHandler)

	r.Get(FormAccountEndpoint, FormAccountController)
	r.Get(FormVerifyEndpoint, VerifyFormHandler)
	r.Post(FormDeleteEndpoint, DeleteHandler)
	r.Post(FormEndpointPrefix + "/cycle", CycleHandler)
	r.Post(FormEndpointPrefix + "/debug", DebugHandler)
	r.Post(FormEndpointPrefix + "/toggle", ToggleHandler)
	r.Get(FormEndpointPrefix + "/bot/{id}", FormSendController)
	r.Get(FormEndpointPrefix + "/bot/{id}/send", FormSendController)
	r.Post(FormEndpointPrefix + "/bot/{id}/send", FormSendController)
	r.Get(FormEndpointPrefix + "/bot/{id}/receive", FormReceiveController)
}

// Authentication manager on forms
func HttpAuthenticatorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())

		if err != nil {
			http.Redirect(w, r, FormLoginEndpoint, http.StatusFound)
			return
		}

		if token == nil || !token.Valid {
			http.Redirect(w, r, FormLoginEndpoint, http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}


// Rrenders route GET "/{prefix}/account"
func FormAccountController(w http.ResponseWriter, r *http.Request) {
	user, err := GetUser(r)
	if err != nil {
		RedirectToLogin(w, r)
	}

	data := QPFormAccountData{
		PageTitle: "Account",
		User:      user,
	}

	bots, err := WhatsAppService.DB.Bot.FindAllForUser(user.ID)
	if err != nil {
		data.ErrorMessage = err.Error()
	} else {
		data.Bots = bots
	}

	templates := template.Must(template.ParseFiles("views/layouts/main.tmpl", "views/account.tmpl"))
	templates.ExecuteTemplate(w, "main", data)
}