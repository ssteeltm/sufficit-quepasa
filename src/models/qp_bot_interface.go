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
	MarkVerified(id string, ok bool) error
	CycleToken(id string) error
	Delete(id string) error
	WebHookUpdate(webhook string, id string) error
	WebHookSincronize(id string) (result string, err error)
	Devel(id string, status bool) error

	SetVersion(id string, version string) (err error)
}
