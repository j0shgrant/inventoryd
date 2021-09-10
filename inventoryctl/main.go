package main

import (
	"context"
	"fmt"
	"github.com/ably/ably-go/ably"
	"os"
	"time"
)

func main() {
	// initialise ably client
	client, err := ably.NewRealtime(ably.WithKey("4JW6ZA.9TGsRA:0677DoZU_HmkH9_1"))
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

	for _, msg := range presenceMsgs {
		fmt.Println(msg.ClientID)
		fmt.Println(msg.Data.(string))
	}
}