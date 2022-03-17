package whatsapp

type WhatsappConnectionState uint

const (
	// Unknown, not treated state
	Unknown WhatsappConnectionState = iota

	// Instantiated
	Created

	// Starting variables
	Starting

	// Connecting to whatsapp servers
	Connecting

	// Stopped requested
	Stopped

	Restarting

	// Connected to whatsapp servers
	Connected

	// Fetching messages from servers
	Fetching

	// Ready and Fully operational
	Ready

	// Finalizing
	Halting

	// Disconnected from whatsapp servers
	Disconnected

	// Failed state, for any reason
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
		"Connecting",
		"Stopped",
		"Restarting",
		"Connected",
		"Fetching",
		"Ready",
		"Halting",
		"Disconnected",
		"Failed",
	}[s]
}
