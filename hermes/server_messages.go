package main

/*
Server supported message types are:
	- "confirmation"
	- "error"
*/
const (
	confMessageType  = "confirmation"
	errorMessageType = "error"
)

type ConfirmationMessage struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func newConfirmationMessage(msg string) ConfirmationMessage {
	return ConfirmationMessage{confMessageType, msg}
}
