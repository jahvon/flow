package logger_test

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/flowexec/tuikit/io"

	"github.com/flowexec/flow/internal/logger"
	"github.com/flowexec/flow/types/config"
)

func TestLogger(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Logger Suite")
}

var _ = Describe("Global Logger", func() {
	BeforeEach(func() {
		logger.Reset()
	})

	It("should initialize the logger once", func() {
		cfg := &config.Config{
			Theme: config.ConfigThemeDefault,
		}
		opts := logger.InitOptions{
			Config:  cfg,
			StdOut:  os.Stdout,
			LogMode: io.Logfmt,
		}

		logger.Init(opts)
		logger1 := logger.Get()

		logger.Init(opts)
		logger2 := logger.Get()

		Expect(logger1).To(Equal(logger2))
	})

	It("should panic when getting logger before initialization", func() {
		Expect(func() {
			logger.Get()
		}).To(Panic())
	})

	It("should allow reinitialization after reset", func() {
		cfg := &config.Config{
			Theme: config.ConfigThemeDefault,
		}
		opts := logger.InitOptions{
			Config:  cfg,
			StdOut:  os.Stdout,
			LogMode: io.Logfmt,
		}

		logger.Init(opts)
		logger1 := logger.Get()

		logger.Reset()

		logger.Init(opts)
		logger2 := logger.Get()

		Expect(logger1).ToNot(Equal(logger2))
	})
})