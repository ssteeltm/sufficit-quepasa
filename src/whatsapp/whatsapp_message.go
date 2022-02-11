package whatsapp

import (
	"time"
)

// Mensagem no formato QuePasa
// Utilizada na API do QuePasa para troca com outros sistemas
type WhatsappMessage struct {

	// original message from source service
	Content interface{} `json:"-"`

	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`

	// Se a msg foi postado em algum grupo ? quem postou !
	Participant WhatsappEndpoint `json:"participant,omitempty"`

	// Fui eu quem enviou a msg ?
	FromMe bool `json:"fromme"`

	// Texto da msg
	Text string `json:"text"`

	Attachment *WhatsappAttachment `json:"attachment,omitempty"`

	Chat WhatsappChat `json:"chat"`
}

//region ORDER BY TIMESTAMP

type ByTimestamp []WhatsappMessage

func (m ByTimestamp) Len() int           { return len(m) }
func (m ByTimestamp) Less(i, j int) bool { return m[i].Timestamp.After(m[j].Timestamp) }
func (m ByTimestamp) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }

//endregion

//region IMPLEMENT WHATSAPP SEND RESPONSE INTERFACE

func (source *WhatsappMessage) GetID() string { return source.ID }

// Get the time of server processed message
func (source *WhatsappMessage) GetTime() time.Time { return source.Timestamp }

// Get the time on unix timestamp format
func (source *WhatsappMessage) GetTimestamp() uint64 { return uint64(source.Timestamp.Unix()) }

//endregion

func (source *WhatsappMessage) GetChatID() string {
	return source.Chat.ID
}

func (source *WhatsappMessage) GetText() string {
	return source.Text
}

func (source *WhatsappMessage) HasAttachment() bool {
	// this attachment is a pointer to correct show info on deserialized
	attach := source.Attachment
	return attach != nil && len(attach.Mimetype) > 0
}

func (source *WhatsappMessage) GetSource() interface{} {
	return source.Content
}
