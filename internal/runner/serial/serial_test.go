package serial_test

import (
	stdCtx "context"
	"errors"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	"github.com/flowexec/flow/internal/runner"
	"github.com/flowexec/flow/internal/runner/engine"
	"github.com/flowexec/flow/internal/runner/engine/mocks"
	"github.com/flowexec/flow/internal/runner/serial"
	testUtils "github.com/flowexec/flow/tests/utils"
	"github.com/flowexec/flow/tools/builder"
	"github.com/flowexec/flow/types/executable"
)

func TestSerialRunner(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Serial Runner Suite")
}

var _ = Describe("SerialRunner", func() {
	var (
		ctx        *testUtils.ContextWithMocks
		serialRnr  runner.Runner
		mockEngine *mocks.MockEngine
	)

	BeforeEach(func() {
		ctx = testUtils.NewContextWithMocks(stdCtx.Background(), GinkgoTB())
		runner.RegisterRunner(ctx.RunnerMock)
		serialRnr = serial.NewRunner()
		engCtl := gomock.NewController(GinkgoT())
		mockEngine = mocks.NewMockEngine(engCtl)
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

	When("Exec", func() {
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

		It("complete successfully when there are no engine errors", func() {
			promptedEnv := make(map[string]string)
			mockCache := ctx.ExecutableCache

			for i, e := range subExecs {
				switch i {
				case 0:
					mockCache.EXPECT().GetExecutableByRef(e.Ref()).Return(e, nil).Times(1)
				case 1:
					mockCache.EXPECT().GetExecutableByRef(e.Ref()).Return(e, nil).Times(1)
				}
			}
			results := engine.ResultSummary{Results: []engine.Result{{}}}
			mockEngine.EXPECT().Execute(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(results).Times(1)
			Expect(serialRnr.Exec(ctx.Ctx, rootExec, mockEngine, promptedEnv)).To(Succeed())
		})

		It("should fail when there is an engine failure", func() {
			mockCache := ctx.ExecutableCache
			for i, e := range subExecs {
				switch i {
				case 0:
					mockCache.EXPECT().GetExecutableByRef(e.Ref()).Return(e, nil).Times(1)
				case 1:
					mockCache.EXPECT().GetExecutableByRef(e.Ref()).Return(e, nil).Times(1)
				}
			}
			results := engine.ResultSummary{Results: []engine.Result{{Error: errors.New("error")}}}
			mockEngine.EXPECT().Execute(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(results).Times(1)
			Expect(serialRnr.Exec(ctx.Ctx, rootExec, mockEngine, make(map[string]string))).ToNot(Succeed())
		})

		It("should skip execution when condition is false", func() {
			serialSpec := rootExec.Serial
			serialSpec.Execs[0].If = "false"
			serialSpec.Execs[1].If = "true"
			mockCache := ctx.ExecutableCache
			for i, e := range subExecs {
				if i == 1 {
					mockCache.EXPECT().GetExecutableByRef(e.Ref()).Return(e, nil).Times(1)
				}
			}
			results := engine.ResultSummary{Results: []engine.Result{{}}}
			mockEngine.EXPECT().Execute(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(results).Times(1)
			Expect(serialRnr.Exec(ctx.Ctx, rootExec, mockEngine, make(map[string]string))).To(Succeed())
		})
	})
})
