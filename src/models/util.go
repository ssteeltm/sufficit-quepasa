package models

import (
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
)

// GetUser gets the user_id from the JWT and finds the
// corresponding user in the database
func GetUser(r *http.Request) (QPUser, error) {
	var user QPUser
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		return user, err
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return user, errors.New("User ID missing")
	}

	return WhatsappService.DB.User.FindByID(userID)
}

// Usado também para identificar o número do bot
// Meramente visual
func GetPhoneByID(id string) (out string, err error) {

	// removing whitespaces
	out = strings.Replace(id, " ", "", -1)
	if strings.Contains(out, "@") {
		// capturando tudo antes do @
		splited := strings.Split(out, "@")
		out = splited[0]

		if strings.Contains(out, ".") {
			// capturando tudo antes do "."
			splited = strings.Split(out, ".")
			out = splited[0]

			return
		}
	}

	re, err := regexp.Compile(`\d*`)
	matches := re.FindAllString(out, -1)
	if len(matches) > 0 {
		out = matches[0]
	}
	return out, err
}

// Getting ChatId from PATH => QUERY => HEADER
func GetChatId(r *http.Request) (result string) {

	// retrieve from url path parameter
	result = chi.URLParam(r, "chatid")
	if len(result) == 0 {

		// retrieve from url query parameter
		if r.URL.Query().Has("chatid") {
			result = r.URL.Query().Get("chatid")
		} else {

			// retrieve from header parameter
			result = r.Header.Get("X-QUEPASA-CHATID")
		}
	}
	return
}
