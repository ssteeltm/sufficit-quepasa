package models

import (
	"errors"
	"time"
)

type QPBotV2 struct {
	ID              string `db:"id" json:"id"`
	Verified        bool   `db:"is_verified" json:"is_verified"`
	Token           string `db:"token" json:"token"`
	UserID          string `db:"user_id" json:"user_id"`
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
func (bot *QPBotV2) GetNumber() string {
	phoneNumber, err := GetPhoneByID(bot.ID)
	if err != nil {
		return ""
	}
	return "+" + phoneNumber
}

func (bot *QPBotV2) GetTimestamp() (timestamp uint64) {
	server, err := GetServerFromID(*&bot.ID)
	if err != nil {
		return
	}

	timestamp = uint64(server.Timestamp.Unix())
	return
}

func (bot *QPBotV2) GetStartedTime() (timestamp time.Time) {
	server, err := GetServerFromID(*&bot.ID)
	if err != nil {
		return
	}

	timestamp = server.Timestamp
	return
}

func (bot *QPBotV2) GetBatteryInfo() (status WhatsAppBateryStatus) {
	server, err := GetServerFromID(*&bot.ID)
	if err != nil {
		return
	}

	status = server.Battery
	return
}

//region SINGLE UPDATES

func (bot *QPBotV2) UpdateVersion(value string) (err error) {
	err = bot.db.UpdateVersion(bot.ID, value)
	if err != nil {
		return
	}

	bot.Version = value
	return
}

func (bot *QPBotV2) UpdateGroups(value bool) (err error) {
	err = bot.db.UpdateGroups(bot.ID, value)
	if err != nil {
		return
	}

	bot.HandleGroups = value
	return
}

func (bot *QPBotV2) UpdateBroadcast(value bool) (err error) {
	err = bot.db.UpdateBroadcast(bot.ID, value)
	if err != nil {
		return
	}

	bot.HandleBroadcast = value
	return
}

func (bot *QPBotV2) UpdateVerified(value bool) (err error) {
	err = bot.db.UpdateVerified(bot.ID, value)
	if err != nil {
		return
	}

	bot.Verified = value
	return
}

func (bot *QPBotV2) UpdateDevel(value bool) (err error) {
	err = bot.db.UpdateDevel(bot.ID, value)
	if err != nil {
		return
	}

	bot.Devel = value
	return
}

func (bot *QPBotV2) UpdateToken(value string) (err error) {
	err = bot.db.UpdateToken(bot.ID, value)
	if err != nil {
		return
	}

	bot.Token = value
	return
}

func (bot *QPBotV2) UpdateWebhook(value string) (err error) {
	return errors.New("not supported anymore")
}

//endregion

func (bot *QPBotV2) IsDevelopmentGlobal() bool {
	return ENV.IsDevelopment()
}

func (bot *QPBotV2) Delete() error {
	return bot.db.Delete(bot.ID)
}

func (bot *QPBotV2) WebHookSincronize() error {
	return errors.New("not supported anymore")
}

//endregion
