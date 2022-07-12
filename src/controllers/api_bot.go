package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
	metrics "github.com/sufficit/sufficit-quepasa-fork/metrics"
	models "github.com/sufficit/sufficit-quepasa-fork/models"
	whatsapp "github.com/sufficit/sufficit-quepasa-fork/whatsapp"
)

// ReceiveAPIHandler renders route GET "/{version}/bot/{token}/receive"
func ReceiveAPIHandler(w http.ResponseWriter, r *http.Request) {

	// setting default reponse type as json
	w.Header().Set("Content-Type", "application/json")

	response := &models.QpReceiveResponse{}

	server, err := GetServer(w, r)
	if err != nil {
		metrics.MessageSendErrors.Inc()
		response.ParseError(err)
		RespondServerError(server, w, response)
		return
	}

	// Evitando tentativa de download de anexos sem o bot estar devidamente sincronizado
	status := server.GetStatus()
	if status != whatsapp.Ready {
		RespondNotReady(w, &ApiServerNotReadyException{Wid: server.GetWid(), Status: status})
		return
	}

	queryValues := r.URL.Query()
	paramTimestamp := queryValues.Get("timestamp")
	timestamp, err := GetTimestamp(paramTimestamp)
	if err != nil {
		metrics.MessageReceiveErrors.Inc()
		response.ParseError(err)
		RespondServerError(server, w, response)
		return
	}

	messages, err := GetMessages(server, timestamp)
	if err != nil {
		metrics.MessageReceiveErrors.Inc()
		response.ParseError(err)
		RespondServerError(server, w, response)
		return
	}

	metrics.MessagesReceived.Add(float64(len(messages)))

	response.Bot = *server.Bot
	response.Messages = messages

	if timestamp > 0 {
		response.ParseSuccess(fmt.Sprintf("getting with timestamp: %v", timestamp))
	} else {
		response.ParseSuccess("getting without filter")
	}
	RespondSuccess(w, response)
}

/*
<summary>
	Renders route POST "/{version}/bot/{token}/sendbinary/{chatId}/{fileName}/{textLabel}"

	Any of then, at this order of priority
	Path parameters: {chatId}
	Path parameters: {fileName}
	Path parameters: {textLabel} only images
	Url parameters: ?chatId={chatId}
	Url parameters: ?fileName={fileName}
	Url parameters: ?textLabel={textLabel} only images
	Header parameters: X-QUEPASA-CHATID = {chatId}
	Header parameters: X-QUEPASA-FILENAME = {fileName}
	Header parameters: X-QUEPASA-TEXTLABEL = {textLabel} only images
</summary>
*/
func SendDocumentFromBinary(w http.ResponseWriter, r *http.Request) {

	// setting default reponse type as json
	w.Header().Set("Content-Type", "application/json")

	response := &models.QpSendResponse{}

	server, err := GetServer(w, r)
	if err != nil {
		metrics.MessageSendErrors.Inc()
		response.ParseError(err)
		RespondServerError(server, w, response)
		return
	}

	// Declare a new request struct.
	var request models.QpSendRequest

	// Getting ChatId parameter
	chatId := chi.URLParam(r, "chatId")
	if len(chatId) == 0 && r.URL.Query().Has("chatId") {
		chatId = r.URL.Query().Get("chatId")
	} else if len(chatId) == 0 {
		chatId = r.Header.Get("X-QUEPASA-CHATID")
		if len(chatId) == 0 {
			metrics.MessageSendErrors.Inc()
			response.ParseError(fmt.Errorf("chat id missing"))
			RespondServerError(server, w, response)
			return
		}
	}

	request.ChatId = chatId

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		metrics.MessageSendErrors.Inc()
		response.ParseError(fmt.Errorf("attachment content missing or read error"))
		RespondServerError(server, w, response)
		return
	}

	request.Content = content

	// Getting FileName parameter
	fileName := chi.URLParam(r, "fileName")
	if len(fileName) == 0 && r.URL.Query().Has("fileName") {
		fileName = r.URL.Query().Get("fileName")
	} else if len(fileName) == 0 {
		fileName = r.Header.Get("X-QUEPASA-FILENAME")
	}

	// Setting filename
	request.FileName = fileName

	// Getting textLabel parameter
	textLabel := chi.URLParam(r, "textLabel")
	if len(textLabel) == 0 && r.URL.Query().Has("textLabel") {
		textLabel = r.URL.Query().Get("textLabel")
	} else if len(textLabel) == 0 {
		textLabel = r.Header.Get("X-QUEPASA-TEXTLABEL")
	}

	request.TextLabel = textLabel
	SendDocument(server, response, request, w)
}

