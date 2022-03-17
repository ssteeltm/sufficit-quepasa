package whatsrhymen

import (
	"fmt"
	"strings"
	"sync"
	"time"

	whatsrhymen "github.com/Rhymen/go-whatsapp"
	log "github.com/sirupsen/logrus"
	whatsapp "github.com/sufficit/sufficit-quepasa-fork/whatsapp"
)

// Must Implement IWhatsappConnection
type WhatsrhymenConnection struct {
	Client         *whatsrhymen.Conn
	Handlers       *WhatsrhymenHandlers
	Session        *whatsrhymen.Session
	WAHandlers     whatsapp.IWhatsappHandlers
	Reconnect      bool
	log            *log.Entry
	failedToken    bool
	syncConnection *sync.Mutex `json:"-"` // Objeto de sinaleiro para evitar chamadas simultâneas a este objeto
}

func (conn *WhatsrhymenConnection) GetVersion() string { return "single" }

func (conn *WhatsrhymenConnection) GetWid() (wid string, err error) {
	if conn.Client == nil {
		if conn.Session != nil {
			wid = FormatWid(conn.Session.Wid)
			return
		}

		err = fmt.Errorf("client not defined on trying to get wid")
	} else {
		if conn.Client.Info == nil {
			if conn.Session != nil {
				wid = FormatWid(conn.Session.Wid)
				return
			}
			err = fmt.Errorf("session|&|client info not defined on trying to get wid")
		} else {
			wid = FormatWid(conn.Client.Info.Wid)
			return
		}
	}

	return
}

func (conn *WhatsrhymenConnection) GetStatus() (state whatsapp.WhatsappConnectionState) {
	if conn != nil {
		state = whatsapp.Created
		if conn.Client != nil {
			state = whatsapp.Stopped
			if conn.Client.GetConnected() {
				state = whatsapp.Connected
				if conn.failedToken {
					state = whatsapp.Failed
				}
				if conn.Client.GetLoggedIn() {
					// indicates that underlying connection is ok & we are loggedin
					// lets see the handlers ....
					state = whatsapp.Fetching

					if conn.WAHandlers != nil && conn.Handlers != nil {
						state = whatsapp.Ready
					}
				}
			} else {
				state = whatsapp.Disconnected
				if conn.failedToken {
					return whatsapp.Failed
				}
			}
		} else {
			state = whatsapp.Starting
			if conn.failedToken {
				return whatsapp.Failed
			}
		}
	}
	return
}

// Retorna algum titulo válido apartir de um jid
func (conn *WhatsrhymenConnection) GetTitle(Wid string) string {
	var result string
	contact, ok := conn.Client.Store.Contacts[Wid]
	if ok {
		result = getContactTitle(contact)
	}
	return result
}

func (conn *WhatsrhymenConnection) Connect() (err error) {
	err = conn.UpdateClient()
	if err != nil {
		return
	}

	// if not logged
	if !conn.Client.GetLoggedIn() {
		conn.log.Debug("restoring logged state with saved session ...")

		// Agora sim, restaura a conexão com o whatsapp apartir de uma seção salva
		_, err = conn.Client.RestoreWithSession(*conn.Session)
		if err != nil {
			conn.log.Errorf("error on restore whatsrhymen: %s", err.Error())
			if strings.Contains(err.Error(), "401") {
				conn.Failed("invalid unauthorized session")
				inner := &whatsapp.UnauthorizedError{Inner: err}

				// avoid recurring erros
				conn.Client.RemoveHandlers()

				err = conn.Client.Logout()
				if err != nil {
					return
				}

				return inner
			}

			conn.Failed("unknown")
			return
		}
	}

	conn.Success()
	return
}

func (conn *WhatsrhymenConnection) FindMessage(imsg whatsapp.IWhatsappMessage) (msg *WhatsrhymenMessage, err error) {
	source := imsg.GetSource()
	if source == nil {
		err = fmt.Errorf("cannot get a valid source for download on: %s", imsg.GetChatID())
		return
	}

	msg, ok := source.(*WhatsrhymenMessage)
	if !ok {
		err = fmt.Errorf("cannot convert interface in WhatsrhymenMessage from: %s", imsg.GetChatID())
		return
	}
	return
}

// func (cli *Client) Download(msg DownloadableMessage) (data []byte, err error)
func (conn *WhatsrhymenConnection) DownloadData(imsg whatsapp.IWhatsappMessage) (data []byte, err error) {
	msg, err := conn.FindMessage(imsg)
	if err != nil {
		return
	}

	conn.log.Tracef("downloading msg from %s", imsg.GetChatID())
	return whatsrhymen.Download(msg.AttachmentInfo.Url, msg.AttachmentInfo.MediaKey, msg.AttachmentInfo.MediaType, msg.AttachmentInfo.Length)
}

// func (cli *Client) Download(msg DownloadableMessage) (data []byte, err error)
func (conn *WhatsrhymenConnection) Download(imsg whatsapp.IWhatsappMessage) (att whatsapp.WhatsappAttachment, err error) {
	msg, err := conn.FindMessage(imsg)
	if err != nil {
		return
	}

	data, err := conn.DownloadData(imsg)
	if err != nil {
		return
	}

	att = *msg.Attachment
	att.SetContent(&data)
	return
}

