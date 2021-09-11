package main

import (
	"context"
	"github.com/ably/ably-go/ably"
	"github.com/j0shgrant/inventoryd/common"
	"go.uber.org/zap"
	"time"
)

type PresenceService struct {
	client                    *ably.Realtime
	channel, clientId, region string
}

func NewPresenceService(key, channel, clientId, region string) (*PresenceService, error) {
	zap.S().Infof("Connecting to Ably with ClientId %s", clientId)
	client, err := ably.NewRealtime(
		ably.WithKey(key),
		ably.WithClientID(clientId),
	)
	if err != nil {
		return nil, err
	}

	zap.S().Infof("Connecting to Channel %s with ClientId %s", channel, clientId)
	svc := &PresenceService{
		client:   client,
		channel:  channel,
		clientId: clientId,
		region:   region,
	}

	return svc, nil
}

func (ps *PresenceService) Update(runningImages map[string]string) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	msg, err := common.EncodeInventoryRecord(ps.clientId, ps.region, runningImages)
	if err != nil {
		return err
	}

	return ps.client.Channels.Get(ps.channel).Presence.Update(ctx, msg)
}

func (ps *PresenceService) Deregister() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	return ps.client.Channels.Get(ps.channel).Presence.Leave(ctx, nil)
}
