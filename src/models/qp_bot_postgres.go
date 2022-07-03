package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type QPBotPostgres struct {
	db *sqlx.DB
}

func (source QPBotPostgres) FindAll() ([]*QPBot, error) {
	bots := []*QPBot{}
	err := source.db.Select(&bots, "SELECT * FROM bots")
	return bots, err
}

func (source QPBotPostgres) FindAllForUser(userID string) ([]QPBot, error) {
	bots := []QPBot{}
	err := source.db.Select(&bots, "SELECT * FROM bots WHERE user_id = $1", userID)
	return bots, err
}

func (source QPBotPostgres) FindByToken(token string) (QPBot, error) {
	var bot QPBot
	err := source.db.Get(&bot, "SELECT * FROM bots WHERE token = $1", token)
	return bot, err
}

func (source QPBotPostgres) FindForUser(userID string, ID string) (QPBot, error) {
	var bot QPBot
	err := source.db.Get(&bot, "SELECT * FROM bots WHERE user_id = $1 AND id = $2", userID, ID)
	return bot, err
}

func (source QPBotPostgres) FindByID(botID string) (QPBot, error) {
	var bot QPBot
	err := source.db.Get(&bot, "SELECT * FROM bots WHERE id = $1", botID)
	return bot, err
}

func (source QPBotPostgres) GetOrCreate(botID string, userID string) (bot QPBot, err error) {
	bot, err = source.FindByID(botID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			bot, err = source.Create(botID, userID)
		}
	}
	return
}

// botID = Wid of whatsapp connection
func (source QPBotPostgres) Create(botID string, userID string) (QPBot, error) {
	var bot QPBot
	token := uuid.New().String()
	now := time.Now()
	query := `INSERT INTO bots
    (id, is_verified, token, user_id, created_at, updated_at, devel)
    VALUES ($1, $2, $3, $4, $5, $6, $7)`
	if _, err := source.db.Exec(query, botID, false, token, userID, now, now, false); err != nil {
		return bot, err
	}

	return source.FindForUser(userID, botID)
}

//region SINGLE UPDATES

/*
UpdateToken(id string, value string) error
UpdateGroups(id string, value bool) error
UpdateBroadcast(id string, value bool) error
UpdateVerified(id string, value bool) error
UpdateDevel(id string, value bool) error
UpdateVersion(id string, value string) error
*/

func (source QPBotPostgres) UpdateToken(id string, value string) error {
	now := time.Now()
	query := "UPDATE bots SET token = $1, updated_at = $2 WHERE id = $3"
	_, err := source.db.Exec(query, value, now, id)
	return err
}

func (source QPBotPostgres) UpdateGroups(id string, value bool) error {
	now := time.Now()
	query := "UPDATE bots SET handlegroups = $1, updated_at = $2 WHERE id = $3"
	_, err := source.db.Exec(query, value, now, id)
	return err
}

func (source QPBotPostgres) UpdateBroadcast(id string, value bool) error {
	now := time.Now()
	query := "UPDATE bots SET handlebroadcast = $1, updated_at = $2 WHERE id = $3"
	_, err := source.db.Exec(query, value, now, id)
	return err
}

func (source QPBotPostgres) UpdateVerified(id string, value bool) error {
	now := time.Now()
	query := "UPDATE bots SET is_verified = $1, updated_at = $2 WHERE id = $3"
	_, err := source.db.Exec(query, value, now, id)
	return err
}

func (source QPBotPostgres) UpdateDevel(id string, value bool) error {
	now := time.Now()
	query := "UPDATE bots SET devel = $1, updated_at = $2 WHERE id = $3"
	_, err := source.db.Exec(query, value, now, id)
	return err
}

func (source QPBotPostgres) UpdateVersion(id string, value string) error {
	now := time.Now()
	query := "UPDATE bots SET version = $1, updated_at = $2 WHERE id = $3"
	_, err := source.db.Exec(query, value, now, id)
	return err
}

//endregion

func (source QPBotPostgres) Delete(id string) error {
	query := "DELETE FROM bots WHERE id = $1"
	_, err := source.db.Exec(query, id)
	return err
}
