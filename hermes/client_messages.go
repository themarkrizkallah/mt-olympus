package main

import (
	"encoding/json"
)

/*
User supported messages are:
	- "subscribe"
	- "unsubscribe"

Alongside the respective channel(s) of course, see channels.go
*/
var acceptedMsgTypes = []string{"subscribe", "unsubscribe"}


type ChannelMessage struct {
	Name       string   `json:"name"`
	ProductIDs []string `json:"product_ids"`
}

type ChannelMessageJSON ChannelMessage

func (cm *ChannelMessage) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, (*ChannelMessageJSON)(cm)); err != nil {
		return err
	}

	// Validate proper channel name
	if !stringFound(cm.Name, acceptedChanNames) {
		return ChannelNameError{cm.Name}
	}

	// Validate product ids
	for _, id := range cm.ProductIDs {
		if !stringFound(id, acceptedProductIDs){
			return ProductIDError{id}
		}
	}

	return nil
}


type SubscribeRequest struct {
	Client *Client
	SubMsg SubscribeMessage
}

type SubscribeMessage struct {
	MsgType     string           `json:"type"`
	ChannelMsgs []ChannelMessage `json:"channels"`
}

type SubscribeMessageJSON SubscribeMessage

func (sm *SubscribeMessage) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, (*SubscribeMessageJSON)(sm)); err != nil {
		return err
	}

	// Validate proper message type
	if !stringFound(sm.MsgType, acceptedMsgTypes) {
		return MessageTypeError{sm.MsgType}
	}

	return nil
}
