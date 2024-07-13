package runner_test

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/runner"
	"github.com/jahvon/flow/internal/runner/mocks"
	"github.com/jahvon/flow/types/executable"
)

func TestRunner(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Runner Suite")
}

var _ = Describe("Runner", func() {
	var (
		ctrl       *gomock.Controller
		mockRunner *mocks.MockRunner
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockRunner = mocks.NewMockRunner(ctrl)
		runner.RegisterRunner(mockRunner)
	})

	AfterEach(func() {
		ctrl.Finish()
		runner.Reset()
	})

	Describe("Exec", func() {
		It("should execute the runner correctly", func() {
			ctx := &context.Context{}
			executable := &executable.Executable{
				Name: "test-executable",
			}
			promptedEnv := make(map[string]string)

			mockRunner.EXPECT().IsCompatible(executable).Return(true)
			mockRunner.EXPECT().Exec(ctx, executable, promptedEnv).Return(nil)
			Expect(runner.Exec(ctx, executable, promptedEnv)).To(Succeed())
		})

		It("should return error when no compatible runner is found", func() {
			ctx := &context.Context{}
			exec := &executable.Executable{
				Name: "test-exec",
			}
			promptedEnv := make(map[string]string)

			mockRunner.EXPECT().IsCompatible(exec).Return(false)

			err := runner.Exec(ctx, exec, promptedEnv)
			Expect(err.Error()).To(ContainSubstring("compatible runner not found"))
		})

		It("should return error when execution times out", func() {
			ctx := &context.Context{}
			exec := &executable.Executable{
				Name:    "test-exec",
				Timeout: 250 * time.Millisecond,
			}
			promptedEnv := make(map[string]string)

			mockRunner.EXPECT().IsCompatible(exec).Return(true)
			mockRunner.EXPECT().Exec(ctx, exec, promptedEnv).DoAndReturn(
				func(ctx *context.Context, executable *executable.Executable, promptedEnv map[string]string) error {
					time.Sleep(2 * time.Second)
					return nil
				})

			err := runner.Exec(ctx, exec, promptedEnv)
			Expect(err.Error()).To(ContainSubstring("timeout"))
		})
	})
})
