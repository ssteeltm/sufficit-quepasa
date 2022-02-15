package whatsapp

type WhatsappConnectionState uint

const (
	Unknown WhatsappConnectionState = iota
	Created
	Starting
	Stopped
	Restarting
	Disconnected
	Connected
	Fetching
	Ready
	Halting
	Failed
)

// EnumIndex - Creating common behavior - give the type a EnumIndex function
func (s WhatsappConnectionState) EnumIndex() int {
	return int(s)
}

func (s WhatsappConnectionState) String() string {
	return [...]string{
		"Unknown",
		"Created",
		"Starting",
		"Stopped",
		"Restarting",
		"Disconnected",
		"Connected",
		"Fetching",
		"Ready",
		"Halting",
		"Failed",
	}[s]
}
