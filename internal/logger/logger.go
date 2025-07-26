package logger

import (
	"os"
	"sync"

	"github.com/flowexec/tuikit/io"
	"golang.org/x/exp/slices"

	"github.com/flowexec/flow/internal/filesystem"
	flowIO "github.com/flowexec/flow/internal/io"
	"github.com/flowexec/flow/types/config"
	"github.com/flowexec/flow/types/executable"
)

var (
	globalLogger io.Logger
	once         sync.Once
)

type InitOptions struct {
	Config  *config.Config
	StdOut  *os.File
	LogMode io.LogMode
}

// Init initializes the global logger with the provided options.
// This function is safe to call multiple times - only the first call will initialize the logger.
func Init(opts InitOptions) {
	once.Do(func() {
		loggerOpts := []io.LoggerOptions{
			io.WithOutput(opts.StdOut),
			io.WithTheme(flowIO.Theme(opts.Config.Theme.String())),
			io.WithMode(opts.LogMode),
		}

		// only create a log archive file for exec commands
		if args := os.Args; len(args) > 0 && slices.Contains(executable.ValidVerbs(), executable.Verb(args[0])) {
			loggerOpts = append(loggerOpts, io.WithArchiveDirectory(filesystem.LogsDir()))
		}

		globalLogger = io.NewLogger(loggerOpts...)
	})
}

// Get returns the global logger instance.
// The logger must be initialized with Init() before calling this function.
func Get() io.Logger {
	if globalLogger == nil {
		panic("global logger not initialized - call logger.Init() first")
	}
	return globalLogger
}

// Reset resets the global logger for testing purposes.
// This should only be used in tests.
func Reset() {
	once = sync.Once{}
	globalLogger = nil
}
