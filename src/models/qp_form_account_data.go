package models

type QPFormAccountData struct {
	PageTitle    string
	ErrorMessage string
	Bots         []QPBot
	User         QPUser
}