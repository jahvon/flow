package cache

import (
	"fmt"

	"github.com/flowexec/flow/internal/context"
	"github.com/flowexec/flow/internal/io/common"
)

func PrintCache(ctx *context.Context, cache map[string]string, format string) {
	if cache == nil {
		return
	}
	logger := ctx.Logger
	output := cacheData{Cache: cache}
	switch common.NormalizeFormat(logger, format) {
	case common.YAMLFormat:
		str, err := output.YAML()
		if err != nil {
			logger.Fatalf("Failed to marshal cache - %v", err)
		}
		_, _ = fmt.Fprint(ctx.StdOut(), str)
	case common.JSONFormat:
		str, err := output.JSON()
		if err != nil {
			logger.Fatalf("Failed to marshal cache - %v", err)
		}
		_, _ = fmt.Fprint(ctx.StdOut(), str)
	}
}
