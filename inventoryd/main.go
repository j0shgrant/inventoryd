package main

import (
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"os"
)

func main() {
	// configure logging
	logger, err := zap.NewProduction()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	zap.ReplaceGlobals(logger)
	defer logger.Sync()

	// load environment variables
	ablyKey := os.Getenv("INVENTORYD_ABLY_KEY")
	if ablyKey == "" {
		_, _ = fmt.Fprintln(os.Stderr, "Environment Variable [INVENTORYD_ABLY_KEY] must be set to use inventoryctl.")
		os.Exit(1)
	}

	// derive channel name (inventoryd:environment:region)
	region := "eu-west-1"
	channel := fmt.Sprintf("inventoryd:production:%s", region)

	// initialise PresenceService
	ps, err := NewPresenceService(ablyKey, channel, uuid.NewString(), region)
	if err != nil {
		zap.S().Fatal(err)
	}

	// initialise DockerService
	cs, err := NewDockerService(ps)
	if err != nil {
		zap.S().Fatal(err)
	}

	// run schedule
	err = cs.Run("* * * * *")
	if err != nil {
		zap.S().Fatal(err)
	}
}
