package library

import(
	"fmt"
	"regexp"
	"strings"
)

var AllowedSuffix = map[string]bool{
	"g.us":           true, // Mensagem para um grupo
	"s.whatsapp.net": true, // Mensagem direta a um usuário
}

// Formata um texto qualquer em formato de destino válido para o sistema do whatsapp
func FormatEndpoint(source string) (recipient string, err error){
	
	// removing whitespaces
	recipient = strings.Replace(source, " ", "", -1)

	// removing starting + from E164 phones
	recipient = strings.TrimLeft(recipient, "+")

	if strings.ContainsAny(recipient, "@") {
		splited := strings.Split(recipient, "@")
		if !AllowedSuffix[splited[1]] {
			err = fmt.Errorf("invalid recipient %s", recipient)
			return
		}

		recipient = splited[0]
		if strings.Contains(recipient, ".") {

			// capturando tudo antes do "."
			splited =  strings.Split(recipient, ".")
			recipient = splited[0]
			return 
		}
	} else {	
		if strings.Contains(recipient, "-"){
			splited := strings.Split(recipient, "-")
			if !IsValidE164(splited[0]) {
				err = fmt.Errorf("contains - but its not a valid group: %s", source)
				return
			}
		} else {
			if !IsValidE164(recipient) {
				err = fmt.Errorf("its not a valid e164 phone: %s", source)
				return
			}
		}
	}

	return
}

var RegexValidE164Test = string(`\d`)
func IsValidE164(phone string) bool {
	regex, err := regexp.Compile(RegexValidE164Test)
	if err != nil { panic("invalid regex on IsValidE164 :: " + RegexValidE164Test) }
	matches := regex.FindAllString(phone, -1)
	if len(matches) >= 9 && len(matches) <= 15 {
		return true
	}
	return false
}