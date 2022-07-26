package models

type QpSendResponseMessage struct {
	Id        string `json:"id,omitempty"`
	Source    string `json:"source,omitempty"`
	Recipient string `json:"recipient,omitempty"`
}
