package models

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"time"

	log "github.com/sirupsen/logrus"
)

// Encaminha msg ao WebHook específicado
func PostToWebHookFromServer(server *QPWhatsappServer, message interface{}) (err error) {
	wid := server.GetWid()
	url := server.Webhook()

	// Ignorando certificado ao realizar o post
	// Não cabe a nós a segurança do cliente
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	if len(url) > 0 {
		return PostToWebHook(wid, url, message)
	}

	return
}

func PostToWebHook(wid string, url string, message interface{}) (err error) {
	typeOfMessage := reflect.TypeOf(message)
	log.Infof("dispatching webhook from: (%s): %s", typeOfMessage, wid)

	payloadJson, _ := json.Marshal(&message)
	log.Debug()

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadJson))
	req.Header.Set("User-Agent", "Quepasa")
	req.Header.Set("X-QUEPASA-BOT", wid)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	client.Timeout = time.Second * 10
	resp, err := client.Do(req)
	if err != nil {
		log.Error("(%s) erro ao postar no webhook: %s", wid, err.Error())
	}
	defer resp.Body.Close()
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
