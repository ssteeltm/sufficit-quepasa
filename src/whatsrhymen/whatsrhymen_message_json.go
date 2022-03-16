package whatsrhymen

type WhatsrhymenMessageJson struct {
	Cmd WhatsrhymenMessageCmd `json:"cmd"`
}

type WhatsrhymenMessageCmd struct {
	Type string `json:"type"`
	Kind string `json:"kind"`
}
