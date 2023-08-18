package context

import (
	"context"

	"github.com/jahvon/tbox/internal/config"
)

const configKey = "config"

type RootCtx struct {
	context.Context
}

func NewRootCtx() *RootCtx {
	ctx := context.Background()
	conf := config.LoadConfig()
	ctx = context.WithValue(ctx, configKey, conf)
	return &RootCtx{ctx}
}

func (r *RootCtx) Config() *config.RootConfig {
	return r.Value(configKey).(*config.RootConfig)
}
