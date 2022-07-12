package models

import (
	"io/ioutil"
	"net/http"
)

type QpSendRequestUrl struct {
	QpSendRequest
	Url string `json:"url"`
}

func (source *QpSendRequestUrl) GenerateContent() (err error) {
	resp, err := http.Get(source.Url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	source.QpSendRequest.Content = content
	return
}
