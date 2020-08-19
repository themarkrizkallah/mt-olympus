package main

import (
	"fmt"
	"strings"
)


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
