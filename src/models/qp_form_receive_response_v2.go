package models

type QPFormReceiveResponseV2 struct {
	Messages []QPMessageV2 `json:"messages"`
	Bot      QPBotV1       `json:"bot"`
}
