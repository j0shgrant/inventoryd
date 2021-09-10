package main

import (
	"context"
	"fmt"
	"github.com/ably/ably-go/ably"
	"github.com/j0shgrant/inventoryd/common"
	"os"
	"time"
)

func main() {
	// load environment variables
	ablyKey := os.Getenv("INVENTORYD_ABLY_KEY")
	if ablyKey == "" {
		_, _ = fmt.Fprintln(os.Stderr, "Environment Variable [INVENTORYD_ABLY_KEY] must be set to use inventoryctl.")
		os.Exit(1)
	}

	// initialise ably client
	client, err := ably.NewRealtime(ably.WithKey(ablyKey))
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// build context+timeout
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancelFunc()

	// list inventory
	presenceMsgs, err := client.Channels.Get("inventoryd").Presence.Get(ctx)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var errors []error
	for _, msg := range presenceMsgs {
		inventoryRecord, err := common.DecodeInventoryRecord(msg.Data)
		if err != nil {
			errors = append(errors, err)
			continue
		}

		fmt.Println(*inventoryRecord)
	}

	for _, err := range errors {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
}
