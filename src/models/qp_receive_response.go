package models

import (
	whatsapp "github.com/sufficit/sufficit-quepasa-fork/whatsapp"
)

type QpReceiveResponse struct {
	QpResponse
	Total    uint64                     `json:"total"`
	Messages []whatsapp.WhatsappMessage `json:"messages,omitempty"`
	Bot      QPBot                      `json:"bot"`
}
