package controllers

import (
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
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

/*
<summary>
	Find a system track identifier to follow the message
	Getting from PATH => QUERY => HEADER
</summary>
*/
func GetTrackId(r *http.Request) (result string) {

	// retrieve from url path parameter
	result = chi.URLParam(r, "trackid")
	if len(result) == 0 {

		// retrieve from url query parameter
		if r.URL.Query().Has("trackid") {
			result = r.URL.Query().Get("trackid")
		} else {

			// retrieve from header parameter
			result = r.Header.Get("X-QUEPASA-TRACKID")
		}
	}
	return
}

// Getting PictureId from PATH => QUERY => HEADER
func GetPictureId(r *http.Request) (result string) {

	// retrieve from url path parameter
	result = chi.URLParam(r, "pictureid")
	if len(result) == 0 {

		// retrieve from url query parameter
		if r.URL.Query().Has("pictureid") {
			result = r.URL.Query().Get("pictureid")
		} else {

			// retrieve from header parameter
			result = r.Header.Get("X-QUEPASA-PICTUREID")
		}
	}
	return
}
