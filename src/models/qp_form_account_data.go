package models

type QPFormAccountData struct {
	PageTitle    string
	ErrorMessage string
	Servers      map[string]*QPWhatsappServer
	User         QPUser
}
