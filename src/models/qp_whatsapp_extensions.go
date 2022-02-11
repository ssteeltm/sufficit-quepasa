package models

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/jinzhu/copier"
	log "github.com/sirupsen/logrus"
	"github.com/skip2/go-qrcode"

	. "github.com/sufficit/sufficit-quepasa-fork/library"
	. "github.com/sufficit/sufficit-quepasa-fork/whatsapp"
)

func TryUpdateHttpChannel(ch chan<- []byte, value []byte) (closed bool) {
	defer func() {
		if recover() != nil {
			// the return result can be altered
			// in a defer function call
			closed = false
		}
	}()

	ch <- value // panic if ch is closed
	return true // <=> closed = false; return
}

// Envia o QRCode para o usuário e aguarda pela resposta
// Retorna um novo BOT
func SignInWithQRCode(user QPUser, out chan<- []byte) (bot *QPBot, err error) {
	con, err := NewConnection("")
	if err != nil {
		return
	}

	log.Info("GetWhatsAppQRChannel ...")
	qrChan := make(chan string)
	go con.GetWhatsAppQRChannel(qrChan)
	for qrBase64 := range qrChan {
		var png []byte
		png, err := qrcode.Encode(qrBase64, qrcode.Medium, 256)
		if err != nil {
			log.Printf("(ERR) Error on QrCode encode :: %v\r", err.Error())
		}
		encodedPNG := base64.StdEncoding.EncodeToString(png)

		if !TryUpdateHttpChannel(out, []byte(encodedPNG)) {
			break
		}
	}

	wid, err := con.GetWid()
	if err != nil {
		return
	}

	if len(wid) == 0 {
		err = fmt.Errorf("invalid wid !")
		return
	}

	// Se chegou até aqui é pq o QRCode foi validado e sincronizado
	databaseBot, err := WhatsAppService.DB.Bot.GetOrCreate(wid, user.ID)
	if err != nil {
		log.Printf("(ERR) Error on get or create bot after login :: %v\r", err.Error())
		return
	}

	bot = &databaseBot

	// Updating connection version information
	bot.Version = con.GetVersion()
	err = WhatsAppService.DB.Bot.SetVersion(bot.ID, bot.Version)

	return
}

func GetDownloadPrefixFromWid(wid string) (path string, err error) {
	server, ok := WhatsAppService.Servers[wid]
	if !ok {
		err = fmt.Errorf("server not found: %s", wid)
		return
	}

	prefix := fmt.Sprintf("/bot/%s/download", server.Bot.Token)
	return prefix, err
}

func ToQPAttachment(source *WhatsappAttachment, id string, wid string) (attach *QPAttachment) {

	// Anexo que devolverá ao utilizador da api, cliente final
	// com Url pública válida sem criptografia
	attach = &QPAttachment{}

	copier.Copy(attach, source)
	url, err := GetDownloadPrefixFromWid(wid)
	if err != nil {
		return
	}

	attach.DirectPath = url + "/" + id
	return
}

func ToQPEndPoint(source WhatsappEndpoint) (endpoint QPEndPoint) {
	WhatsappID := source.ID
	if !strings.Contains(WhatsappID, "@") {
		if strings.Contains(WhatsappID, "-") {
			WhatsappID = WhatsappID + "@g.us"
		} else {
			WhatsappID = WhatsappID + "@s.whatsapp.net"
		}
	}

	endpoint.ID = WhatsappID
	endpoint.Title = source.Title
	if len(endpoint.Title) == 0 {
		endpoint.Title = source.UserName
	}

	return
}

func ChatToQPEndPoint(source WhatsappChat) (endpoint QPEndPoint) {
	WhatsappID := source.ID
	if !strings.Contains(WhatsappID, "@") {
		if strings.Contains(WhatsappID, "-") {
			WhatsappID = WhatsappID + "@g.us"
		} else {
			WhatsappID = WhatsappID + "@s.whatsapp.net"
		}
	}

	endpoint.ID = WhatsappID
	endpoint.Title = source.Title
	return
}

func ToWhatsappMessage(destination string, text string, attach *WhatsappAttachment) (msg *WhatsappMessage, err error) {
	recipient, err := FormatEndpoint(destination)
	if err != nil {
		return
	}

	chat := WhatsappChat{ID: recipient}
	msg = &WhatsappMessage{}
	msg.Text = text
	msg.Chat = chat
	if attach != nil {
		msg.Attachment = attach
	}
	return

}
