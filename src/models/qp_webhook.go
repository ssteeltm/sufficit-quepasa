package models

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"time"

	log "github.com/sirupsen/logrus"
)

type QpWebhook struct {
	Url     string     `db:"url" json:"url"`                   // destination
	Failure *time.Time `db:"failure" json:"failure,omitempty"` // first failure time
}

var ErrInvalidResponse error = errors.New("the requested url do not return 200 status code")

func (source *QpWebhook) Post(wid string, url string, message interface{}) (err error) {
	typeOfMessage := reflect.TypeOf(message)
	log.Infof("dispatching webhook from: (%s): %s, to: %s", typeOfMessage, wid, url)

	payloadJson, _ := json.Marshal(&message)
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

	if resp.StatusCode != 200 {
		err = ErrInvalidResponse
	}

	if err != nil {
		if source.Failure == nil {
			time := time.Now().UTC()
			source.Failure = &time
		}
	} else {
		source.Failure = nil
	}

	return
}
