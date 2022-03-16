package whatsapp

type State int64

const (
	// since iota starts with 0, the first value
	// defined here will be the default
	Undefined State = iota
	Verified
	Autumn
	Winter
	Spring
)