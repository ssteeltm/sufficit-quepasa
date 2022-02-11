package controllers

import (
	"sort"
	"strconv"
	"time"

	. "github.com/sufficit/sufficit-quepasa-fork/models"
	. "github.com/sufficit/sufficit-quepasa-fork/whatsapp"
)

// Retrieve messages with timestamp parameter
// Sorting then
func GetMessages(server *QPWhatsappServer, timestamp string) (messages []WhatsappMessage, err error) {

	searchTimestamp, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		if len(timestamp) > 0 {
			return
		} else {
			searchTimestamp = 0
		}
	}

	searchTime := time.Unix(searchTimestamp, 0)
	messages, err = server.GetMessages(searchTime)
	if err != nil {
		return
	}

	sort.Sort(ByTimestamp(messages))
	return
}
