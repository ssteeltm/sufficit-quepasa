package models

type QpResponse struct {
	Success bool   `json:"success"`
	Status  string `json:"status,omitempty"`
}

func (source QpResponse) Error() (message string) {
	if !source.Success {
		message = source.Status
	}
	return
}

func (source *QpResponse) ParseError(err error) {
	source.Success = false
	source.Status = err.Error()
}

func (source *QpResponse) ParseSuccess(status string) {
	source.Success = true
	source.Status = status
}