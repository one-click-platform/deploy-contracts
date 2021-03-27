package main

import (
	"github.com/one-click-platform/deploy-contracts/internal/config"
	"gitlab.com/distributed_lab/kit/kv"
)

func main() {
	cfg := config.NewConfig(kv.MustFromEnv())
	log := cfg.Log()

	log.Info("finished")
}
