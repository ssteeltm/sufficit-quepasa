package models

import (
	whatsapp "github.com/sufficit/sufficit-quepasa-fork/whatsapp"
)

type QpReceiveResponse struct {
	QpResponse
	Messages []whatsapp.WhatsappMessage `json:"messages,omitempty"`
	Bot      QPBot                      `json:"bot"`
}
