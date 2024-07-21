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
			rootExec *executable.Executable
			subExecs executable.ExecutableList
		)

		BeforeEach(func() {
			ns := "examples"
			rootExec = builder.SerialExecByRef(
				builder.WithNamespaceName(ns),
				builder.WithWorkspaceName(ctx.Ctx.CurrentWorkspace.AssignedName()),
				builder.WithWorkspacePath(ctx.Ctx.CurrentWorkspace.Location()),
			)
			execFlowfile := builder.ExamplesExecFlowFile(
				builder.WithNamespaceName(ns),
				builder.WithWorkspaceName(ctx.Ctx.CurrentWorkspace.AssignedName()),
				builder.WithWorkspacePath(ctx.Ctx.CurrentWorkspace.Location()),
			)
			subExecs = testUtils.FindSubExecs(rootExec, executable.FlowFileList{execFlowfile})

			runner.RegisterRunner(serialRnr)
			runner.RegisterRunner(ctx.RunnerMock)
			ctx.RunnerMock.EXPECT().IsCompatible(rootExec).Return(false).AnyTimes()
		})

		It("should execute in order", func() {
			promptedEnv := make(map[string]string)
			mockRunner := ctx.RunnerMock
			mockCache := ctx.ExecutableCache

			for _, e := range subExecs {
				isSerialExec := testUtils.ExecWithRef(e.Ref())
				mockCache.EXPECT().GetExecutableByRef(ctx.Logger, e.Ref()).Return(e, nil).Times(1)
				mockRunner.EXPECT().IsCompatible(isSerialExec).Return(true).Times(1)
				mockRunner.EXPECT().Exec(ctx.Ctx, isSerialExec, promptedEnv).Return(nil).Times(1)
			}

			Expect(serialRnr.Exec(ctx.Ctx, rootExec, promptedEnv)).To(Succeed())
		})

		It("should fail fast when enabled", func() {
			promptedEnv := make(map[string]string)
			mockRunner := ctx.RunnerMock
			mockCache := ctx.ExecutableCache

			for i, e := range subExecs {
				switch i {
				case 0:
					isSerialExec := testUtils.ExecWithRef(e.Ref())
					mockCache.EXPECT().GetExecutableByRef(ctx.Logger, e.Ref()).Return(e, nil).Times(1)
					mockRunner.EXPECT().IsCompatible(isSerialExec).Return(true).Times(1)
					mockRunner.EXPECT().Exec(ctx.Ctx, isSerialExec, promptedEnv).Return(nil).Times(1)
				case 1:
					isSerialExec := testUtils.ExecWithRef(e.Ref())
					mockCache.EXPECT().GetExecutableByRef(ctx.Logger, e.Ref()).Return(e, nil).Times(1)
					mockRunner.EXPECT().IsCompatible(isSerialExec).Return(true).Times(1)
					mockRunner.EXPECT().Exec(ctx.Ctx, isSerialExec, promptedEnv).Return(errors.New("error")).Times(1)
				}
			}

			rootExec.Serial.FailFast = true
			Expect(serialRnr.Exec(ctx.Ctx, rootExec, promptedEnv)).ToNot(Succeed())
		})

		It("should not fail fast when disabled", func() {
			promptedEnv := make(map[string]string)
			mockRunner := ctx.RunnerMock
			mockLogger := ctx.Logger
			mockCache := ctx.ExecutableCache

			for i, e := range subExecs {
				switch i {
				case 0, 2:
					isSerialExec := testUtils.ExecWithRef(e.Ref())
					mockCache.EXPECT().GetExecutableByRef(ctx.Logger, e.Ref()).Return(e, nil).Times(1)
					mockRunner.EXPECT().IsCompatible(isSerialExec).Return(true).Times(1)
					mockRunner.EXPECT().Exec(ctx.Ctx, isSerialExec, promptedEnv).Return(nil).Times(1)
				case 1:
					isSerialExec := testUtils.ExecWithRef(e.Ref())
					mockCache.EXPECT().GetExecutableByRef(ctx.Logger, e.Ref()).Return(e, nil).Times(1)
					mockRunner.EXPECT().IsCompatible(isSerialExec).Return(true).Times(1)
					mockRunner.EXPECT().Exec(ctx.Ctx, isSerialExec, promptedEnv).Return(errors.New("error")).Times(1)
					mockLogger.EXPECT().Error(gomock.Any(), gomock.Any()).Times(1)
				}
			}

			Expect(serialRnr.Exec(ctx.Ctx, rootExec, promptedEnv)).ToNot(Succeed())
		})
	})

	When("Executables with ref configs", func() {
		var (
			rootExec *executable.Executable
			subExecs executable.ExecutableList
		)
		BeforeEach(func() {
			ns := "examples"
			rootExec = builder.SerialExecByRefConfig(
				builder.WithNamespaceName(ns),
				builder.WithWorkspaceName(ctx.Ctx.CurrentWorkspace.AssignedName()),
				builder.WithWorkspacePath(ctx.Ctx.CurrentWorkspace.Location()),
			)
			execFlowfile := builder.ExamplesExecFlowFile(
				builder.WithNamespaceName(ns),
				builder.WithWorkspaceName(ctx.Ctx.CurrentWorkspace.AssignedName()),
				builder.WithWorkspacePath(ctx.Ctx.CurrentWorkspace.Location()),
			)
			subExecs = testUtils.FindSubExecs(rootExec, executable.FlowFileList{execFlowfile})

			runner.RegisterRunner(serialRnr)
			runner.RegisterRunner(ctx.RunnerMock)
			ctx.RunnerMock.EXPECT().IsCompatible(rootExec).Return(false).AnyTimes()
		})

		It("should execute in order", func() {
			promptedEnv := make(map[string]string)
			mockRunner := ctx.RunnerMock
			mockCache := ctx.ExecutableCache

			for i, e := range subExecs {
				switch i {
				case 0:
					isSerialExec := testUtils.ExecWithRef(e.Ref())
					mockCache.EXPECT().GetExecutableByRef(ctx.Logger, e.Ref()).Return(e, nil).Times(1)
					mockRunner.EXPECT().IsCompatible(isSerialExec).Return(true).Times(1)
					mockRunner.EXPECT().Exec(ctx.Ctx, isSerialExec, promptedEnv).Return(nil).Times(1)
				case 1:
					isSerialExec := testUtils.ExecWithRef(e.Ref())
					serialPrompt := map[string]string{"ARG1": "hello", "ARG2": "123"}
					mockCache.EXPECT().GetExecutableByRef(ctx.Logger, e.Ref()).Return(e, nil).Times(1)
					mockRunner.EXPECT().IsCompatible(isSerialExec).Return(true).Times(1)
					mockRunner.EXPECT().Exec(ctx.Ctx, isSerialExec, serialPrompt).Return(nil).Times(1)
				case 2:
					isSerialExec := testUtils.ExecWithCmd(e.Exec.Cmd)
					mockRunner.EXPECT().IsCompatible(isSerialExec).Return(true).Times(1)
					mockRunner.EXPECT().Exec(ctx.Ctx, isSerialExec, promptedEnv).Return(nil).Times(1)
				}
			}
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

					for i, e := range subExecs {
						switch i {
						case 0:
							isSerialExec := testUtils.ExecWithRef(e.Ref())
							mockCache.EXPECT().GetExecutableByRef(ctx.Logger, e.Ref()).Return(e, nil).Times(1)
							mockRunner.EXPECT().IsCompatible(isSerialExec).Return(true).Times(1)
							mockRunner.EXPECT().Exec(ctx.Ctx, isSerialExec, promptedEnv).Return(nil).Times(1)
						case 1:
							isSerialExec := testUtils.ExecWithRef(e.Ref())
							serialPrompt := map[string]string{"ARG1": "hello", "ARG2": "123"}
							mockCache.EXPECT().GetExecutableByRef(ctx.Logger, e.Ref()).Return(e, nil).Times(1)
							mockRunner.EXPECT().IsCompatible(isSerialExec).Return(true).Times(3)
							mockRunner.EXPECT().Exec(ctx.Ctx, isSerialExec, serialPrompt).Return(errors.New("error")).Times(3)
							mockLogger.EXPECT().Warnx("retrying", "ref", e.Ref()).Times(2)
							mockLogger.EXPECT().
								Errorx("retries exceeded", "err", gomock.Any(), "ref", e.Ref(), "max", 2).
								Times(1)
						case 2:
							isSerialExec := testUtils.ExecWithCmd(e.Exec.Cmd)
							mockRunner.EXPECT().IsCompatible(isSerialExec).Return(true).Times(1)
							mockRunner.EXPECT().Exec(ctx.Ctx, isSerialExec, promptedEnv).Return(nil).Times(1)
						}
					}

					rootExec.Serial.FailFast = false
					Expect(serialRnr.Exec(ctx.Ctx, rootExec, make(map[string]string))).ToNot(Succeed())
				})
			})

			When("fail fast is enabled", func() {
				It("should fail fast after max attempts when enabled", func() {
					mockRunner := ctx.RunnerMock
					mockCache := ctx.ExecutableCache
					mockLogger := ctx.Logger
					promptedEnv := make(map[string]string)

					for i, e := range subExecs {
						switch i {
						case 0:
							isSerialExec := testUtils.ExecWithRef(e.Ref())
							mockCache.EXPECT().GetExecutableByRef(ctx.Logger, e.Ref()).Return(e, nil).Times(1)
							mockRunner.EXPECT().IsCompatible(isSerialExec).Return(true).Times(1)
							mockRunner.EXPECT().Exec(ctx.Ctx, isSerialExec, promptedEnv).Return(nil).Times(1)
						case 1:
							isSerialExec := testUtils.ExecWithRef(e.Ref())
							serialPrompt := map[string]string{"ARG1": "hello", "ARG2": "123"}
							mockCache.EXPECT().GetExecutableByRef(ctx.Logger, e.Ref()).Return(e, nil).Times(1)
							mockRunner.EXPECT().IsCompatible(isSerialExec).Return(true).Times(3)
							mockRunner.EXPECT().Exec(ctx.Ctx, isSerialExec, serialPrompt).Return(errors.New("error")).Times(3)
							mockLogger.EXPECT().Warnx("retrying", "ref", e.Ref()).Times(2)
						}
					}

					rootExec.Serial.FailFast = true
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

					for i, e := range subExecs {
						switch i {
						case 0:
							isSerialExec := testUtils.ExecWithRef(e.Ref())
							mockCache.EXPECT().GetExecutableByRef(ctx.Logger, e.Ref()).Return(e, nil).Times(1)
							mockRunner.EXPECT().IsCompatible(isSerialExec).Return(true).Times(1)
							mockRunner.EXPECT().Exec(ctx.Ctx, isSerialExec, promptedEnv).Return(nil).Times(1)
						case 1:
							isSerialExec := testUtils.ExecWithRef(e.Ref())
							serialPrompt := map[string]string{"ARG1": "hello", "ARG2": "123"}
							mockCache.EXPECT().GetExecutableByRef(ctx.Logger, e.Ref()).Return(e, nil).Times(1)
							mockRunner.EXPECT().IsCompatible(isSerialExec).Return(true).Times(1)
							mockRunner.EXPECT().Exec(ctx.Ctx, isSerialExec, serialPrompt).Return(errors.New("error")).Times(1)
							mockLogger.EXPECT().
								Errorx("execution error", "err", gomock.Any(), "ref", e.Ref()).
								Times(1)
						case 2:
							isSerialExec := testUtils.ExecWithCmd(e.Exec.Cmd)
							mockRunner.EXPECT().IsCompatible(isSerialExec).Return(true).Times(1)
							mockRunner.EXPECT().Exec(ctx.Ctx, isSerialExec, promptedEnv).Return(nil).Times(1)
						}
					}

					Expect(serialRnr.Exec(ctx.Ctx, rootExec, make(map[string]string))).ToNot(Succeed())
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

					for i, e := range subExecs {
						switch i {
						case 0:
							isSerialExec := testUtils.ExecWithRef(e.Ref())
							mockCache.EXPECT().GetExecutableByRef(ctx.Logger, e.Ref()).Return(e, nil).Times(1)
							mockRunner.EXPECT().IsCompatible(isSerialExec).Return(true).Times(1)
							mockRunner.EXPECT().Exec(ctx.Ctx, isSerialExec, promptedEnv).Return(nil).Times(1)
						case 1:
							isSerialExec := testUtils.ExecWithRef(e.Ref())
							serialPrompt := map[string]string{"ARG1": "hello", "ARG2": "123"}
							mockCache.EXPECT().GetExecutableByRef(ctx.Logger, e.Ref()).Return(e, nil).Times(1)
							mockRunner.EXPECT().IsCompatible(isSerialExec).Return(true).Times(1)
							mockRunner.EXPECT().Exec(ctx.Ctx, isSerialExec, serialPrompt).Return(errors.New("error")).Times(1)
						}
					}

					Expect(serialRnr.Exec(ctx.Ctx, rootExec, make(map[string]string))).ToNot(Succeed())
				})
			})
		})
	})
})
