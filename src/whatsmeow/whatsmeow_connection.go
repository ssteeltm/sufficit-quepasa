package whatsmeow

import (
	"context"
	"fmt"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"

	whatsapp "github.com/sufficit/sufficit-quepasa-fork/whatsapp"
	. "go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	. "go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
)

// Must Implement IWhatsappConnection
type WhatsmeowConnection struct {
	Client   *Client
	Handlers *WhatsmeowHandlers
	waLogger waLog.Logger
	logger   *log.Logger
	log      *log.Entry
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

func (conn *WhatsmeowConnection) GetStatus() (state whatsapp.WhatsappConnectionState) {
	if conn != nil {
		state = whatsapp.Created
		if conn.Client != nil {
			if conn.Client.IsConnected() {
				state = whatsapp.Connected
				if conn.Client.IsLoggedIn() {
					state = whatsapp.Ready
				}
			} else {
				if conn.Client.IsLoggedIn() {
					state = whatsapp.Disconnected
				} else {
					state = whatsapp.Failed
				}
			}
		}
	}
	return
}

// Retorna algum titulo válido apartir de um jid
func (conn *WhatsmeowConnection) GetTitle(Wid string) string {
	jid := NewJID(Wid, "")
	var result string
	contact, err := conn.Client.Store.Contacts.GetContact(jid)
	if err == nil {
		result = contact.PushName
	}

	return result
}

func (conn *WhatsmeowConnection) Connect() (err error) {
	conn.log.Info("starting whatsmeow connecting ...")

	err = conn.Client.Connect()
	if err != nil {
		if strings.Contains(err.Error(), "401") {
			return &whatsapp.UnauthorizedError{Inner: err}
		}
		return
	}

	return
}

// func (cli *Client) Download(msg DownloadableMessage) (data []byte, err error)
func (conn *WhatsmeowConnection) DownloadData(imsg whatsapp.IWhatsappMessage) (data []byte, err error) {
	msg := imsg.GetSource()
	downloadable, ok := msg.(DownloadableMessage)
	if !ok {
		conn.log.Debug("not downloadable, trying default message")
		waMsg, ok := msg.(*waProto.Message)
		if !ok {
			err = fmt.Errorf("parameter msg cannot be converted to an original message")
			return
		}
		return conn.Client.DownloadAny(waMsg)
	}
	return conn.Client.Download(downloadable)
}

func (conn *WhatsmeowConnection) Download(imsg whatsapp.IWhatsappMessage) (att whatsapp.WhatsappAttachment, err error) {
	data, err := conn.DownloadData(imsg)
	if err != nil {
		return
	}

	att.SetContent(&data)
	return
}

// Default SEND method using WhatsappMessage Interface
func (conn *WhatsmeowConnection) Send(msg whatsapp.WhatsappMessage) (whatsapp.IWhatsappSendResponse, error) {

	response := &whatsapp.WhatsappSendResponse{}
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
	formatedDestiantion, _ := whatsapp.FormatEndpoint(msg.GetChatID())
	jid, err := ParseJID(formatedDestiantion)
	if err != nil {
		conn.log.Infof("Send error on get jid: %s", err)
		return response, err
	}

	// Generating a new unique MessageID
	response.ID = GenerateMessageID()

	timestamp, err := conn.Client.SendMessage(jid, response.ID, newMessage)
	if err != nil {
		conn.log.Infof("Send error: %s", err)
		return response, err
	}

	response.Timestamp = timestamp

	conn.log.Infof("Send: %s, on: %s", response.ID, response.Timestamp)
	return response, err
}

// func (cli *Client) Upload(ctx context.Context, plaintext []byte, appInfo MediaType) (resp UploadResponse, err error)
func (conn *WhatsmeowConnection) UploadAttachment(msg whatsapp.WhatsappMessage) (result *waProto.Message, err error) {

	content := *msg.Attachment.GetContent()
	if content == nil {
		err = fmt.Errorf("null content")
		return
	}

	mediaType := GetMediaTypeFromString(msg.Attachment.Mimetype)

	response, err := conn.Client.Upload(context.Background(), content, mediaType)
	if err != nil {
		return
	}

	result = NewWhatsmeowMessageAttachment(response, msg.Attachment)
	return
}

func (conn *WhatsmeowConnection) Disconnect() error {
	conn.Client.Disconnect()
	return nil
}

func (conn *WhatsmeowConnection) Delete() error {
	conn.Client.Disconnect()
	return conn.Client.Store.Delete()
}

func (conn *WhatsmeowConnection) GetWhatsAppQRChannel(result chan<- string) (err error) {
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

func (conn *WhatsmeowConnection) UpdateHandler(handlers whatsapp.IWhatsappHandlers) {
	conn.Handlers.WAHandlers = handlers
}

//endregion

func (conn *WhatsmeowConnection) LogLevel(level log.Level) {
	conn.logger.SetLevel(level)
}

func (conn *WhatsmeowConnection) PrintStatus() {
	/*
		conn.log.Warnf("STATUS IS CONNECTED: %v", conn.Client.IsConnected())
		conn.log.Warnf("STATUS IS LOGGED IN: %v", conn.Client.IsLoggedIn())

		conn.Client.SendPresence(types.PresenceAvailable)
		conn.Client.SetPassive(false)
	*/
}
