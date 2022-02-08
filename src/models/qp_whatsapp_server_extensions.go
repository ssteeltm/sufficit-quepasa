package models

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Encaminha msg ao WebHook específicado
func PostToWebHookFromServer(server *QPWhatsappServer, message interface{}) error {
	bot := server.Bot
	if bot == nil {
		return fmt.Errorf("cannot find bot for server")
	}

	if len(bot.WebHook) > 0 {
		log.Info("dispatching webhook from: ", server.GetWid())

		payloadJson, _ := json.Marshal(&message)
		requestBody := bytes.NewBuffer(payloadJson)

		// Ignorando certificado ao realizar o post
		// Não cabe a nós a segurança do cliente
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		resp, err := http.Post(bot.WebHook, "application/json", requestBody)
		if err != nil {
			log.Printf("(%s) erro ao postar no webhook: %s", bot.GetNumber(), err.Error())
		} else {
			if resp != nil {
				defer resp.Body.Close()
				if resp.StatusCode == 422 {
					body, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						log.Printf("(%s) erro ao ler resposta do webhook: %s", bot.GetNumber(), err.Error())
					} else {
						if body != nil && strings.Contains(string(body), "invalid callback token") {

							// Sincroniza o token mais novo
							bot.WebHookSincronize()

							// Preenche o body novamente pois foi esvaziado na requisição anterior
							requestBody = bytes.NewBuffer(payloadJson)
							http.Post(bot.WebHook, "application/json", requestBody)
						}
					}
				}
			}
		}
	}
	return nil
}

//region FIND|SEARCH WHATSAPP SERVER

var ServerNotFoundError error = errors.New("the requested whatsapp server was not founded")

func GetServerFromID(source string) (server *QPWhatsappServer, err error) {
	server, ok := WhatsAppService.Servers[source]
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
	for _, item := range WhatsAppService.Servers {
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

//endregion
