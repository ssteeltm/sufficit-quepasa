package models

type IQPBot interface {
	FindAll() ([]*QPBot, error)
	FindAllForUser(userID string) ([]QPBot, error)
	FindByToken(token string) (QPBot, error)
	FindForUser(userID string, ID string) (QPBot, error)
	FindByID(botID string) (QPBot, error)
	GetOrCreate(botID string, userID string) (bot QPBot, err error)
	Create(botID string, userID string) (QPBot, error)

	/// FORWARDING ---
	UpdateToken(id string, value string) error
	UpdateGroups(id string, value bool) error
	UpdateBroadcast(id string, value bool) error
	UpdateVerified(id string, value bool) error
	UpdateWebhook(id string, value string) error
	UpdateDevel(id string, value bool) error
	UpdateVersion(id string, value string) error

	Delete(id string) error
	WebHookSincronize(id string) (result string, err error)
}
