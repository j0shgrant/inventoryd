package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ably/ably-go/ably"
	"log"
	"os"
	"strings"
	"time"
)

type Channel struct {
	ID string `json:"channelId"`
}

func main() {
	// load environment variables
	ablyKey := os.Getenv("INVENTORYD_ABLY_KEY")
	if ablyKey == "" {
		_, _ = fmt.Fprintln(os.Stderr, "Environment Variable [INVENTORYD_ABLY_KEY] must be set to use inventoryctl.")
		os.Exit(1)
	}

	// initialise rest client
	rest, err := ably.NewREST(ably.WithKey(ablyKey), ably.WithUseBinaryProtocol(false))
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// build context+timeout
	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()

	// list channels
	var allChannels []Channel
	pages, err := rest.Request("get", "/channels").Pages(ctx)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for {
		if !pages.Next(ctx) {
			break
		}

		var channels []Channel
		err = pages.Items(&channels)
		if err != nil {
			log.Fatal(err)
		}

		allChannels = append(allChannels, channels...)
	}

	// filter out invalid channels
	var channelReferences []ChannelReference
	var invalidChannelNames []string
	for _, channel := range allChannels {
		channelNameTokens := strings.Split(channel.ID, ":")

		// channel name is invalid if it does meet pattern [inventoryd:<ENVIRONMENT>:<REGION>]
		if len(channelNameTokens) != 3 {
			invalidChannelNames = append(invalidChannelNames, channel.ID)
			continue
		}
		if channelNameTokens[0] != "inventoryd" {
			invalidChannelNames = append(invalidChannelNames, channel.ID)
			continue
		}

		channelReferences = append(channelReferences, ChannelReference{
			Environment: channelNameTokens[1],
			Region:      channelNameTokens[2],
		})
	}

	// handle invalid channels
	if len(invalidChannelNames) > 0 {
		allInvalidChannels := strings.Join(invalidChannelNames, ",")
		_, _ = fmt.Fprintf(os.Stderr, "The following invalid channel names did not meet the format [inventoryd:environment:region]: [%s]\n", allInvalidChannels)
	}

	// initialise RealtimeService
	rs, err := NewRealtimeService(ablyKey)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// list all present instances
	output, batchErrors := rs.BatchPresence(channelReferences, ctx)
	if len(batchErrors) > 0 {
		for _, err := range batchErrors {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}
	}

	// print output
	outputBytes, err := json.Marshal(output)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(string(outputBytes))
}
