package controllers

import (
	"fmt"
	"html/template"
	"net/http"

	. "github.com/sufficit/sufficit-quepasa-fork/metrics"
	. "github.com/sufficit/sufficit-quepasa-fork/models"
)

// FormReceiveController renders route GET "/bot/{botID}/receive"
func FormReceiveController(w http.ResponseWriter, r *http.Request) {
	data := QPFormReceiveData{PageTitle: "Receive", FormAccountEndpoint: FormAccountEndpoint}

	bot, err := GetBotFromRequest(r)
	if err != nil {
		data.ErrorMessage = err.Error()
	} else {
		data.Number = bot.GetNumber()
		data.Token = bot.Token
		data.DownloadPrefix = GetDownloadPrefix(bot.Token)
	}

	// Evitando tentativa de download de anexos sem o bot estar devidamente sincronizado
	if bot.GetStatus() != "ready" {
		RespondNotReady(w, fmt.Errorf("bot not ready yet ! try later."))
		return
	}

	queryValues := r.URL.Query()
	timestamp := queryValues.Get("timestamp")

	messages, err := RetrieveMessages(bot.ID, timestamp)
	if err != nil {
		MessageReceiveErrors.Inc()
		data.ErrorMessage = err.Error()
	}

	data.Messages = messages

	MessagesReceived.Add(float64(len(messages)))

	templates := template.Must(template.ParseFiles(
		"views/layouts/main.tmpl",
		"views/bot/receive.tmpl"))
	templates.ExecuteTemplate(w, "main", data)
}
