package models

// Requisição no formato QuePasa
// Utilizada na API do QuePasa para atualizar um WebHook de algum BOT
type QPWebhookRequest struct {
	Url string `json:"url"`
}
