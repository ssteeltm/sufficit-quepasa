package models

import (
	"time"
)

type QPBot struct {
	ID              string `db:"id" json:"id"`
	Verified        bool   `db:"is_verified" json:"is_verified"`
	Token           string `db:"token" json:"token"`
	UserID          string `db:"user_id" json:"user_id"`
	Webhook         string `db:"webhook" json:"webhook,omitempty"`
	CreatedAt       string `db:"created_at" json:"created_at"`
	UpdatedAt       string `db:"updated_at" json:"updated_at"`
	Devel           bool   `db:"devel" json:"devel"`
	Version         string `db:"version" json:"version,omitempty"`
	HandleGroups    bool   `db:"handlegroups" json:"handlegroups,omitempty"`
	HandleBroadcast bool   `db:"handlebroadcast" json:"handlebroadcast,omitempty"`

	db IQPBot
}

//region DATABASE METHODS

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
		go WhatsappService.AppendNewServer(bot)
	} else {
		if server.GetStatus() == Stopped || server.GetStatus() == Created {
			err = server.Start()
		} else {
			err = server.Shutdown()
		}
	}
	return
}

//region SINGLE UPDATES

func (bot *QPBot) UpdateVersion(value string) (err error) {
	err = bot.db.UpdateVersion(bot.ID, value)
	if err != nil {
		return
	}

	bot.Version = value
	return
}

func (bot *QPBot) UpdateGroups(value bool) (err error) {
	err = bot.db.UpdateGroups(bot.ID, value)
	if err != nil {
		return
	}

	bot.HandleGroups = value
	return
}

func (bot *QPBot) UpdateBroadcast(value bool) (err error) {
	err = bot.db.UpdateBroadcast(bot.ID, value)
	if err != nil {
		return
	}

	bot.HandleBroadcast = value
	return
}

func (bot *QPBot) UpdateVerified(value bool) (err error) {
	err = bot.db.UpdateVerified(bot.ID, value)
	if err != nil {
		return
	}

	bot.Verified = value
	return
}

func (bot *QPBot) UpdateDevel(value bool) (err error) {
	err = bot.db.UpdateDevel(bot.ID, value)
	if err != nil {
		return
	}

	bot.Devel = value
	return
}

func (bot *QPBot) UpdateToken(value string) (err error) {
	err = bot.db.UpdateToken(bot.ID, value)
	if err != nil {
		return
	}

	bot.Token = value
	return
}

func (bot *QPBot) UpdateWebhook(value string) (err error) {
	err = bot.db.UpdateWebhook(bot.ID, value)
	if err != nil {
		return
	}

	bot.Webhook = value
	return
}

//endregion

func (bot *QPBot) IsDevelopmentGlobal() bool {
	return ENV.IsDevelopment()
}

func (bot *QPBot) Delete() error {
	return bot.db.Delete(bot.ID)
}

func (bot *QPBot) WebHookSincronize() error {
	webhook, err := bot.db.WebHookSincronize(bot.ID)
	bot.Webhook = webhook
	return err
}

//endregion
