package parallel_test

import (
	stdCtx "context"
	"errors"
	"testing"

	tuikitIOMocks "github.com/jahvon/tuikit/io/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/runner"
	"github.com/jahvon/flow/internal/runner/mocks"
	"github.com/jahvon/flow/internal/runner/parallel"
	testRunner "github.com/jahvon/flow/tests/utils"
	"github.com/jahvon/flow/tools/builder"
	"github.com/jahvon/flow/types/executable"
)

func TestParallelRunner(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "parallel Runner Suite")
}

var _ = Describe("ParallelRunner", func() {
	var (
		ctx         *context.Context
		mockLogger  *tuikitIOMocks.MockLogger
		parallelRnr runner.Runner
	)

	BeforeEach(func() {
		ctrl := gomock.NewController(GinkgoT())
		ctx, mockLogger = testRunner.NewTestContextWithMocks(stdCtx.Background(), GinkgoT(), ctrl)
		parallelRnr = parallel.NewRunner()
	})

	Context("Title", func() {
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

	Context("Exec", func() {
		var (
			rootExec                                          *executable.Executable
			isParallelExec1, isParallelExec2, isParallelExec3 gomock.Matcher
			mockRunner                                        *mocks.MockRunner
		)

		BeforeEach(func() {
			ctrl := gomock.NewController(GinkgoT())
			mockRunner = mocks.NewMockRunner(ctrl)

			ns := "examples"
			rootExec = builder.ParallelExecByRef(
				builder.WithNamespaceName(ns),
				builder.WithWorkspaceName(ctx.CurrentWorkspace.AssignedName()),
				builder.WithWorkspacePath(ctx.CurrentWorkspace.Location()),
				builder.WithFlowFilePath(""),
			)
			parallelSpec := rootExec.Parallel
			isParallelExec1 = gomock.Cond(func(e any) bool { return isExecutableWithRef(e, parallelSpec.Refs[0]) })
			isParallelExec2 = gomock.Cond(func(e any) bool { return isExecutableWithRef(e, parallelSpec.Refs[1]) })
			isParallelExec3 = gomock.Cond(func(e any) bool { return isExecutableWithRef(e, parallelSpec.Refs[2]) })

			runner.RegisterRunner(parallelRnr)
			runner.RegisterRunner(mockRunner)
			mockRunner.EXPECT().IsCompatible(rootExec).Return(false).AnyTimes()

			Expect(mockLogger).To(Not(BeNil()))
		})

		AfterEach(func() {
			runner.Reset()
		})

		It("should execute in order", func() {
			promptedEnv := make(map[string]string)

			mockRunner.EXPECT().IsCompatible(isParallelExec1).Return(true).Times(1)
			mockRunner.EXPECT().Exec(ctx, isParallelExec1, promptedEnv).Return(nil).Times(1)

			mockRunner.EXPECT().IsCompatible(isParallelExec2).Return(true).Times(1)
			mockRunner.EXPECT().Exec(ctx, isParallelExec2, promptedEnv).Return(nil).Times(1)

			mockRunner.EXPECT().IsCompatible(isParallelExec3).Return(true).Times(1)
			mockRunner.EXPECT().Exec(ctx, isParallelExec3, promptedEnv).Return(nil).Times(1)

			rootExec.Parallel.MaxThreads = 1
			Expect(parallelRnr.Exec(ctx, rootExec, promptedEnv)).To(Succeed())

		})

		It("should fail fast when enabled", func() {
			promptedEnv := make(map[string]string)

			mockRunner.EXPECT().IsCompatible(isParallelExec1).Return(true).MaxTimes(1)
			mockRunner.EXPECT().Exec(ctx, isParallelExec1, promptedEnv).Return(nil).MaxTimes(1)

			mockRunner.EXPECT().IsCompatible(isParallelExec2).Return(true).MaxTimes(1)
			mockRunner.EXPECT().Exec(ctx, isParallelExec2, promptedEnv).Return(errors.New("error")).Times(1)

			mockRunner.EXPECT().IsCompatible(isParallelExec3).Return(true).MaxTimes(1)
			mockRunner.EXPECT().Exec(ctx, isParallelExec3, promptedEnv).Return(nil).MaxTimes(1)

			rootExec.Parallel.FailFast = true
			Expect(parallelRnr.Exec(ctx, rootExec, promptedEnv)).ToNot(Succeed())
		})

		It("should not fail fast when disabled", func() {
			promptedEnv := make(map[string]string)

			mockRunner.EXPECT().IsCompatible(isParallelExec1).Return(true).Times(1)
			mockRunner.EXPECT().Exec(ctx, isParallelExec1, promptedEnv).Return(nil).Times(1)

			mockRunner.EXPECT().IsCompatible(isParallelExec2).Return(true).Times(1)
			mockRunner.EXPECT().Exec(ctx, isParallelExec2, promptedEnv).Return(errors.New("error")).Times(1)

			mockRunner.EXPECT().IsCompatible(isParallelExec3).Return(true).Times(1)
			mockRunner.EXPECT().Exec(ctx, isParallelExec3, promptedEnv).Return(nil).Times(1)
			mockLogger.EXPECT().Error(gomock.Any(), gomock.Any()).Times(1)

			Expect(parallelRnr.Exec(ctx, rootExec, promptedEnv)).ToNot(Succeed())
		})
	})
})

func isExecutableWithRef(e any, ref executable.Ref) bool {
	exec, ok := e.(*executable.Executable)
	if !ok {
		return false
	}
	// fmt.Println("want: ", ref, "got: ", exec.Ref())
	return exec.Ref() == ref
}
