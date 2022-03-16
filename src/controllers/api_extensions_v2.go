package controllers

import (
	"sort"
	"strconv"
	"time"

	. "github.com/sufficit/sufficit-quepasa-fork/models"
)

// Retrieve messages with timestamp parameter
// Sorting then
func GetMessagesToAPIV2(server *QPWhatsappServer, timestamp string) (messages []QPMessageV2, err error) {

	searchTimestamp, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		if len(timestamp) > 0 {
			return
		} else {
			searchTimestamp = 0
		}
	}

	searchTime := time.Unix(searchTimestamp, 0)
	messages, err = GetMessagesFromServerV2(server, searchTime)
	if err != nil {
		return
	}

	sort.Sort(ByTimestampV2(messages))
	return
}
