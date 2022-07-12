package models

type QpInfoResponse struct {
	QpResponse
	Server QPWhatsappServer `json:"server,omitempty"`
}

func (source *QpInfoResponse) ParseSuccess(server QPWhatsappServer) {
	source.QpResponse.ParseSuccess("follow server information")
	source.Server = server
}
