package serial_test

import (
	stdCtx "context"
	"errors"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	"github.com/jahvon/flow/internal/runner"
	"github.com/jahvon/flow/internal/runner/serial"
	testUtils "github.com/jahvon/flow/tests/utils"
	"github.com/jahvon/flow/tools/builder"
	"github.com/jahvon/flow/types/executable"
)

func TestSerialRunner(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Serial Runner Suite")
}

var _ = Describe("SerialRunner", func() {
	var (
		ctx       *testUtils.ContextWithMocks
		serialRnr runner.Runner
	)

	BeforeEach(func() {
		ctx = testUtils.NewContextWithMocks(stdCtx.Background(), GinkgoT())
		runner.RegisterRunner(ctx.RunnerMock)
		serialRnr = serial.NewRunner()
	})

	AfterEach(func() {
		runner.Reset()
	})

	Context("Name", func() {
		It("should return the correct runner name", func() {
			Expect(serialRnr.Name()).To(Equal("serial"))
		})
	})

	Context("IsCompatible", func() {
		It("should return false when executable is nil", func() {
			Expect(serialRnr.IsCompatible(nil)).To(BeFalse())
		})

		It("should return false when executable type is nil", func() {
			executable := &executable.Executable{}
			Expect(serialRnr.IsCompatible(executable)).To(BeFalse())
		})

		It("should return true when executable type is serial", func() {
			executable := &executable.Executable{
				Serial: &executable.SerialExecutableType{},
			}
			Expect(serialRnr.IsCompatible(executable)).To(BeTrue())
		})
	})

	When("Executables with ref", func() {
		var (
			rootExec                                    *executable.Executable
			isSerialExec1, isSerialExec2, isSerialExec3 gomock.Matcher
		)

		BeforeEach(func() {
			ns := "examples"
			rootExec = builder.SerialExecByRef(
				builder.WithNamespaceName(ns),
				builder.WithWorkspaceName(ctx.Ctx.CurrentWorkspace.AssignedName()),
				builder.WithWorkspacePath(ctx.Ctx.CurrentWorkspace.Location()),
			)
			serialSpec := rootExec.Serial
			isSerialExec1 = testUtils.ExecWithRef(serialSpec.Refs[0])
			isSerialExec2 = testUtils.ExecWithRef(serialSpec.Refs[1])
			isSerialExec3 = testUtils.ExecWithRef(serialSpec.Refs[2])

			runner.RegisterRunner(serialRnr)
			runner.RegisterRunner(ctx.RunnerMock)
			ctx.RunnerMock.EXPECT().IsCompatible(rootExec).Return(false).AnyTimes()
		})

		It("should execute in order", func() {
			promptedEnv := make(map[string]string)

			mockRunner := ctx.RunnerMock
			mockRunner.EXPECT().IsCompatible(isSerialExec1).Return(true).Times(1)
			mockRunner.EXPECT().Exec(ctx.Ctx, isSerialExec1, promptedEnv).Return(nil).Times(1)

			mockRunner.EXPECT().IsCompatible(isSerialExec2).Return(true).Times(1)
			mockRunner.EXPECT().Exec(ctx.Ctx, isSerialExec2, promptedEnv).Return(nil).Times(1)

			mockRunner.EXPECT().IsCompatible(isSerialExec3).Return(true).Times(1)
			mockRunner.EXPECT().Exec(ctx.Ctx, isSerialExec3, promptedEnv).Return(nil).Times(1)

			Expect(serialRnr.Exec(ctx.Ctx, rootExec, promptedEnv)).To(Succeed())
		})

		It("should fail fast when enabled", func() {
			promptedEnv := make(map[string]string)

			mockRunner := ctx.RunnerMock
			mockRunner.EXPECT().IsCompatible(isSerialExec1).Return(true).Times(1)
			mockRunner.EXPECT().Exec(ctx.Ctx, isSerialExec1, promptedEnv).Return(nil).Times(1)

			mockRunner.EXPECT().IsCompatible(isSerialExec2).Return(true).Times(1)
			mockRunner.EXPECT().Exec(ctx.Ctx, isSerialExec2, promptedEnv).Return(errors.New("error")).Times(1)

			rootExec.Serial.FailFast = true
			Expect(serialRnr.Exec(ctx.Ctx, rootExec, promptedEnv)).ToNot(Succeed())
		})

		FIt("should not fail fast when disabled", func() {
			promptedEnv := make(map[string]string)
			mockRunner := ctx.RunnerMock
			mockLogger := ctx.Logger
			mockCache := ctx.ExecutableCache

			mockCache.EXPECT().GetExecutableByRef(ctx.Logger, gomock.Any()).Return(builder.SimpleExec(), nil).Times(3)
			mockRunner.EXPECT().IsCompatible(isSerialExec1).Return(true).Times(1)
			mockRunner.EXPECT().Exec(ctx.Ctx, isSerialExec1, promptedEnv).Return(nil).Times(1)

			// mockCache.EXPECT().GetExecutableByRef(ctx.Logger, isSerialExec2).Return(subExec, nil).Times(1)
			mockRunner.EXPECT().IsCompatible(isSerialExec2).Return(true).Times(1)
			mockRunner.EXPECT().Exec(ctx.Ctx, isSerialExec2, promptedEnv).Return(errors.New("error")).Times(1)

			// mockCache.EXPECT().GetExecutableByRef(ctx.Logger, isSerialExec3).Return(subExec, nil).Times(1)
			mockRunner.EXPECT().IsCompatible(isSerialExec3).Return(true).Times(1)
			mockRunner.EXPECT().Exec(ctx.Ctx, isSerialExec3, promptedEnv).Return(nil).Times(1)
			mockLogger.EXPECT().Error(gomock.Any(), gomock.Any()).Times(1)

			Expect(serialRnr.Exec(ctx.Ctx, rootExec, promptedEnv)).ToNot(Succeed())
		})
	})

	When("Executables with ref configs", func() {
		var (
			rootExec                                    *executable.Executable
			isSerialExec1, isSerialExec2, isSerialExec3 gomock.Matcher
		)
		BeforeEach(func() {
			ns := "examples"
			rootExec = builder.SerialExecByRefConfig(
				builder.WithNamespaceName(ns),
				builder.WithWorkspaceName(ctx.Ctx.CurrentWorkspace.AssignedName()),
				builder.WithWorkspacePath(ctx.Ctx.CurrentWorkspace.Location()),
			)
			serialSpec := rootExec.Serial
			isSerialExec1 = testUtils.ExecWithRef(serialSpec.Execs[0].Ref)
			isSerialExec2 = testUtils.ExecWithRef(serialSpec.Execs[1].Ref)
			isSerialExec3 = testUtils.ExecWithCmd(serialSpec.Execs[2].Cmd)

			runner.RegisterRunner(serialRnr)
			runner.RegisterRunner(ctx.RunnerMock)
			ctx.RunnerMock.EXPECT().IsCompatible(rootExec).Return(false).AnyTimes()
		})

		It("should execute in order", func() {
			promptedEnv := make(map[string]string)

			mockRunner := ctx.RunnerMock
			mockRunner.EXPECT().IsCompatible(isSerialExec1).Return(true).Times(1)
			mockRunner.EXPECT().Exec(ctx.Ctx, isSerialExec1, promptedEnv).Return(nil).Times(1)

			mockRunner.EXPECT().IsCompatible(isSerialExec2).Return(true).Times(1)
			mockRunner.EXPECT().Exec(ctx.Ctx, isSerialExec2, promptedEnv).Return(nil).Times(1)

			mockRunner.EXPECT().IsCompatible(isSerialExec3).Return(true).Times(1)
			mockRunner.EXPECT().Exec(ctx.Ctx, isSerialExec3, promptedEnv).Return(nil).Times(1)

			Expect(serialRnr.Exec(ctx.Ctx, rootExec, promptedEnv)).To(Succeed())
		})

		Context("when retries are set on a failed ref config", func() {
			BeforeEach(func() {
				rootExec.Serial.Execs[1].Retries = 2
			})

			When("fail fast is disabled", func() {
				It("should be retried until attempted max times", func() {
					mockRunner := ctx.RunnerMock
					mockCache := ctx.ExecutableCache
					mockLogger := ctx.Logger
					promptedEnv := make(map[string]string)

					serial1Exec := rootExec.Serial.Execs[0]
					mockCache.EXPECT().GetExecutableByRef(ctx.Logger, serial1Exec.Ref).Return(serial1Exec, nil).Times(1)
					mockRunner.EXPECT().IsCompatible(isSerialExec1).Return(true).Times(1)
					mockRunner.EXPECT().Exec(ctx, isSerialExec1, promptedEnv).Return(nil).Times(1)

					serial2Exec := rootExec.Serial.Execs[1]
					serial2Prompt := map[string]string{"ARG1": "hello", "ARG2": "123"}
					mockCache.EXPECT().GetExecutableByRef(ctx.Logger, serial2Exec.Ref).Return(serial2Exec, nil).Times(3)
					mockRunner.EXPECT().IsCompatible(isSerialExec2).Return(true).Times(3)
					mockRunner.EXPECT().Exec(ctx.Ctx, isSerialExec2, serial2Prompt).Return(errors.New("error")).Times(3)
					mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any()).Times(2)
					mockLogger.EXPECT().Error(gomock.Any(), gomock.Any()).Times(1)

					serial3Exec := rootExec.Serial.Execs[2]
					mockCache.EXPECT().GetExecutableByRef(ctx.Logger, serial3Exec.Ref).Return(serial3Exec, nil).Times(1)
					mockRunner.EXPECT().IsCompatible(isSerialExec3).Return(true).Times(1)
					mockRunner.EXPECT().Exec(ctx.Ctx, isSerialExec3, promptedEnv).Return(nil).Times(1)

					Expect(serialRnr.Exec(ctx.Ctx, rootExec, make(map[string]string))).To(Succeed())
				})
			})

			When("fail fast is enabled", func() {
				BeforeEach(func() {
					rootExec.Serial.FailFast = true
				})

				It("should fail fast after max attempts when enabled", func() {
					mockRunner := ctx.RunnerMock
					mockCache := ctx.ExecutableCache
					mockLogger := ctx.Logger
					promptedEnv := make(map[string]string)

					serial1Exec := rootExec.Serial.Execs[0]
					mockCache.EXPECT().GetExecutableByRef(ctx.Logger, serial1Exec.Ref).Return(serial1Exec, nil).Times(1)
					mockRunner.EXPECT().IsCompatible(isSerialExec1).Return(true).Times(1)
					mockRunner.EXPECT().Exec(ctx, isSerialExec1, promptedEnv).Return(nil).Times(1)

					serial2Exec := rootExec.Serial.Execs[1]
					serial2Prompt := map[string]string{"ARG1": "hello", "ARG2": "123"}
					mockCache.EXPECT().GetExecutableByRef(ctx.Logger, serial2Exec.Ref).Return(serial2Exec, nil).Times(3)
					mockRunner.EXPECT().IsCompatible(isSerialExec2).Return(true).Times(3)
					mockRunner.EXPECT().Exec(ctx.Ctx, isSerialExec2, serial2Prompt).Return(errors.New("error")).Times(3)
					mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any()).Times(2)
					mockLogger.EXPECT().Error(gomock.Any(), gomock.Any()).Times(1)

					Expect(serialRnr.Exec(ctx.Ctx, rootExec, make(map[string]string))).ToNot(Succeed())
				})
			})
		})

		Context("when retries are not enabled on a failed ref config", func() {
			When("fail fast is disabled", func() {
				It("should be retried until attempted max times", func() {
					mockRunner := ctx.RunnerMock
					mockCache := ctx.ExecutableCache
					mockLogger := ctx.Logger
					promptedEnv := make(map[string]string)

					serial1Exec := rootExec.Serial.Execs[0]
					mockCache.EXPECT().GetExecutableByRef(ctx.Logger, serial1Exec.Ref).Return(serial1Exec, nil).Times(1)
					mockRunner.EXPECT().IsCompatible(isSerialExec1).Return(true).Times(1)
					mockRunner.EXPECT().Exec(ctx, isSerialExec1, promptedEnv).Return(nil).Times(1)

					serial2Exec := rootExec.Serial.Execs[1]
					serial2Prompt := map[string]string{"ARG1": "hello", "ARG2": "123"}
					mockCache.EXPECT().GetExecutableByRef(ctx.Logger, serial2Exec.Ref).Return(serial2Exec, nil).Times(1)
					mockRunner.EXPECT().IsCompatible(isSerialExec2).Return(true).Times(1)
					mockRunner.EXPECT().Exec(ctx.Ctx, isSerialExec2, serial2Prompt).Return(errors.New("error")).Times(1)
					mockLogger.EXPECT().Error(gomock.Any(), gomock.Any()).Times(1)

					serial3Exec := rootExec.Serial.Execs[2]
					mockCache.EXPECT().GetExecutableByRef(ctx.Logger, serial3Exec.Ref).Return(serial3Exec, nil).Times(1)
					mockRunner.EXPECT().IsCompatible(isSerialExec3).Return(true).Times(1)
					mockRunner.EXPECT().Exec(ctx.Ctx, isSerialExec3, promptedEnv).Return(nil).Times(1)

					Expect(serialRnr.Exec(ctx.Ctx, rootExec, make(map[string]string))).To(Succeed())
				})
			})

			When("fail fast is enabled", func() {
				BeforeEach(func() {
					rootExec.Serial.FailFast = true
				})

				It("should fail fast after max attempts when enabled", func() {
					mockRunner := ctx.RunnerMock
					mockCache := ctx.ExecutableCache
					promptedEnv := make(map[string]string)

					serial1Exec := rootExec.Serial.Execs[0]
					mockCache.EXPECT().GetExecutableByRef(ctx.Logger, serial1Exec.Ref).Return(serial1Exec, nil).Times(1)
					mockRunner.EXPECT().IsCompatible(isSerialExec1).Return(true).Times(1)
					mockRunner.EXPECT().Exec(ctx, isSerialExec1, promptedEnv).Return(nil).Times(1)

					serial2Exec := rootExec.Serial.Execs[1]
					serial2Prompt := map[string]string{"ARG1": "hello", "ARG2": "123"}
					mockCache.EXPECT().GetExecutableByRef(ctx.Logger, serial2Exec.Ref).Return(serial2Exec, nil).Times(1)
					mockRunner.EXPECT().IsCompatible(isSerialExec2).Return(true).Times(1)
					mockRunner.EXPECT().Exec(ctx.Ctx, isSerialExec2, serial2Prompt).Return(errors.New("error")).Times(1)

					Expect(serialRnr.Exec(ctx.Ctx, rootExec, make(map[string]string))).ToNot(Succeed())
				})
			})
		})
	})
})
