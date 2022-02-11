package models

import (
	. "github.com/sufficit/sufficit-quepasa-fork/whatsapp"
)

type QPFormReceiveResponse struct {
	Messages []WhatsappMessage `json:"messages"`
	Bot      QPBot             `json:"bot"`
}
