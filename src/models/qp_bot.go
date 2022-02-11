package models

import (
	"time"

	. "github.com/sufficit/sufficit-quepasa-fork/whatsapp"
)

type QPBot struct {
	ID        string `db:"id" json:"id"`
	Verified  bool   `db:"is_verified" json:"is_verified"`
	Token     string `db:"token" json:"token"`
	UserID    string `db:"user_id" json:"user_id"`
	WebHook   string `db:"webhook" json:"webhook,omitempty"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
	Devel     bool   `db:"devel" json:"devel"`
	Version   string `db:"version" json:"version,omitempty"`
}

// Traduz o Wid para um número de telefone em formato E164
func (bot *QPBot) GetNumber() string {
	phoneNumber, err := GetPhoneByID(bot.ID)
	if err != nil {
		return ""
	}
	return "+" + phoneNumber
}

func (bot *QPBot) GetStatus() string {
	server, err := GetServerFromBot(*bot)
	if err != nil {
		return Unknown.String()
	}

	return server.GetStatus().String()
}

func (bot *QPBot) GetTimestamp() (timestamp uint64) {
	server, err := GetServerFromBot(*bot)
	if err != nil {
		return
	}

	timestamp = uint64(server.Timestamp.Unix())
	return
}

func (bot *QPBot) GetStartedTime() (timestamp time.Time) {
	server, err := GetServerFromBot(*bot)
	if err != nil {
		return
	}

	timestamp = server.Timestamp
	return
}

func (bot *QPBot) GetBatteryInfo() (status WhatsAppBateryStatus) {
	server, err := GetServerFromBot(*bot)
	if err != nil {
		return
	}

	status = server.Battery
	return
}

func (bot *QPBot) Toggle() (err error) {
	server, err := GetServerFromBot(*bot)
	if err != nil {
		go WhatsAppService.AppendNewServer(bot)
	} else {
		if server.GetStatus() == Stopped || server.GetStatus() == Created {
			err = server.Start()
		} else {
			err = server.Shutdown()
		}
	}
	return
}

func (bot *QPBot) IsDevelopmentGlobal() bool {
	return ENV.IsDevelopment()
}

func (bot *QPBot) MarkVerified(ok bool) error {
	return WhatsAppService.DB.Bot.MarkVerified(bot.ID, ok)
}

func (bot *QPBot) CycleToken() error {
	return WhatsAppService.DB.Bot.CycleToken(bot.ID)
}

func (bot *QPBot) Delete() error {
	return WhatsAppService.DB.Bot.Delete(bot.ID)
}

func (bot *QPBot) WebHookUpdate() error {
	return WhatsAppService.DB.Bot.WebHookUpdate(bot.WebHook, bot.ID)
}

func (bot *QPBot) WebHookSincronize() error {
	webhook, err := WhatsAppService.DB.Bot.WebHookSincronize(bot.ID)
	bot.WebHook = webhook
	return err
}

func (bot *QPBot) ToggleDevel() (err error) {
	if bot.Devel {
		err = WhatsAppService.DB.Bot.Devel(bot.ID, false)
		bot.Devel = false
	} else {
		err = WhatsAppService.DB.Bot.Devel(bot.ID, true)
		bot.Devel = true
	}
	return err
}

func (bot *QPBot) Send(msg *WhatsappMessage) (err error) {
	return SendMessageFromBot(bot, msg)
}
