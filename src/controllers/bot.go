package controllers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"os/signal"
    "time"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sufficit/sufficit-quepasa-fork/models"
)

//
// Metrics
//

var messagesSent = promauto.NewCounter(prometheus.CounterOpts{
	Name: "quepasa_sent_messages_total",
	Help: "Total sent messages",
})

var messageSendErrors = promauto.NewCounter(prometheus.CounterOpts{
	Name: "quepasa_send_message_errors_total",
	Help: "Total message send errors",
})

var messagesReceived = promauto.NewCounter(prometheus.CounterOpts{
	Name: "quepasa_received_messages_total",
	Help: "Total messages received",
})

var messageReceiveErrors = promauto.NewCounter(prometheus.CounterOpts{
	Name: "quepasa_receive_message_errors_total",
	Help: "Total message receive errors",
})

//
// Cycle
//

// CycleHandler renders route POST "/bot/cycle"
func CycleHandler(w http.ResponseWriter, r *http.Request) {
	user, err := models.GetUser(r)
	if err != nil {
		redirectToLogin(w, r)
		return
	}

	r.ParseForm()
	botID := r.Form.Get("botID")
	bot, err := models.WhatsAppService.DB.Bot.FindForUser(user.ID, botID)
	if err != nil {
		return
	}

	err = bot.CycleToken()
	if err != nil {
		return
	}

	http.Redirect(w, r, "/account", http.StatusFound)
}

// DebugHandler renders route POST "/bot/debug"
func DebugHandler(w http.ResponseWriter, r *http.Request) {
	user, err := models.GetUser(r)
	if err != nil {
		redirectToLogin(w, r)
		return
	}

	r.ParseForm()
	botID := r.Form.Get("botID")
	bot, err := models.WhatsAppService.DB.Bot.FindForUser(user.ID, botID)
	if err != nil {
		return
	}

	err = bot.ToggleDevel()
	if err != nil {
		log.Printf("(%s)(ERR) Debug Handler :: '%s',", bot.GetNumber(), err)
		return
	}

	http.Redirect(w, r, "/account", http.StatusFound)
}

// ToggleHandler renders route POST "/bot/toggle"
func ToggleHandler(w http.ResponseWriter, r *http.Request) {
	user, err := models.GetUser(r)
	if err != nil {
		redirectToLogin(w, r)
		return
	}

	r.ParseForm()
	botID := r.Form.Get("botID")
	bot, err := models.WhatsAppService.DB.Bot.FindForUser(user.ID, botID)
	if err != nil {
		return
	}

	err = bot.Toggle()
	if err != nil {
		log.Print("(%s)(ERR) Toggle Handler :: '%s',", bot.GetNumber(), err)
		return
	}

	http.Redirect(w, r, "/account", http.StatusFound)
}

//
// Verify
//

type verifyFormData struct {
	PageTitle    string
	ErrorMessage string
	Bot          models.QPBot
	Protocol     string
	Host         string
}

// VerifyFormHandler renders route GET "/bot/verify"
func VerifyFormHandler(w http.ResponseWriter, r *http.Request) {
	data := verifyFormData{
		PageTitle: "Verify To Add or Update",
		Protocol:  webSocketProtocol(),
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
	user, err := models.GetUser(r)
	if err != nil {
		log.Print("Connection upgrade error (not logged): ", err)
		redirectToLogin(w, r)
		return
	}

	done = make(chan interface{}) // Channel to indicate that the receiverHandler is done
    interrupt = make(chan os.Signal) // Channel to listen for interrupt signal to terminate gracefully 
    signal.Notify(interrupt, os.Interrupt) // Notify the interrupt channel for SIGINT

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Connection upgrade error: ", err)
		return
	}

	defer conn.Close()
	go receiveWebScoketHandler(user, conn)

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

func receiveWebScoketHandler(user models.QPUser,  connection *websocket.Conn) {
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
			bot, err := models.SignInWithQRCode(user, out)
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
				err = bot.MarkVerified(true)
				if err != nil {
					log.Printf("(%s)(ERR) Error on update verified state :: %s", bot.GetNumber(), err)
				}
		
				err = connection.WriteMessage(websocket.TextMessage, []byte("complete"))
				if err != nil {
					log.Printf("(%s)(ERR) Error on write complete message after qrcode verified :: %s", bot.GetNumber(), err)
				}
		
				go models.WhatsAppService.AppendNewServer(bot)
				return
			}
		}else{
			log.Printf("Received Unknown msg from WebSocket: %s\n", msg)
		}
    }
}

