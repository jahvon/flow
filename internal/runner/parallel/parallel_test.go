package parallel_test

import (
	stdCtx "context"
	"errors"
	"testing"

	tuikitIOMocks "github.com/jahvon/tuikit/io/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/runner"
	"github.com/jahvon/flow/internal/runner/mocks"
	"github.com/jahvon/flow/internal/runner/parallel"
	examples_test "github.com/jahvon/flow/tests/examples"
	testRunner "github.com/jahvon/flow/tests/runner"
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
		ctx, mockLogger = testRunner.NewTestContextWithMockLogger(stdCtx.Background(), GinkgoT(), ctrl)
		parallelRnr = parallel.NewRunner()
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
			executable := &config.Executable{}
			Expect(parallelRnr.IsCompatible(executable)).To(BeFalse())
		})

		It("should return true when executable type is parallel", func() {
			executable := &config.Executable{
				Type: &config.ExecutableTypeSpec{
					Parallel: &config.ParallelExecutableType{},
				},
			}
			Expect(parallelRnr.IsCompatible(executable)).To(BeTrue())
		})
	})

	Context("Exec", func() {
		var (
			rootExec, parallelExec1, parallelExec2, parallelExec3 *config.Executable
			isParallelExec1, isParallelExec2, isParallelExec3     gomock.Matcher
			mockRunner                                            *mocks.MockRunner
		)

		BeforeEach(func() {
			ctrl := gomock.NewController(GinkgoT())
			mockRunner = mocks.NewMockRunner(ctrl)
			// waitGroup = &sync.WaitGroup{}

			setExecCtx := func(exec *config.Executable) {
				exec.SetContext(
					ctx.CurrentWorkspace.AssignedName(),
					ctx.CurrentWorkspace.Location(),
					"examples",
					"",
				)
			}

			rootExec = examples_test.ParallelExecRoot
			parallelExec1 = examples_test.ParallelExec1
			setExecCtx(parallelExec1)
			parallelExec2 = examples_test.ParallelExec2
			setExecCtx(parallelExec2)
			parallelExec3 = examples_test.ParallelExec3
			setExecCtx(parallelExec3)

			isParallelExec1 = gomock.Cond(func(e any) bool { return isExecutableWithRef(e, parallelExec1.Ref()) })
			isParallelExec2 = gomock.Cond(func(e any) bool { return isExecutableWithRef(e, parallelExec2.Ref()) })
			isParallelExec3 = gomock.Cond(func(e any) bool { return isExecutableWithRef(e, parallelExec3.Ref()) })

			runner.RegisterRunner(parallelRnr)
			runner.RegisterRunner(mockRunner)

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

			Expect(parallelRnr.Exec(ctx, examples_test.ParallelExecRootWithMaxThreads, promptedEnv)).To(Succeed())

		})

		It("should fail fast when enabled", func() {
			promptedEnv := make(map[string]string)

			mockRunner.EXPECT().IsCompatible(gomock.Any()).Return(true).MinTimes(2).MaxTimes(4)
			mockRunner.EXPECT().Exec(ctx, gomock.Any(), promptedEnv).Return(nil).MaxTimes(3)
			mockRunner.EXPECT().Exec(ctx, gomock.Any(), promptedEnv).Return(errors.New("error")).Times(1)

			Expect(parallelRnr.Exec(ctx, examples_test.ParallelExecRootWithExit, promptedEnv)).ToNot(Succeed())
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

func isExecutableWithRef(e any, ref config.Ref) bool {
	exec, ok := e.(*config.Executable)
	if !ok {
		return false
	}
	// fmt.Println("want: ", ref, "got: ", exec.Ref())
	return exec.Ref() == ref
}
