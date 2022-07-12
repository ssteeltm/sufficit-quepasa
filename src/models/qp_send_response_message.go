package models

type QpSendResponseMessage struct {
	Id        string `json:"id"`
	Source    string `json:"source"`
	Recipient string `json:"recipient"`
}
