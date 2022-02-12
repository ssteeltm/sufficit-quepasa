package models

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Encaminha msg ao WebHook específicado
func PostToWebHookFromServer(server *QPWhatsappServer, message interface{}) (err error) {
	wid := server.GetWid()
	url := server.Webhook()

	if len(url) > 0 {
		return PostToWebHook(wid, url, message)
	}

	return
}

func PostToWebHook(wid string, url string, message interface{}) (err error) {
	log.Info("dispatching webhook from: ", wid)

	payloadJson, _ := json.Marshal(&message)
	log.Debug(string(payloadJson))

	requestBody := bytes.NewBuffer(payloadJson)

	// Ignorando certificado ao realizar o post
	// Não cabe a nós a segurança do cliente
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp, err := http.Post(url, "application/json", requestBody)
	if err != nil {
		log.Error("(%s) erro ao postar no webhook: %s", wid, err.Error())
	} else {
		if resp != nil {
			defer resp.Body.Close()
			if resp.StatusCode == 422 {
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Printf("(%s) erro ao ler resposta do webhook: %s", wid, err.Error())
				} else {
					if body != nil && strings.Contains(string(body), "invalid callback token") {
						// Preenche o body novamente pois foi esvaziado na requisição anterior
						requestBody = bytes.NewBuffer(payloadJson)
						http.Post(url, "application/json", requestBody)
					}
				}
			}
		}
	}

	return
}

//region FIND|SEARCH WHATSAPP SERVER

var ServerNotFoundError error = errors.New("the requested whatsapp server was not founded")

func GetServerFromID(source string) (server *QPWhatsappServer, err error) {
	server, ok := WhatsappService.Servers[source]
	if !ok {
		err = ServerNotFoundError
		return
	}
	return
}

func GetServerFromBot(source QPBot) (server *QPWhatsappServer, err error) {
	return GetServerFromID(source.ID)
}

func GetServerFromToken(token string) (server *QPWhatsappServer, err error) {
	for _, item := range WhatsappService.Servers {
		if item.Bot != nil && item.Bot.Token == token {
			server = item
			break
		}
	}
	if server == nil {
		err = ServerNotFoundError
	}
	return
}

func GetServersForUserID(userid string) (servers map[string]*QPWhatsappServer) {
	return WhatsappService.GetServersForUser(userid)
}

func GetServersForUser(user QPUser) (servers map[string]*QPWhatsappServer) {
	return GetServersForUserID(user.ID)
}

//endregion
