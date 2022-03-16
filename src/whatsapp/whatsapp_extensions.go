package whatsapp

import (
	"fmt"
	"regexp"
	"strings"
)

var AllowedSuffix = map[string]bool{
	"g.us":           true, // Mensagem para um grupo
	"s.whatsapp.net": true, // Mensagem direta a um usuário
}

func PhoneToWid(source string) (destination string) {

	// removing starting + from E164 phones
	destination = strings.TrimLeft(source, "+")

	if !strings.ContainsAny(destination, "@") {
		return destination + "@s.whatsapp.net"
	}
	return
}

// Formata um texto qualquer em formato de destino válido para o sistema do whatsapp
func FormatEndpoint(source string) (destination string, err error) {

	// removing whitespaces
	destination = strings.Replace(source, " ", "", -1)

	// if have a + as prefix, is a phone number
	if strings.HasPrefix(source, "+") {
		destination = PhoneToWid(source)
		return
	}

	if strings.ContainsAny(destination, "@") {
		splited := strings.Split(destination, "@")
		if !AllowedSuffix[splited[1]] {
			err = fmt.Errorf("invalid recipient %s", destination)
			return
		}

		return
	} else {
		if strings.Contains(destination, "-") {
			splited := strings.Split(destination, "-")
			if !IsValidE164(splited[0]) {
				err = fmt.Errorf("contains - but its not a valid group: %s", source)
				return
			}
		} else {
			if IsValidE164(destination) {
				destination = PhoneToWid(destination)
				return
			} else {
				destination = destination + "@g.us"
				return
			}
		}
	}

	return
}

var RegexValidE164Test = string(`\d`)

func IsValidE164(phone string) bool {
	regex, err := regexp.Compile(RegexValidE164Test)
	if err != nil {
		panic("invalid regex on IsValidE164 :: " + RegexValidE164Test)
	}
	matches := regex.FindAllString(phone, -1)
	if len(matches) >= 9 && len(matches) <= 15 {
		return true
	}
	return false
}
