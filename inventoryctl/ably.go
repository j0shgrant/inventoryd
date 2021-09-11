package main

import (
	"context"
	"fmt"
	"github.com/ably/ably-go/ably"
	"github.com/j0shgrant/inventoryd/common"
)

type ChannelReference struct {
	Environment string
	Region      string
}

type RealtimeService struct {
	client *ably.Realtime
}

func NewRealtimeService(key string) (*RealtimeService, error) {
	client, err := ably.NewRealtime(ably.WithKey(key))
	if err != nil {
		return nil, err
	}

	svc := &RealtimeService{
		client: client,
	}

	return svc, nil
}

func (rs *RealtimeService) BatchPresence(channelRefs []ChannelReference, ctx context.Context) (Output, []error) {
	var errors []error
	messagesByChannel := make(map[ChannelReference][]*ably.PresenceMessage)

	// get all presence messages by channel
	for _, channelRef := range channelRefs {
		// reconstitute channel name
		channel := fmt.Sprintf("inventoryd:%s:%s", channelRef.Environment, channelRef.Region)

		// list present channel members
		presenceMessages, err := rs.client.Channels.Get(channel).Presence.Get(ctx)
		if err != nil {
			errors = append(errors, err)
			continue
		}

		messagesByChannel[channelRef] = presenceMessages
	}

	// build output
	var jsonErrors []error
	output := Output{}
	for channelRef, presenceMessages := range messagesByChannel {
		for _, msg := range presenceMessages {
			inventoryRecord, err := common.DecodeInventoryRecord(msg.Data)
			if err != nil {
				jsonErrors = append(jsonErrors, err)
				continue
			}

			environment, ok := output[channelRef.Environment]
			if !ok {
				output[channelRef.Environment] = make(Environment)
				environment = output[channelRef.Environment]
			}

			region, ok := environment[channelRef.Region]
			if !ok {
				environment[channelRef.Region] = make(Region)
				region = environment[channelRef.Region]
			}

			region[inventoryRecord.ID] = Server{
				Images: inventoryRecord.Images,
				Tags:   inventoryRecord.Tags,
			}
		}
	}

	errors = append(errors, jsonErrors...)

	return output, errors
}