// Default SEND method using WhatsappMessage Interface
func (conn *WhatsrhymenConnection) Send(msg whatsapp.WhatsappMessage) (whatsapp.IWhatsappSendResponse, error) {

	if msg.Type == whatsapp.UnknownMessageType {
		if msg.HasAttachment() {
			msg.Type = MessageTypeFromAttachment(*msg.Attachment)
		}
	}

	switch msg.Type {
	case whatsapp.AudioMessageType:
		return conn.SendAudio(msg)
	case whatsapp.ImageMessageType:
		return conn.SendImage(msg)
	case whatsapp.DocumentMessageType:
		return conn.SendDocument(msg)
	case whatsapp.DiscardMessageType:
		// discarding
		return &whatsapp.WhatsappSendResponse{}, nil
	default:
		return conn.SendText(msg)
	}
}

func (conn *WhatsrhymenConnection) Disconnect() (err error) {
	conn.log.Info("disconnect requested")
	if conn.Client.GetConnected() {
		session, err := conn.Client.Disconnect()
		if err != nil {
			return err
		}

		conn.Session = &session
	}

	conn.failedToken = false
	return
}

func (conn *WhatsrhymenConnection) Delete() (err error) {
	// getting widclient
	wid, err := conn.GetWid()
	if err != nil {
		return
	}

	conn.Dispose()

	// deleting from store
	return WhatsrhymenService.Delete(wid)
}

func (conn *WhatsrhymenConnection) GetWhatsAppQRChannel(result chan<- string) (err error) {
	err = conn.EnsureUnderlying()
	if err != nil {
		return
	}

	session, err := conn.Client.Login(result)
	if err != nil {
		return
	}

	conn.log.Debug("GetWhatsAppQRChannel: login successful")

	// Se chegou até aqui é pq o QRCode foi validado e sincronizado
	// Saving session data
	err = WhatsrhymenService.UpdateSession(session)
	if err != nil {
		return
	}

	conn.Session = &session
	conn.Success()
	return
}

func (conn *WhatsrhymenConnection) UpdateLog(entry *log.Entry) {
	conn.log = entry
}

func (conn *WhatsrhymenConnection) UpdateHandler(handlers whatsapp.IWhatsappHandlers) {
	conn.WAHandlers = handlers
}

//endregion

func (conn *WhatsrhymenConnection) Failed(reason string) {
	conn.log.Debugf("updating failed token, reason: %s", reason)
	conn.failedToken = true

	if conn.Reconnect && !strings.Contains(reason, "invalid") {
		conn.log.Debugf("trying to auto reconnect")
		go conn.EnsureUnderlying()
	}
}

func (conn *WhatsrhymenConnection) Success() {
	conn.log.Debugf("updating success token")
	conn.failedToken = false
}

// Ensure a valid underlying whatsapp server connection
func (conn *WhatsrhymenConnection) EnsureUnderlying() (err error) {
	conn.syncConnection.Lock()

	if conn.Client == nil {
		connection, err := conn.GetWhatsAppClient()
		if err == nil {
			conn.Client = connection
			conn.log.Debugf("updating whatsrhymen connection")
		}
	} else if !conn.Client.GetConnected() {
		err = conn.Client.Restore()
		if err != nil {
			conn.Failed("restoring")
		}
	}

	conn.syncConnection.Unlock()
	return
}

func (conn *WhatsrhymenConnection) EnsureHandlers() (err error) {

	if conn.Handlers == nil {
		conn.Handlers = &WhatsrhymenHandlers{
			Connection: conn,
			log:        conn.log,
		}
	}

	if !conn.Handlers.IsRegistered() {
		err = conn.Handlers.Register()
		if err != nil {
			return err
		}
	}

	return
}

func (conn *WhatsrhymenConnection) UpdateClient() (err error) {
	err = conn.EnsureUnderlying()
	if err != nil {
		return
	}

	err = conn.EnsureHandlers()
	return
}

func (conn *WhatsrhymenConnection) GetWhatsAppClient() (client *whatsrhymen.Conn, err error) {
	client, err = whatsrhymen.NewConn(20 * time.Second)

	showing := whatsapp.WhatsappWebAppName + " Single"
	if len(whatsapp.WhatsappWebAppSystem) > 0 {
		showing += " " + whatsapp.WhatsappWebAppSystem
	}

	client.SetClientName(showing, whatsapp.WhatsappWebAppName, whatsapp.WhatsappWebAppVersion)
	client.SetClientVersion(2, 2208, 7)
	//client.SetClientVersion(2, 2142, 12)

	log.Debugf("debug client version :: %v", client.GetClientVersion())
	return
}

func (conn *WhatsrhymenConnection) Dispose() {
	// desabling auto reconnect
	conn.Reconnect = false

	_ = conn.DisposeEnsureLogOut()
	if conn.Client != nil {
		conn.Client.RemoveHandlers()
		if conn.Client.GetConnected() {
			if conn.Client.GetLoggedIn() {
				_ = conn.Client.Logout()
			} else {
				_, _ = conn.Client.Disconnect()
			}
		}
		conn.Client = nil
	}

	conn.Handlers = nil
	conn.Session = nil
	conn.WAHandlers = nil
	conn.log = nil
	conn = nil
}

func (conn *WhatsrhymenConnection) DisposeEnsureLogOut() (err error) {
	if conn.Session != nil {
		err = conn.EnsureUnderlying()
		if err != nil {
			return
		}

		conn.Client.RemoveHandlers()
		err = conn.Client.Logout()
		if err != nil {
			conn.log.Errorf("erro on trying to logout after a dispose connection: %s", err.Error())
		}
		return
	}
	return
}
