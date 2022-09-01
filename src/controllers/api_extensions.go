package controllers

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	models "github.com/sufficit/sufficit-quepasa/models"
	whatsapp "github.com/sufficit/sufficit-quepasa/whatsapp"
)

func GetTimestamp(timestamp string) (result int64, err error) {
	if len(timestamp) > 0 {
		result, err = strconv.ParseInt(timestamp, 10, 64)
		if err != nil {
			if len(timestamp) > 0 {
				return
			} else {
				result = 0
			}
		}
	}
	return
}

// Retrieve messages with timestamp parameter
// Sorting then
func GetMessages(server *models.QPWhatsappServer, timestamp int64) (messages []whatsapp.WhatsappMessage) {
	searchTime := time.Unix(timestamp, 0)
	messages = server.GetMessages(searchTime)
	sort.Sort(whatsapp.ByTimestamp(messages))
	return
}

// Getting chatId from PATH => QUERY => HEADER
func GetChatId(r *http.Request) (chatId string) {

	// retrieve from url path parameter
	chatId = chi.URLParam(r, "chatid")
	if len(chatId) == 0 {

		// retrieve from url query parameter
		if r.URL.Query().Has("chatid") {
			chatId = r.URL.Query().Get("chatid")
		} else {

			// retrieve from header parameter
			chatId = r.Header.Get("X-QUEPASA-CHATID")
		}
	}
	return
}

func EnsureChatId(sendR *models.QpSendRequest, r *http.Request) (err error) {
	if len(sendR.ChatId) == 0 {
		sendR.ChatId = GetChatId(r)
	}

	if len(sendR.ChatId) == 0 {
		err = fmt.Errorf("chat id missing")
	}
	return
}
