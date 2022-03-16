package controllers

import (
	"sort"

	. "github.com/sufficit/sufficit-quepasa-fork/models"
)

// Retrieve messages with timestamp parameter
// Sorting then
func GetMessagesV1(source *QPBot, timestamp string) (messages []QPMessageV1, err error) {

	messages, err = GetMessagesFromBotV1(*source, timestamp)
	if err != nil {
		return
	}

	sort.Sort(ByTimestampV1(messages))
	return
}