//
// Send
//

type sendFormData struct {
	PageTitle    string
	MessageId    string
	ErrorMessage string
	Bot          models.QPBot
}

func renderSendForm(w http.ResponseWriter, data sendFormData) {
	templates := template.Must(template.ParseFiles("views/layouts/main.tmpl", "views/bot/send.tmpl"))
	templates.ExecuteTemplate(w, "main", data)
}

// SendFormHandler renders route GET "/bot/{botID}/send"
func SendFormHandler(w http.ResponseWriter, r *http.Request) {
	data := sendFormData{
		PageTitle: "Send",
	}

	bot, err := findBot(r)
	if err != nil {
		data.ErrorMessage = err.Error()
		renderSendForm(w, data)
		return
	}

	data.Bot = bot
	renderSendForm(w, data)
}

// SendHandler renders route POST "/bot/{botID}/send"
// Vindo do formulário de testes
func SendHandler(w http.ResponseWriter, r *http.Request) {
	data := sendFormData{
		PageTitle: "Send",
	}
	bot, err := findBot(r)
	if err != nil {
		data.ErrorMessage = err.Error()
		renderSendForm(w, data)
		return
	}

	r.ParseForm()
	recipient := r.Form.Get("recipient")
	message := r.Form.Get("message")

	messageID, err := models.SendMessageFromBOT(bot.ID, recipient, message, models.QPAttachment{})
	if err != nil {
		messageSendErrors.Inc()
		data.ErrorMessage = err.Error()
		renderSendForm(w, data)
		return
	}

	data.MessageId = messageID

	messagesSent.Inc()

	renderSendForm(w, data)
}

//
// Receive
//

type receiveResponse struct {
	Messages []models.QPMessage `json:"messages"`
	Bot      models.QPBot       `json:"bot"`
}

type receiveFormData struct {
	PageTitle    string
	ErrorMessage string
	Number       string
	Messages     []models.QPMessage
}

// ReceiveFormHandler renders route GET "/bot/{botID}/receive"
func ReceiveFormHandler(w http.ResponseWriter, r *http.Request) {
	data := receiveFormData{
		PageTitle: "Receive",
	}

	bot, err := findBot(r)
	if err != nil {
		data.ErrorMessage = err.Error()
	} else {
		data.Number = bot.GetNumber()
	}

	// Evitando tentativa de download de anexos sem o bot estar devidamente sincronizado
	if bot.GetStatus() != "ready" {
		respondNotReady(w, fmt.Errorf("bot not ready yet ! try later."))
		return
	}

	queryValues := r.URL.Query()
	timestamp := queryValues.Get("timestamp")

	messages, err := models.RetrieveMessages(bot.ID, timestamp)
	if err != nil {
		messageReceiveErrors.Inc()
		data.ErrorMessage = err.Error()
	}

	data.Messages = messages

	messagesReceived.Add(float64(len(messages)))

	templates := template.Must(template.ParseFiles(
		"views/layouts/main.tmpl",
		"views/bot/receive.tmpl"))
	templates.ExecuteTemplate(w, "main", data)
}

//
// Delete
//

// DeleteHandler renders route POST "/bot/{botID}/delete"
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	user, err := models.GetUser(r)
	if err != nil {
		return
	}

	r.ParseForm()
	botID := r.Form.Get("botID")

	bot, err := models.WhatsAppService.DB.Bot.FindForUser(user.ID, botID)
	if err != nil {
		return
	}

	if err := models.WhatsAppService.DB.Store.Delete(bot.ID); err != nil {
		return
	}

	if err := bot.Delete(); err != nil {
		return
	}

	http.Redirect(w, r, "/account", http.StatusFound)
}

//
// Helpers
//

func findBot(r *http.Request) (models.QPBot, error) {
	var bot models.QPBot
	user, err := models.GetUser(r)
	if err != nil {
		return bot, err
	}

	botID := chi.URLParam(r, "botID")

	return models.WhatsAppService.DB.Bot.FindForUser(user.ID, botID)
}
