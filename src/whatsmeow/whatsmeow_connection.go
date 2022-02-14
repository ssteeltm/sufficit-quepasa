package whatsmeow

import (
	"context"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"

	. "github.com/sufficit/sufficit-quepasa-fork/whatsapp"
	. "go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	. "go.mau.fi/whatsmeow/types"
)

// Must Implement IWhatsappConnection
type WhatsmeowConnection struct {
	Client   *Client
	Handlers *WhatsmeowHandlers
}

//region IMPLEMENT INTERFACE WHATSAPP CONNECTION

func (conn *WhatsmeowConnection) GetVersion() string { return "multi" }

func (conn *WhatsmeowConnection) GetWid() (wid string, err error) {
	if conn.Client == nil {
		err = fmt.Errorf("client not defined on trying to get wid")
	} else {
		if conn.Client.Store == nil {
			err = fmt.Errorf("device store not defined on trying to get wid")
		} else {
			if conn.Client.Store.ID == nil {
				err = fmt.Errorf("device id not defined on trying to get wid")
			} else {
				wid = conn.Client.Store.ID.User
			}
		}
	}

	return
}

// Retorna algum titulo válido apartir de um jid
func (conn *WhatsmeowConnection) GetTitle() string {
	var result string
	/*
		contact, ok := store.Contacts[jid]
		if ok {
			result = getContactTitle(contact)
		}
	*/
	return result
}

func (conn *WhatsmeowConnection) Connect() (err error) {
	wid, err := conn.GetWid()
	if err != nil {
		return
	}

	log.Println("(%s) starting whatsmeow connecting ...", wid)

	err = conn.Client.Connect()
	if err != nil {
		return
	}

	return
}

// func (cli *Client) Download(msg DownloadableMessage) (data []byte, err error)
func (conn *WhatsmeowConnection) Download(imsg IWhatsappMessage) (data []byte, err error) {
	msg := imsg.GetSource()
	downloadable, ok := msg.(DownloadableMessage)
	if !ok {
		log.Debug("not downloadable, trying default message")
		waMsg, ok := msg.(*waProto.Message)
		if !ok {
			err = fmt.Errorf("parameter msg cannot be converted to an original message")
			return
		}
		return conn.Client.DownloadAny(waMsg)
	}
	return conn.Client.Download(downloadable)
}

// Default SEND method using WhatsappMessage Interface
func (conn *WhatsmeowConnection) Send(msg WhatsappMessage) (IWhatsappSendResponse, error) {

	response := &WhatsappSendResponse{}
	var err error

	messageText := msg.GetText()

	var newMessage *waProto.Message
	if !msg.HasAttachment() {
		internal := &waProto.ExtendedTextMessage{Text: &messageText}
		newMessage = &waProto.Message{ExtendedTextMessage: internal}
	} else {
		newMessage, err = conn.UploadAttachment(msg)
		if err != nil {
			return response, err
		}
	}

	// Formatting destination accordly
	formatedDestiantion := FormatEndpoint(msg.GetChatID())

	jid, err := ParseJID(formatedDestiantion)
	if err != nil {
		log.Printf("Send error on get jid: %s", err)
		return response, err
	}

	// Generating a new unique MessageID
	response.ID = GenerateMessageID()

	timestamp, err := conn.Client.SendMessage(jid, response.ID, newMessage)
	if err != nil {
		log.Printf("Send error: %s", err)
		return response, err
	}

	response.Timestamp = timestamp

	log.Printf("Send: %s, on: %s", response.ID, response.Timestamp)
	return response, err
}

// func (cli *Client) Upload(ctx context.Context, plaintext []byte, appInfo MediaType) (resp UploadResponse, err error)
func (conn *WhatsmeowConnection) UploadAttachment(msg WhatsappMessage) (result *waProto.Message, err error) {

	content := *msg.Attachment.GetContent()
	if content == nil {
		err = fmt.Errorf("null content")
		return
	}

	mediaType := GetMediaType(msg.Attachment.Mimetype)

	response, err := conn.Client.Upload(context.Background(), content, mediaType)
	if err != nil {
		return
	}

	result = NewWhatsmeowMessageAttachment(response, msg.Attachment)
	return
}

func (conn *WhatsmeowConnection) Disconnect() error {
	return nil
}

func (conn *WhatsmeowConnection) Delete() error {
	conn.Client.Disconnect()
	return conn.Client.Store.Delete()
}

func (conn *WhatsmeowConnection) GetWhatsAppQRChannel(result chan string) (err error) {
	// No ID stored, new login
	qrChan, _ := conn.Client.GetQRChannel(context.Background())
	err = conn.Client.Connect()
	if err != nil {
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)

	for evt := range qrChan {
		if evt.Event == "code" {
			result <- evt.Code
		} else {
			wg.Done()
			break
		}
	}

	wg.Wait()
	close(result)
	return
}

func (conn *WhatsmeowConnection) UpdateHandler(handlers IWhatsappHandlers) {
	conn.Handlers.WAHandlers = handlers
}

//endregion
