package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/j0shgrant/inventoryd/inventoryd/internal"
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

	// initialise PresenceService
	ps, err := internal.NewPresenceService("4JW6ZA.9TGsRA:0677DoZU_HmkH9_1", "inventoryd", uuid.NewString())
	if err != nil {
		zap.S().Fatal(err)
	}

	// initialise DockerService
	cs, err := internal.NewDockerService(ps)
	if err != nil {
		zap.S().Fatal(err)
	}

	// run schedule
	err = cs.Run("* * * * *")
	if err != nil {
		zap.S().Fatal(err)
	}
}
