package parallel_test

import (
	stdCtx "context"
	"errors"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	"github.com/jahvon/flow/internal/runner"
	"github.com/jahvon/flow/internal/runner/parallel"
	testUtils "github.com/jahvon/flow/tests/utils"
	"github.com/jahvon/flow/tools/builder"
	"github.com/jahvon/flow/types/executable"
)

func TestParallelRunner(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "parallel Runner Suite")
}

var _ = Describe("ParallelRunner", func() {
	var (
		ctx         *testUtils.ContextWithMocks
		parallelRnr runner.Runner
	)

	BeforeEach(func() {
		ctx = testUtils.NewContextWithMocks(stdCtx.Background(), GinkgoT())
		runner.RegisterRunner(ctx.RunnerMock)
		parallelRnr = parallel.NewRunner()
	})

	AfterEach(func() {
		runner.Reset()
	})

	Context("Name", func() {
		It("should return the correct runner name", func() {
			Expect(parallelRnr.Name()).To(Equal("parallel"))
		})
	})

	Context("IsCompatible", func() {
		It("should return false when executable is nil", func() {
			Expect(parallelRnr.IsCompatible(nil)).To(BeFalse())
		})

		It("should return false when executable type is nil", func() {
			executable := &executable.Executable{}
			Expect(parallelRnr.IsCompatible(executable)).To(BeFalse())
		})

		It("should return true when executable type is parallel", func() {
			executable := &executable.Executable{
				Parallel: &executable.ParallelExecutableType{},
			}
			Expect(parallelRnr.IsCompatible(executable)).To(BeTrue())
		})
	})

	When("Exec", func() {
		var (
			rootExec *executable.Executable
			subExecs executable.ExecutableList
		)

		BeforeEach(func() {
			ns := "examples"
			rootExec = builder.ParallelExecByRefConfig(
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

			runner.RegisterRunner(parallelRnr)
			runner.RegisterRunner(ctx.RunnerMock)
			ctx.RunnerMock.EXPECT().IsCompatible(rootExec).Return(false).AnyTimes()
		})

		It("should execute all sub execs", func() {
			promptedEnv := make(map[string]string)
			mockRunner := ctx.RunnerMock
			mockCache := ctx.ExecutableCache

			for i, e := range subExecs {
				switch i {
				case 0:
					isParallelExec := testUtils.ExecWithRef(e.Ref())
					mockCache.EXPECT().GetExecutableByRef(ctx.Logger, e.Ref()).Return(e, nil).Times(1)
					mockRunner.EXPECT().IsCompatible(isParallelExec).Return(true).Times(1)
					mockRunner.EXPECT().Exec(ctx.Ctx, isParallelExec, promptedEnv).Return(nil).Times(1)
				case 1:
					isParallelExec := testUtils.ExecWithRef(e.Ref())
					parallelPrompt := map[string]string{"ARG1": "hello", "ARG2": "123"}
					mockCache.EXPECT().GetExecutableByRef(ctx.Logger, e.Ref()).Return(e, nil).Times(1)
					mockRunner.EXPECT().IsCompatible(isParallelExec).Return(true).Times(1)
					mockRunner.EXPECT().Exec(ctx.Ctx, isParallelExec, parallelPrompt).Return(nil).Times(1)
				case 2:
					isParallelExec := testUtils.ExecWithCmd(e.Exec.Cmd)
					mockRunner.EXPECT().IsCompatible(isParallelExec).Return(true).Times(1)
					mockRunner.EXPECT().Exec(ctx.Ctx, isParallelExec, promptedEnv).Return(nil).Times(1)
				}
			}
			Expect(parallelRnr.Exec(ctx.Ctx, rootExec, promptedEnv)).To(Succeed())
		})

		Context("when retries are set on a failed ref config", func() {
			BeforeEach(func() {
				rootExec.Parallel.Execs[1].Retries = 2
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
							isParallelExec := testUtils.ExecWithRef(e.Ref())
							mockCache.EXPECT().GetExecutableByRef(ctx.Logger, e.Ref()).Return(e, nil).MaxTimes(1)
							mockRunner.EXPECT().IsCompatible(isParallelExec).Return(true).MaxTimes(1)
							mockRunner.EXPECT().
								Exec(ctx.Ctx, isParallelExec, promptedEnv).Return(nil).MaxTimes(1)
						case 1:
							isParallelExec := testUtils.ExecWithRef(e.Ref())
							parallelPrompt := map[string]string{"ARG1": "hello", "ARG2": "123"}
							mockCache.EXPECT().GetExecutableByRef(ctx.Logger, e.Ref()).Return(e, nil).Times(1)
							mockRunner.EXPECT().IsCompatible(isParallelExec).Return(true).Times(3)
							mockRunner.EXPECT().
								Exec(ctx.Ctx, isParallelExec, parallelPrompt).
								Return(errors.New("error")).Times(3)
							mockLogger.EXPECT().Warnx("retrying", "ref", e.Ref()).Times(2)
							mockLogger.EXPECT().
								Errorx("retries exceeded", "err", gomock.Any(), "ref", e.Ref(), "max", 2).
								Times(1)
						case 2:
							isParallelExec := testUtils.ExecWithCmd(e.Exec.Cmd)
							mockRunner.EXPECT().IsCompatible(isParallelExec).Return(true).MaxTimes(1)
							mockRunner.EXPECT().
								Exec(ctx.Ctx, isParallelExec, promptedEnv).Return(nil).MaxTimes(1)
						}
					}

					rootExec.Parallel.FailFast = false
					Expect(parallelRnr.Exec(ctx.Ctx, rootExec, make(map[string]string))).ToNot(Succeed())
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
							isParallelExec := testUtils.ExecWithRef(e.Ref())
							mockCache.EXPECT().GetExecutableByRef(ctx.Logger, e.Ref()).Return(e, nil).MaxTimes(1)
							mockRunner.EXPECT().IsCompatible(isParallelExec).Return(true).MaxTimes(1)
							mockRunner.EXPECT().
								Exec(ctx.Ctx, isParallelExec, promptedEnv).Return(nil).MaxTimes(1)
						case 1:
							isParallelExec := testUtils.ExecWithRef(e.Ref())
							parallelPrompt := map[string]string{"ARG1": "hello", "ARG2": "123"}
							mockCache.EXPECT().GetExecutableByRef(ctx.Logger, e.Ref()).Return(e, nil).Times(1)
							mockRunner.EXPECT().IsCompatible(isParallelExec).Return(true).Times(3)
							mockRunner.EXPECT().
								Exec(ctx.Ctx, isParallelExec, parallelPrompt).
								Return(errors.New("error")).Times(3)
							mockLogger.EXPECT().Warnx("retrying", "ref", e.Ref()).Times(2)
						case 2:
							isParallelExec := testUtils.ExecWithCmd(e.Exec.Cmd)
							mockRunner.EXPECT().IsCompatible(isParallelExec).Return(true).MaxTimes(1)
							mockRunner.EXPECT().
								Exec(ctx.Ctx, isParallelExec, promptedEnv).Return(nil).MaxTimes(1)
						}
					}

					rootExec.Parallel.FailFast = true
					Expect(parallelRnr.Exec(ctx.Ctx, rootExec, make(map[string]string))).ToNot(Succeed())
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
							isParallelExec := testUtils.ExecWithRef(e.Ref())
							mockCache.EXPECT().GetExecutableByRef(ctx.Logger, e.Ref()).Return(e, nil).MaxTimes(1)
							mockRunner.EXPECT().IsCompatible(isParallelExec).Return(true).MaxTimes(1)
							mockRunner.EXPECT().Exec(ctx.Ctx, isParallelExec, promptedEnv).
								Return(nil).MaxTimes(1)
						case 1:
							isParallelExec := testUtils.ExecWithRef(e.Ref())
							parallelPrompt := map[string]string{"ARG1": "hello", "ARG2": "123"}
							mockCache.EXPECT().GetExecutableByRef(ctx.Logger, e.Ref()).Return(e, nil).Times(1)
							mockRunner.EXPECT().IsCompatible(isParallelExec).Return(true).Times(1)
							mockRunner.EXPECT().Exec(ctx.Ctx, isParallelExec, parallelPrompt).
								Return(errors.New("error")).Times(1)
							mockLogger.EXPECT().Errorx("execution error", "err", gomock.Any(), "ref", e.Ref()).
								Times(1)
						case 2:
							isParallelExec := testUtils.ExecWithCmd(e.Exec.Cmd)
							mockRunner.EXPECT().IsCompatible(isParallelExec).Return(true).MaxTimes(1)
							mockRunner.EXPECT().Exec(ctx.Ctx, isParallelExec, promptedEnv).
								Return(nil).MaxTimes(1)
						}
					}

					Expect(parallelRnr.Exec(ctx.Ctx, rootExec, make(map[string]string))).ToNot(Succeed())
				})
			})

			When("fail fast is enabled", func() {
				BeforeEach(func() {
					rootExec.Parallel.FailFast = true
				})

				It("should fail fast after max attempts when enabled", func() {
					mockRunner := ctx.RunnerMock
					mockCache := ctx.ExecutableCache
					promptedEnv := make(map[string]string)

					for i, e := range subExecs {
						switch i {
						case 0:
							isParallelExec := testUtils.ExecWithRef(e.Ref())
							mockCache.EXPECT().GetExecutableByRef(ctx.Logger, e.Ref()).Return(e, nil).MaxTimes(1)
							mockRunner.EXPECT().IsCompatible(isParallelExec).Return(true).MaxTimes(1)
							mockRunner.EXPECT().Exec(ctx.Ctx, isParallelExec, promptedEnv).
								Return(nil).MaxTimes(1)
						case 1:
							isParallelExec := testUtils.ExecWithRef(e.Ref())
							parallelPrompt := map[string]string{"ARG1": "hello", "ARG2": "123"}
							mockCache.EXPECT().GetExecutableByRef(ctx.Logger, e.Ref()).Return(e, nil).Times(1)
							mockRunner.EXPECT().IsCompatible(isParallelExec).Return(true).Times(1)
							mockRunner.EXPECT().Exec(ctx.Ctx, isParallelExec, parallelPrompt).
								Return(errors.New("error")).Times(1)
						case 2:
							isParallelExec := testUtils.ExecWithCmd(e.Exec.Cmd)
							mockRunner.EXPECT().IsCompatible(isParallelExec).Return(true).MaxTimes(1)
							mockRunner.EXPECT().Exec(ctx.Ctx, isParallelExec, promptedEnv).
								Return(nil).MaxTimes(1)
						}
					}

					Expect(parallelRnr.Exec(ctx.Ctx, rootExec, make(map[string]string))).ToNot(Succeed())
				})
			})
		})
	})
})
