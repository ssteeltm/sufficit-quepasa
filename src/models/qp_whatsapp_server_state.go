package models

type QPWhatsappState uint64

const (
	Unknown QPWhatsappState = iota
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

func (s QPWhatsappState) String() string {
	switch s {
	case Created:
		return "Created"
	case Starting:
		return "Starting"
	case Stopped:
		return "Stopped"
	case Restarting:
		return "Restarting"
	case Disconnected:
		return "Disconnected"
	case Connected:
		return "Connected"
	case Fetching:
		return "Fetching"
	case Ready:
		return "Ready"
	case Halting:
		return "Halting"
	case Failed:
		return "Failed"
	}
	return "Unknown"
}
