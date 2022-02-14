package library

import (
	"reflect"
	"regexp"
)

// Validate email string
func IsValidEMail(s string) bool {
	var rx = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if len(s) < 255 && rx.MatchString(s) {
		return true
	}

	return false
}

// Returns a string representation of source type interface
func GetTypeString(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}
