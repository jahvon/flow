package serial_test

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
	"github.com/jahvon/flow/internal/runner/serial"
	examples_test "github.com/jahvon/flow/tests/examples"
	testRunner "github.com/jahvon/flow/tests/runner"
)

func TestSerialRunner(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Serial Runner Suite")
}

var _ = Describe("SerialRunner", func() {
	var (
		ctx        *context.Context
		mockLogger *tuikitIOMocks.MockLogger
		serialRnr  runner.Runner
	)

	BeforeEach(func() {
		ctrl := gomock.NewController(GinkgoT())
		ctx, mockLogger = testRunner.NewTestContextWithMockLogger(stdCtx.Background(), GinkgoT(), ctrl)
		serialRnr = serial.NewRunner()
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
			executable := &config.Executable{}
			Expect(serialRnr.IsCompatible(executable)).To(BeFalse())
		})

		It("should return true when executable type is serial", func() {
			executable := &config.Executable{
				Type: &config.ExecutableTypeSpec{
					Serial: &config.SerialExecutableType{},
				},
			}
			Expect(serialRnr.IsCompatible(executable)).To(BeTrue())
		})
	})

	Context("Exec", func() {
		var (
			rootExec, serialExec1, serialExec2, serialExec3 *config.Executable
			isSerialExec1, isSerialExec2, isSerialExec3     gomock.Matcher
			mockRunner                                      *mocks.MockRunner
		)

		BeforeEach(func() {
			ctrl := gomock.NewController(GinkgoT())
			mockRunner = mocks.NewMockRunner(ctrl)

			setExecCtx := func(exec *config.Executable) {
				exec.SetContext(
					ctx.CurrentWorkspace.AssignedName(),
					ctx.CurrentWorkspace.Location(),
					"examples",
					"",
				)
			}

			rootExec = examples_test.SerialExecRoot
			serialExec1 = examples_test.SerialExec1
			setExecCtx(serialExec1)
			serialExec2 = examples_test.SerialExec2
			setExecCtx(serialExec2)
			serialExec3 = examples_test.SerialExec3
			setExecCtx(serialExec3)

			isSerialExec1 = gomock.Cond(func(e any) bool { return isExecutableWithRef(e, serialExec1.Ref()) })
			isSerialExec2 = gomock.Cond(func(e any) bool { return isExecutableWithRef(e, serialExec2.Ref()) })
			isSerialExec3 = gomock.Cond(func(e any) bool { return isExecutableWithRef(e, serialExec3.Ref()) })

			runner.RegisterRunner(serialRnr)
			runner.RegisterRunner(mockRunner)

			Expect(mockLogger).To(Not(BeNil()))
		})

		AfterEach(func() {
			runner.Reset()
		})

		It("should execute in order", func() {
			promptedEnv := make(map[string]string)

			// mockRunner.EXPECT().IsCompatible(rootExec).Return(true).Times(1)
			// mockRunner.EXPECT().Exec(ctx, rootExec, make(map[string]string)).Return(nil).Times(1)

			// isRootExec := gomock.Cond(func(e any) bool { return isExecutableWithRef(e, "run examples:serial") })

			mockRunner.EXPECT().IsCompatible(isSerialExec1).Return(true).Times(1)
			mockRunner.EXPECT().Exec(ctx, isSerialExec1, promptedEnv).Return(nil).Times(1)

			mockRunner.EXPECT().IsCompatible(isSerialExec2).Return(true).Times(1)
			mockRunner.EXPECT().Exec(ctx, isSerialExec2, promptedEnv).Return(nil).Times(1)

			mockRunner.EXPECT().IsCompatible(isSerialExec3).Return(true).Times(1)
			mockRunner.EXPECT().Exec(ctx, isSerialExec3, promptedEnv).Return(nil).Times(1)

			Expect(serialRnr.Exec(ctx, rootExec, promptedEnv)).To(Succeed())
		})

		It("should fail fast when enabled", func() {
			promptedEnv := make(map[string]string)

			mockRunner.EXPECT().IsCompatible(isSerialExec1).Return(true).Times(1)
			mockRunner.EXPECT().Exec(ctx, isSerialExec1, promptedEnv).Return(nil).Times(1)

			mockRunner.EXPECT().IsCompatible(isSerialExec2).Return(true).Times(1)
			mockRunner.EXPECT().Exec(ctx, isSerialExec2, promptedEnv).Return(errors.New("error")).Times(1)

			Expect(serialRnr.Exec(ctx, examples_test.SerialWithExitRoot, promptedEnv)).ToNot(Succeed())
		})

		It("should not fail fast when disabled", func() {
			promptedEnv := make(map[string]string)

			mockRunner.EXPECT().IsCompatible(isSerialExec1).Return(true).Times(1)
			mockRunner.EXPECT().Exec(ctx, isSerialExec1, promptedEnv).Return(nil).Times(1)

			mockRunner.EXPECT().IsCompatible(isSerialExec2).Return(true).Times(1)
			mockRunner.EXPECT().Exec(ctx, isSerialExec2, promptedEnv).Return(errors.New("error")).Times(1)

			mockRunner.EXPECT().IsCompatible(isSerialExec3).Return(true).Times(1)
			mockRunner.EXPECT().Exec(ctx, isSerialExec3, promptedEnv).Return(nil).Times(1)
			mockLogger.EXPECT().Error(gomock.Any(), gomock.Any()).Times(1)

			Expect(serialRnr.Exec(ctx, examples_test.SerialExecRoot, promptedEnv)).ToNot(Succeed())
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
