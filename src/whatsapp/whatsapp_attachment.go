package whatsapp

type WhatsappAttachment struct {
	content *[]byte `json:"-"`

	Mimetype   string `json:"mime"`
	FileLength uint64 `json:"filelength"`

	// document
	FileName string `json:"filename,omitempty"`

	// video | image
	JpegThumbnail string `json:"thumbnail,omitempty"`

	// audio
	Seconds uint32 `json:"seconds,omitempty"`
}

func (source *WhatsappAttachment) GetContent() *[]byte {
	return source.content
}

func (source *WhatsappAttachment) SetContent(content *[]byte) {
	source.content = content
}
