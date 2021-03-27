package config

import (
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
)

var Version string

type config struct {
	comfig.Logger
	getter kv.Getter
	Node
}

type Config interface {
	comfig.Logger
	Node
}

func NewConfig(getter kv.Getter) Config {
	return &config{
		getter:         getter,
		Logger:         comfig.NewLogger(getter, comfig.LoggerOpts{Release: Version}),
		Node:           NewNode(getter),
	}
}
