package serial_test

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
	"github.com/jahvon/flow/internal/runner/serial"
	testRunner "github.com/jahvon/flow/tests/runner"
	"github.com/jahvon/flow/tools/builder"
	"github.com/jahvon/flow/types/executable"
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

	Context("Title", func() {
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

	Context("Exec", func() {
		var (
			rootExec                                    *executable.Executable
			isSerialExec1, isSerialExec2, isSerialExec3 gomock.Matcher
			mockRunner                                  *mocks.MockRunner
		)

		BeforeEach(func() {
			ctrl := gomock.NewController(GinkgoT())
			mockRunner = mocks.NewMockRunner(ctrl)

			ns := "examples"
			rootExec = builder.SerialExecByRef(
				builder.WithNamespaceName(ns),
				builder.WithWorkspaceName(ctx.CurrentWorkspace.AssignedName()),
				builder.WithWorkspacePath(ctx.CurrentWorkspace.Location()),
			)
			serialSpec := rootExec.Serial
			isSerialExec1 = gomock.Cond(func(e any) bool { return isExecutableWithRef(e, serialSpec.Refs[0]) })
			isSerialExec2 = gomock.Cond(func(e any) bool { return isExecutableWithRef(e, serialSpec.Refs[1]) })
			isSerialExec3 = gomock.Cond(func(e any) bool { return isExecutableWithRef(e, serialSpec.Refs[2]) })

			runner.RegisterRunner(serialRnr)
			runner.RegisterRunner(mockRunner)
			mockRunner.EXPECT().IsCompatible(rootExec).Return(false).AnyTimes()

			Expect(mockLogger).To(Not(BeNil()))
		})

		AfterEach(func() {
			runner.Reset()
		})

		It("should execute in order", func() {
			promptedEnv := make(map[string]string)

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

			rootExec.Serial.FailFast = true
			Expect(serialRnr.Exec(ctx, rootExec, promptedEnv)).ToNot(Succeed())
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

			Expect(serialRnr.Exec(ctx, rootExec, promptedEnv)).ToNot(Succeed())
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
