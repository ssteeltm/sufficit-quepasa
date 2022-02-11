package models

// Usado apenas para compatibilidade de retorno das API antigas
type QPBotV1 struct {
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
