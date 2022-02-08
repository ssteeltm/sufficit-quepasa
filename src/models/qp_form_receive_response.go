package models

type QPFormReceiveResponse struct {
	Messages []QPMessageV1 `json:"messages"`
	Bot      QPBot         `json:"bot"`
}
