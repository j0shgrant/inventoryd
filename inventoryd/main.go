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
	defer func() {
		_ = logger.Sync()
	}()

	// load environment variables
	ablyKey := os.Getenv("INVENTORYD_ABLY_KEY")
	if ablyKey == "" {
		zap.S().Fatal("environment variable [INVENTORYD_ABLY_KEY] must be set to use inventoryd")
	}
	environment := os.Getenv("INVENTORYD_ENVIRONMENT")
	if environment == "" {
		zap.S().Fatal("environment variable [INVENTORYD_ENVIRONMENT] must be set to use inventoryd")
	}
	region := os.Getenv("INVENTORYD_REGION")
	if region == "" {
		zap.S().Fatal("environment variable [INVENTORYD_REGION] must be set to use inventoryd")
	}
	schedule := os.Getenv("INVENTORYD_CRON_SCHEDULE")
	if region == "" {
		zap.S().Fatal("environment variable [INVENTORYD_CRON_SCHEDULE] must be set to use inventoryd")
	}

	// log config
	zap.S().Infof("Starting inventoryd with config:")
	zap.S().Infof("Environment: %s", environment)
	zap.S().Infof("Region: %s", region)
	zap.S().Infof("Schedule: %s", schedule)

	// derive channel name (inventoryd:environment:region)
	channel := fmt.Sprintf("inventoryd:%s:%s", environment, region)

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
	err = cs.Run(schedule)
	if err != nil {
		zap.S().Fatal(err)
	}
}
