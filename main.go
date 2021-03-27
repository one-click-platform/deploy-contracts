package main

import (
	"context"

	"github.com/one-click-platform/deploy-contracts/internal/config"
	"github.com/one-click-platform/deploy-contracts/internal/deployer"
	"gitlab.com/distributed_lab/kit/kv"
)

func main() {
	cfg := config.NewConfig(kv.MustFromEnv())
	log := cfg.Log()

	ctx := context.Background()

	dep, err := deployer.New(ctx, cfg)
	if err != nil {
		log.WithError(err).Fatal("failed to create deployer")
	}

	err = dep.Run(ctx, tasks())
	if err != nil {
		log.WithError(err).Fatal("failed to run deployer tasks")
	}

	log.Info("finished")
}
