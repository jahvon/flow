package engine_test

import (
	"context"
	"errors"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/flowexec/flow/internal/runner/engine"
)

func TestEngine_Execute(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Execute Engine Suite")
}

var _ = Describe("e.Execute", func() {
	var (
		eng    engine.Engine
		ctx    context.Context
		cancel context.CancelFunc
	)

	BeforeEach(func() {
		eng = engine.NewExecEngine()
		ctx, cancel = context.WithCancel(context.Background())
	})

	AfterEach(func() {
		cancel()
	})

	Context("Parallel execution", func() {
		It("should execute execs in parallel", func() {
			execs := []engine.Exec{
				{ID: "exec1", Function: func() error { time.Sleep(100 * time.Millisecond); return nil }},
				{ID: "exec2", Function: func() error { return nil }},
			}

			start := time.Now()
			ff := false
			summary := eng.Execute(ctx, execs, engine.WithMode(engine.Parallel), engine.WithFailFast(&ff))
			duration := time.Since(start)

			Expect(summary.Results).To(HaveLen(2))
			Expect(summary.Results[0].Error).NotTo(HaveOccurred())
			Expect(summary.Results[1].Error).NotTo(HaveOccurred())
			Expect(duration).To(BeNumerically("<", 200*time.Millisecond))
		})

		It("should handle exec failures with fail fast", func() {
			execs := []engine.Exec{
				{ID: "exec1", Function: func() error { return errors.New("error") }},
				{ID: "exec2", Function: func() error { time.Sleep(100 * time.Millisecond); return nil }},
			}

			ff := true
			summary := eng.Execute(ctx, execs, engine.WithMode(engine.Parallel), engine.WithFailFast(&ff))

			Expect(summary.Results).To(HaveLen(2))
			Expect(summary.Results[0].Error).To(HaveOccurred())
			Expect(summary.Results[1].Error).ToNot(HaveOccurred())
			Expect(summary.HasErrors()).To(BeTrue())
		})

		It("should limit the number of concurrent execs", func() {
			execs := []engine.Exec{
				{ID: "exec1", Function: func() error { time.Sleep(100 * time.Millisecond); return nil }},
				{ID: "exec2", Function: func() error { time.Sleep(100 * time.Millisecond); return nil }},
				{ID: "exec3", Function: func() error { time.Sleep(100 * time.Millisecond); return nil }},
				{ID: "exec4", Function: func() error { time.Sleep(100 * time.Millisecond); return nil }},
				{ID: "exec5", Function: func() error { time.Sleep(100 * time.Millisecond); return nil }},
			}

			start := time.Now()
			ff := false
			summary := eng.Execute(ctx, execs,
				engine.WithMode(engine.Parallel), engine.WithFailFast(&ff), engine.WithMaxThreads(2))
			duration := time.Since(start)

			Expect(summary.Results).To(HaveLen(5))
			Expect(summary.Results[0].Error).NotTo(HaveOccurred())
			Expect(summary.Results[1].Error).NotTo(HaveOccurred())
			Expect(summary.Results[2].Error).NotTo(HaveOccurred())
			Expect(summary.Results[3].Error).NotTo(HaveOccurred())
			Expect(summary.Results[4].Error).NotTo(HaveOccurred())
			Expect(duration).To(BeNumerically(">=", 250*time.Millisecond))
		})
	})

	Context("Serial execution", func() {
		It("should execute execs serially", func() {
			execs := []engine.Exec{
				{ID: "exec1", Function: func() error { time.Sleep(100 * time.Millisecond); return nil }},
				{ID: "exec2", Function: func() error { time.Sleep(110 * time.Millisecond); return nil }},
			}

			start := time.Now()
			ff := false
			summary := eng.Execute(ctx, execs, engine.WithMode(engine.Serial), engine.WithFailFast(&ff))
			duration := time.Since(start)

			Expect(summary.Results).To(HaveLen(2))
			Expect(summary.Results[0].Error).NotTo(HaveOccurred())
			Expect(summary.Results[1].Error).NotTo(HaveOccurred())
			Expect(duration).To(BeNumerically(">=", 200*time.Millisecond))
		})

		It("should handle exec failures with fail fast", func() {
			execs := []engine.Exec{
				{ID: "exec1", Function: func() error { return errors.New("error") }},
				{ID: "exec2", Function: func() error { return nil }},
			}

			ff := true
			summary := eng.Execute(ctx, execs, engine.WithMode(engine.Serial), engine.WithFailFast(&ff))

			Expect(summary.Results).To(HaveLen(1))
			Expect(summary.Results[0].Error).To(HaveOccurred())
			Expect(summary.HasErrors()).To(BeTrue())
		})
	})
})
