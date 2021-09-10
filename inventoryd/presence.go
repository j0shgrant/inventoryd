package main

import (
	"context"
	"encoding/json"
	"github.com/ably/ably-go/ably"
	"github.com/j0shgrant/inventoryd/common"
	"go.uber.org/zap"
	"time"
)

type PresenceService struct {
	channel          *ably.RealtimeChannel
	clientId, region string
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
		channel:  client.Channels.Get(channel),
		clientId: clientId,
		region:   region,
	}

	return svc, nil
}

func (ps *PresenceService) Register(runningImages map[string]string) error {
	msg, err := json.Marshal(runningImages)
	if err != nil {
		return err
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	return ps.channel.Presence.Enter(ctx, msg)
}

func (ps *PresenceService) Update(runningImages map[string]string) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	msg, err := common.EncodeInventoryRecord(ps.clientId, ps.region, runningImages)
	if err != nil {
		return err
	}

	return ps.channel.Presence.Update(ctx, msg)
}

func (ps *PresenceService) Deregister() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	return ps.channel.Presence.Leave(ctx, nil)
}
