package controllers

import (
	"sort"
	"strconv"
	"time"

	models "github.com/sufficit/sufficit-quepasa-fork/models"
	whatsapp "github.com/sufficit/sufficit-quepasa-fork/whatsapp"
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