/*
<summary>
	Renders route POST "/{version}/bot/{token}/sendencoded"

	Body parameter: {chatId}
	Body parameter: {fileName}
	Body parameter: {textLabel} only images
	Body parameter: {content}
</summary>
*/
func SendDocumentFromEncoded(w http.ResponseWriter, r *http.Request) {

	// setting default reponse type as json
	w.Header().Set("Content-Type", "application/json")

	response := &models.QpSendResponse{}

	server, err := GetServer(w, r)
	if err != nil {
		metrics.MessageSendErrors.Inc()
		response.ParseError(err)
		RespondServerError(server, w, response)
		return
	}

	// Declare a new request struct.
	var request models.QpSendRequestEncoded

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		metrics.MessageSendErrors.Inc()
		response.ParseError(err)
		RespondServerError(server, w, response)
		return
	}

	chatId, err := whatsapp.FormatEndpoint(request.ChatId)
	if err != nil {
		metrics.MessageSendErrors.Inc()
		response.ParseError(err)
		RespondServerError(server, w, response)
		return
	}

	// formatted chat id
	request.ChatId = chatId

	// base 64 content to byte array
	err = request.GenerateContent()
	if err != nil {
		metrics.MessageSendErrors.Inc()
		response.ParseError(err)
		RespondServerError(server, w, response)
		return
	}

	SendDocument(server, response, request.QpSendRequest, w)
}

/*
<summary>
	Renders route POST "/{version}/bot/{token}/sendurl"

	Body parameter: {url}
	Body parameter: {chatId}
	Body parameter: {fileName}
	Body parameter: {textLabel} only images
</summary>
*/
func SendDocumentFromUrl(w http.ResponseWriter, r *http.Request) {

	// setting default reponse type as json
	w.Header().Set("Content-Type", "application/json")

	response := &models.QpSendResponse{}

	server, err := GetServer(w, r)
	if err != nil {
		metrics.MessageSendErrors.Inc()
		response.ParseError(err)
		RespondServerError(server, w, response)
		return
	}

	// Declare a new request struct.
	var request models.QpSendRequestUrl

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		metrics.MessageSendErrors.Inc()
		response.ParseError(err)
		RespondServerError(server, w, response)
		return
	}

	chatId, err := whatsapp.FormatEndpoint(request.ChatId)
	if err != nil {
		metrics.MessageSendErrors.Inc()
		response.ParseError(err)
		RespondServerError(server, w, response)
		return
	}

	// formatted chat id
	request.ChatId = chatId

	// base 64 content to byte array
	err = request.GenerateContent()
	if err != nil {
		metrics.MessageSendErrors.Inc()
		response.ParseError(err)
		RespondServerError(server, w, response)
		return
	}

	SendDocument(server, response, request.QpSendRequest, w)
}

func SendDocument(server *models.QPWhatsappServer, response *models.QpSendResponse, request models.QpSendRequest, w http.ResponseWriter) {
	attach, err := request.ToWhatsappAttachment()
	if err != nil {
		metrics.MessageSendErrors.Inc()
		response.ParseError(err)
		RespondServerError(server, w, response)
		return
	}

	sendResponse, err := server.SendAttachment(request.ChatId, request.TextLabel, attach)
	if err != nil {
		metrics.MessageSendErrors.Inc()
		response.ParseError(err)
		RespondServerError(server, w, response)
		return
	}

	metrics.MessagesSent.Inc()
	response.ParseSuccess(sendResponse)
	RespondSuccess(w, response)
}
