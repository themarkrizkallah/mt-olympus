package main

import (
	"errors"
	"fmt"
	"strings"
)

var GenericError error = errors.New("bad request")

type MessageTypeError struct {
	givenType string
}

func (mte MessageTypeError) Error() string {
	return fmt.Sprintf(
		"Accepted message types are [%s], user passed in: %s",
		strings.Join(acceptedMsgTypes[:], ","),
		mte.givenType)
}

type ChannelNameError struct {
	givenName string
}

func (cne ChannelNameError) Error() string {
	return fmt.Sprintf(
		"Accepted channel names are [%s], user passed in: %s",
		strings.Join(acceptedChanNames[:], ","),
		cne.givenName)
}

type ProductIDError struct {
	givenID string
}

func (pie ProductIDError) Error() string {
	return fmt.Sprintf(
		"Accepted product_ids are [%s], user passed in: %s",
		strings.Join(acceptedProductIDs[:], ","),
		pie.givenID)
}

// Return true if e is an error defined in this file, false otherwise
func definedErrorType(e error) bool{
	if _, ok := e.(*MessageTypeError); ok {
		return true
	} else if _, ok = e.(*ChannelNameError); ok {
		return true
	} else if _, ok = e.(*ProductIDError); ok {
		return true
	}
	return false
}