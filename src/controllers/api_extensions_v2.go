package controllers

import (
	"sort"

	. "github.com/sufficit/sufficit-quepasa-fork/models"
)

// Retrieve messages with timestamp parameter
// Sorting then
func GetMessagesV2(source *QPBot, timestamp string) (messages []QPMessageV2, err error) {

	messages, err = GetMessagesFromBotV2(*source, timestamp)
	if err != nil {
		return
	}

	sort.Sort(ByTimestampV2(messages))
	return
}
