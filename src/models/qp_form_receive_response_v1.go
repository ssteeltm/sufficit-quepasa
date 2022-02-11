package models

type QPFormReceiveResponseV1 struct {
	Messages []QPMessageV1 `json:"messages"`
	Bot      QPBot         `json:"bot"`
}
