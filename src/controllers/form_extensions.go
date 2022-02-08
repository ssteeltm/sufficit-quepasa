package controllers

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	. "github.com/sufficit/sufficit-quepasa-fork/models"
)

func RedirectToLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, FormLoginEndpoint, http.StatusFound)
}

// Google chrome bloqueou wss, portanto retornaremos sempre ws apatir de agora
func WebSocketProtocol() string {
	protocol := "ws"
	isSecure, _ := GetEnvBool("WEBSOCKETSSL", false)
	if isSecure {
		protocol = "wss"
	}

	return protocol
}

//
// Cycle
//

// CycleHandler renders route POST "/bot/cycle"
func CycleHandler(w http.ResponseWriter, r *http.Request) {
	user, err := GetUser(r)
	if err != nil {
		RedirectToLogin(w, r)
		return
	}

	r.ParseForm()
	botID := r.Form.Get("botID")
	bot, err := WhatsAppService.DB.Bot.FindForUser(user.ID, botID)
	if err != nil {
		return
	}

	err = bot.CycleToken()
	if err != nil {
		return
	}

	http.Redirect(w, r, FormAccountEndpoint, http.StatusFound)
}

// DebugHandler renders route POST "/bot/debug"
func DebugHandler(w http.ResponseWriter, r *http.Request) {
	user, err := GetUser(r)
	if err != nil {
		RedirectToLogin(w, r)
		return
	}

	r.ParseForm()
	botID := r.Form.Get("botID")
	bot, err := WhatsAppService.DB.Bot.FindForUser(user.ID, botID)
	if err != nil {
		return
	}

	err = bot.ToggleDevel()
	if err != nil {
		log.Printf("(%s)(ERR) Debug Handler :: '%s',", bot.GetNumber(), err)
		return
	}

	http.Redirect(w, r, FormAccountEndpoint, http.StatusFound)
}

// ToggleHandler renders route POST "/bot/toggle"
func ToggleHandler(w http.ResponseWriter, r *http.Request) {
	user, err := GetUser(r)
	if err != nil {
		RedirectToLogin(w, r)
		return
	}

	r.ParseForm()
	botID := r.Form.Get("botID")
	bot, err := WhatsAppService.DB.Bot.FindForUser(user.ID, botID)
	if err != nil {
		return
	}

	err = bot.Toggle()
	if err != nil {
		log.Print("(%s)(ERR) Toggle Handler :: '%s',", bot.GetNumber(), err)
		return
	}

	http.Redirect(w, r, FormAccountEndpoint, http.StatusFound)
}

//
// Verify
//

type verifyFormData struct {
	PageTitle    string
	ErrorMessage string
	Bot          QPBot
	Protocol     string
	Host         string
}

// VerifyFormHandler renders route GET "/bot/verify"
func VerifyFormHandler(w http.ResponseWriter, r *http.Request) {
	data := verifyFormData{
		PageTitle: "Verify To Add or Update",
		Protocol:  WebSocketProtocol(),
		Host:      r.Host,
	}

	templates := template.Must(template.ParseFiles(
		"views/layouts/main.tmpl",
		"views/bot/verify.tmpl",
	))
	templates.ExecuteTemplate(w, "main", data)
}

var done chan interface{}
var interrupt chan os.Signal
var upgrader = websocket.Upgrader{}

// VerifyHandler renders route GET "/bot/verify/ws"
func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	user, err := GetUser(r)
	if err != nil {
		log.Print("Connection upgrade error (not logged): ", err)
		RedirectToLogin(w, r)
		return
	}

	done = make(chan interface{})          // Channel to indicate that the receiverHandler is done
	interrupt = make(chan os.Signal)       // Channel to listen for interrupt signal to terminate gracefully
	signal.Notify(interrupt, os.Interrupt) // Notify the interrupt channel for SIGINT

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Connection upgrade error: ", err)
		return
	}

	defer conn.Close()
	go receiveWebSocketHandler(user, conn)

	// Our main loop for the client
	// We send our relevant packets here
	for {
		select {
		case <-time.After(time.Duration(1) * time.Millisecond * 1000):
			// Send an echo packet every second
			err := conn.WriteMessage(websocket.TextMessage, []byte("echo"))
			if err != nil {
				//log.Println("Error during writing to websocket:", err)
				return
			}

		case <-interrupt:
			// We received a SIGINT (Ctrl + C). Terminate gracefully...
			log.Println("Received SIGINT interrupt signal. Closing all pending connections")

			// Close our websocket connection
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Error during closing websocket:", err)
				return
			}

			select {
			case <-done:
				log.Println("Receiver Channel Closed! Exiting....")
			case <-time.After(time.Duration(1) * time.Second):
				log.Println("Timeout in closing receiving channel. Exiting....")
			}
			return
		}
	}
}

func receiveWebSocketHandler(user QPUser, connection *websocket.Conn) {
	defer close(done)
	for {
		mt, msg, err := connection.ReadMessage()
		if err != nil {
			log.Println("Error in receive:", err)
			return
		}

		if strings.EqualFold(string(msg), "start") {
			out := make(chan []byte)
			go func() {
				defer close(out)
				err = connection.WriteMessage(mt, <-out)
				if err != nil {
					log.Println("Write message error: ", err)
				}
			}()

			// Exibindo código QR
			bot, err := SignInWithQRCode(user, out)
			//panic("finalizando")
			if err != nil {
				if strings.Contains(err.Error(), "timed out") {
					err = connection.WriteMessage(mt, []byte("timeout"))
					if err != nil {
						// log.Println("Write message error after timeout: ", err)
						return
					}
				} else {
					log.Println("SignInWithQRCode Unknown Error:", err)
				}
			} else {
				log.Printf("(%s) SignInWithQRCode success ...", bot.GetNumber())

				// Marking as verified
				err = bot.MarkVerified(true)
				if err != nil {
					log.Printf("(%s)(ERR) Error on update verified state :: %s", bot.GetNumber(), err)
				}

				err = connection.WriteMessage(websocket.TextMessage, []byte("complete"))
				if err != nil {
					log.Printf("(%s)(ERR) Error on write complete message after qrcode verified :: %s", bot.GetNumber(), err)
				}

				go WhatsAppService.AppendNewServer(bot)
				return
			}
		} else {
			log.Printf("Received Unknown msg from WebSocket: %s\n", msg)
		}
	}
}

//
// Delete
//

// DeleteHandler renders route POST "/bot/{botID}/delete"
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	user, err := GetUser(r)
	if err != nil {
		return
	}

	r.ParseForm()
	botID := r.Form.Get("botID")

	bot, err := WhatsAppService.DB.Bot.FindForUser(user.ID, botID)
	if err != nil {
		return
	}

	if err := WhatsAppService.DB.Store.Delete(bot.ID); err != nil {
		return
	}

	if err := bot.Delete(); err != nil {
		return
	}

	http.Redirect(w, r, FormAccountEndpoint, http.StatusFound)
}

//
// Helpers
//

// Returns bot from http form request using E164 id
func GetBotFromRequest(r *http.Request) (QPBot, error) {
	var bot QPBot
	user, err := GetUser(r)
	if err != nil {
		return bot, err
	}

	botID := chi.URLParam(r, "id")

	return WhatsAppService.DB.Bot.FindForUser(user.ID, botID)
}

// Returns bot from http form request using E164 id
func GetServerFromRequest(r *http.Request) (*QPWhatsappServer, error) {
	wid := chi.URLParam(r, "id")
	return GetServerFromID(wid)
}
