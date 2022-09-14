package models

type QpResponseInterface interface {
	IsSuccess() bool
	ParseSuccess(string)
	GetStatusMessage() string
}
