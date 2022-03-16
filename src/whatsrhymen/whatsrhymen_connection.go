package whatsrhymen

import (
	"fmt"

	whatsrhymen "github.com/Rhymen/go-whatsapp"
	log "github.com/sirupsen/logrus"
	whatsapp "github.com/sufficit/sufficit-quepasa-fork/whatsapp"
)

// Must Implement IWhatsappConnection
type WhatsrhymenConnection struct {
	Client   *whatsrhymen.Conn
	Handlers *WhatsrhymenHandlers
	Session  *whatsrhymen.Session

	logger *log.Logger
	log    *log.Entry
}

//wac.AddHandler(WhatsrhymenHandlers{})
//wac, err := whatsapp.NewConn(20 * time.Second)
//region IMPLEMENT INTERFACE WHATSAPP CONNECTION

func (conn *WhatsrhymenConnection) GetVersion() string { return "single" }

func (conn *WhatsrhymenConnection) GetWid() (wid string, err error) {
	if conn.Client == nil {
		err = fmt.Errorf("client not defined on trying to get wid")
	} else {
		if conn.Client.Info == nil {
			if conn.Session != nil {
				wid = conn.Session.Wid
				return
			}
			err = fmt.Errorf("session|&|client info not defined on trying to get wid")
		} else {
			wid = conn.Client.Info.Wid
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
			if conn.Client.Info != nil {
				if conn.Client.Info.Connected {
					state = whatsapp.Connected
					if conn.Handlers != nil {
						state = whatsapp.Ready
					}
				} else {
					state = whatsapp.Disconnected
				}
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
	conn.log.Info("starting whatsrhymen connecting ...")

	// Agora sim, restaura a conexão com o whatsapp apartir de uma seção salva
	_, err = conn.Client.RestoreWithSession(*conn.Session)
	if err != nil {
		conn.log.Printf("(ERR) Error on restore session :: %s", err)
		return
	}

	return
}

// func (cli *Client) Download(msg DownloadableMessage) (data []byte, err error)
func (conn *WhatsrhymenConnection) Download(imsg whatsapp.IWhatsappMessage) (data []byte, err error) {
	conn.log.Info(imsg.GetSource())
	return
}

// Default SEND method using WhatsappMessage Interface
func (conn *WhatsrhymenConnection) Send(msg whatsapp.WhatsappMessage) (whatsapp.IWhatsappSendResponse, error) {
	conn.log.Debug("sending message")
	return nil, nil
}

// func (cli *Client) Upload(ctx context.Context, plaintext []byte, appInfo MediaType) (resp UploadResponse, err error)
func (conn *WhatsrhymenConnection) UploadAttachment(msg whatsapp.WhatsappMessage) (err error) {
	return
}

func (conn *WhatsrhymenConnection) Disconnect() (err error) {
	session, err := conn.Client.Disconnect()
	if err != nil {
		return
	}

	conn.Session = &session
	return
}

func (conn *WhatsrhymenConnection) Delete() error {
	return nil
}

func (conn *WhatsrhymenConnection) GetWhatsAppQRChannel(result chan<- string) (err error) {
	session, err := conn.Client.Login(result)
	if err != nil {
		return
	}

	log.Printf("login successful, session")
	conn.Session = &session

	// Se chegou até aqui é pq o QRCode foi validado e sincronizado
	// Saving session data
	err = WhatsrhymenService.Container.Update(session)
	return
}

func (conn *WhatsrhymenConnection) UpdateHandler(handlers whatsapp.IWhatsappHandlers) {
	conn.Handlers.WAHandlers = handlers
}

//endregion

func (conn *WhatsrhymenConnection) LogLevel(level log.Level) {

}

func (conn *WhatsrhymenConnection) PrintStatus() {
	/*
		conn.log.Warnf("STATUS IS CONNECTED: %v", conn.Client.IsConnected())
		conn.log.Warnf("STATUS IS LOGGED IN: %v", conn.Client.IsLoggedIn())

		conn.Client.SendPresence(types.PresenceAvailable)
		conn.Client.SetPassive(false)
	*/
}
