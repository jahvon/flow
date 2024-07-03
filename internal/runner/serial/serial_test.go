package serial_test

import (
	stdCtx "context"
	"errors"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/builder"
	"github.com/jahvon/flow/internal/runner"
	"github.com/jahvon/flow/internal/runner/serial"
	testRunner "github.com/jahvon/flow/tests/runner"
)

func TestSerialRunner(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Serial Runner Suite")
}

var _ = Describe("SerialRunner", func() {
	var (
		ctx       *testRunner.ContextWithMocks
		serialRnr runner.Runner
	)

	BeforeEach(func() {
		ctx = testRunner.NewContextWithMocks(stdCtx.Background(), GinkgoT())
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

	When("Exec with ref", func() {
		var (
			subject           *config.Executable
			expectedSubExecs  config.ExecutableList
			expectedPromptEnv map[string]string
		)

		BeforeEach(func() {
			subject, expectedSubExecs = builder.SerialExecByRef(ctx.Ctx, "serial-exec", "definition-path")
			ws := ctx.Ctx.CurrentWorkspace
			subject.SetContext(ws.AssignedName(), ws.Location(), "", "")
			expectedPromptEnv = map[string]string{}
		})

		It("should execute in order", func() {
			mockRunner := ctx.RunnerMock
			mockCache := ctx.ExecutableCache

			for _, e := range expectedSubExecs {
				mockRunner.EXPECT().IsCompatible(testRunner.ExecWithRef(e.Ref())).
					Return(true).Times(1)
				mockRunner.EXPECT().Exec(ctx.Ctx, testRunner.ExecWithRef(e.Ref()), expectedPromptEnv).
					Return(nil).Times(1)
				mockCache.EXPECT().GetExecutableByRef(ctx.Logger, e.Ref()).
					Return(e, nil).Times(1)
			}

			Expect(serialRnr.Exec(ctx.Ctx, subject, make(map[string]string))).To(Succeed())
		})

		It("should fail fast when enabled", func() {
			subject.Type.Serial.FailFast = true
			mockRunner := ctx.RunnerMock
			mockCache := ctx.ExecutableCache
			for i, e := range expectedSubExecs {
				if i == 1 {
					mockRunner.EXPECT().IsCompatible(testRunner.ExecWithRef(e.Ref())).
						Return(true).Times(1)
					mockRunner.EXPECT().Exec(ctx.Ctx, testRunner.ExecWithRef(e.Ref()), expectedPromptEnv).
						Return(errors.New("error")).Times(1)
					mockCache.EXPECT().GetExecutableByRef(ctx.Logger, e.Ref()).
						Return(e, nil).Times(1)
					break
				} else {
					mockRunner.EXPECT().IsCompatible(testRunner.ExecWithRef(e.Ref())).
						Return(true).Times(1)
					mockRunner.EXPECT().Exec(ctx.Ctx, testRunner.ExecWithRef(e.Ref()), expectedPromptEnv).
						Return(nil).Times(1)
					mockCache.EXPECT().GetExecutableByRef(ctx.Logger, e.Ref()).
						Return(e, nil).Times(1)
				}
			}

			Expect(serialRnr.Exec(ctx.Ctx, subject, make(map[string]string))).ToNot(Succeed())
		})

		It("should not fail fast when disabled", func() {
			mockRunner := ctx.RunnerMock
			mockCache := ctx.ExecutableCache
			for i, e := range expectedSubExecs {
				if i == 1 {
					By("failing on the second executable")
					mockRunner.EXPECT().IsCompatible(testRunner.ExecWithRef(e.Ref())).
						Return(true).Times(1)
					mockRunner.EXPECT().Exec(ctx.Ctx, testRunner.ExecWithRef(e.Ref()), expectedPromptEnv).
						Return(errors.New("error")).Times(1)
					mockCache.EXPECT().GetExecutableByRef(ctx.Logger, e.Ref()).
						Return(e, nil).Times(1)
					ctx.Logger.EXPECT().Error(gomock.Any(), gomock.Any()).Times(1)
				} else {
					mockRunner.EXPECT().IsCompatible(testRunner.ExecWithRef(e.Ref())).
						Return(true).Times(1)
					mockRunner.EXPECT().Exec(ctx.Ctx, testRunner.ExecWithRef(e.Ref()), expectedPromptEnv).
						Return(nil).Times(1)
					mockCache.EXPECT().GetExecutableByRef(ctx.Logger, e.Ref()).
						Return(e, nil).Times(1)
				}
			}

			Expect(serialRnr.Exec(ctx.Ctx, subject, make(map[string]string))).ToNot(Succeed())
		})
	})

	When("Executables are defined with ref configs", func() {
		var (
			subject          *config.Executable
			expectedSubExecs config.ExecutableList
		)
		BeforeEach(func() {
			subject, expectedSubExecs = builder.SerialExecByRefConfig(ctx.Ctx, "serial-exec", "definition-path")
		})

		It("should execute in order", func() {
			promptedEnv := make(map[string]string)

			mockRunner := ctx.RunnerMock
			mockCache := ctx.ExecutableCache
			for _, e := range expectedSubExecs {
				mockRunner.EXPECT().IsCompatible(testRunner.ExecWithRef(e.Ref())).
					Return(true).Times(1)
				mockRunner.EXPECT().Exec(ctx.Ctx, testRunner.ExecWithRef(e.Ref()), promptedEnv).
					Return(nil).Times(1)
				mockCache.EXPECT().GetExecutableByRef(ctx.Logger, e.Ref()).
					Return(e, nil).Times(1)
			}
			mockRunner.EXPECT().IsCompatible(testRunner.ExecWithCmd(builder.RefConfigCmd)).
				Return(true).Times(1)
			mockRunner.EXPECT().Exec(ctx.Ctx, testRunner.ExecWithCmd(builder.RefConfigCmd), promptedEnv).
				Return(nil).Times(1)

			Expect(serialRnr.Exec(ctx.Ctx, subject, promptedEnv)).To(Succeed())
		})

		Context("when retries are set on a failed ref config", func() {
			BeforeEach(func() {
				subject.Type.Serial.Executables[1].Retries = 2
			})

			It("should be retried until attempted max times", func() {
				mockRunner := ctx.RunnerMock
				mockCache := ctx.ExecutableCache
				for i, e := range expectedSubExecs {
					if i == 1 {
						expectedPromptEnv := map[string]string{"ARG1": "arg1"}
						mockRunner.EXPECT().Exec(ctx.Ctx, testRunner.ExecWithRef(e.Ref()), expectedPromptEnv).
							Return(errors.New("error")).Times(3)
						mockCache.EXPECT().GetExecutableByRef(ctx.Logger, e.Ref()).
							Return(e, nil).Times(3)
						mockRunner.EXPECT().IsCompatible(testRunner.ExecWithRef(e.Ref())).
							Return(true).Times(3)
						ctx.Logger.EXPECT().Warnf(gomock.Any(), gomock.Any()).Times(2)
						ctx.Logger.EXPECT().Error(gomock.Any(), gomock.Any()).Times(1)
					} else {
						expectedPromptEnv := make(map[string]string)
						mockRunner.EXPECT().Exec(ctx.Ctx, testRunner.ExecWithRef(e.Ref()), expectedPromptEnv).
							Return(nil).Times(1)
						mockCache.EXPECT().GetExecutableByRef(ctx.Logger, e.Ref()).
							Return(e, nil).Times(1)
						mockRunner.EXPECT().IsCompatible(testRunner.ExecWithRef(e.Ref())).
							Return(true).Times(1)
					}
				}
				mockRunner.EXPECT().IsCompatible(testRunner.ExecWithCmd(builder.RefConfigCmd)).
					Return(true).Times(1)
				mockRunner.EXPECT().Exec(ctx.Ctx, testRunner.ExecWithCmd(builder.RefConfigCmd), make(map[string]string)).
					Return(nil).Times(1)

				Expect(serialRnr.Exec(ctx.Ctx, subject, make(map[string]string))).To(Succeed())
			})

			It("should not fail fast when disabled", func() {
				mockRunner := ctx.RunnerMock
				// mockCache := ctx.ExecutableCache
				expectedPromptEnv := make(map[string]string)
				for i, e := range expectedSubExecs {
					if i == 1 {
						By("failing on the second executable")
						mockRunner.EXPECT().IsCompatible(testRunner.ExecWithRef(e.Ref())).
							Return(true).Times(2)
						mockRunner.EXPECT().Exec(ctx, testRunner.ExecWithRef(e.Ref()), expectedPromptEnv).
							Return(errors.New("error")).Times(2)
						ctx.Logger.EXPECT().Error(gomock.Any(), gomock.Any()).Times(2)
					} else {
						mockRunner.EXPECT().IsCompatible(testRunner.ExecWithRef(e.Ref())).
							Return(true).Times(1)
						mockRunner.EXPECT().Exec(ctx, testRunner.ExecWithRef(e.Ref()), expectedPromptEnv).
							Return(nil).Times(1)
					}
				}

				Expect(serialRnr.Exec(ctx.Ctx, subject, make(map[string]string))).ToNot(Succeed())
			})
		})

		Context("when retries are not enabled", func() {
			It("should fail fast when enabled", func() {
				mockRunner := ctx.RunnerMock
				expectedPromptEnv := make(map[string]string)
				for i, e := range expectedSubExecs {
					if i == 1 {
						By("failing on the second executable")
						mockRunner.EXPECT().IsCompatible(testRunner.ExecWithRef(e.Ref())).
							Return(true).Times(1)
						mockRunner.EXPECT().Exec(ctx, testRunner.ExecWithRef(e.Ref()), expectedPromptEnv).
							Return(errors.New("error")).Times(1)
					} else {
						mockRunner.EXPECT().IsCompatible(testRunner.ExecWithRef(e.Ref())).
							Return(true).Times(1)
						mockRunner.EXPECT().Exec(ctx, testRunner.ExecWithRef(e.Ref()), expectedPromptEnv).
							Return(nil).Times(1)
					}
				}

				Expect(serialRnr.Exec(ctx.Ctx, subject, make(map[string]string))).ToNot(Succeed())
			})

			It("should not fail fast when disabled", func() {
				mockRunner := ctx.RunnerMock
				expectedPromptEnv := make(map[string]string)
				for i, e := range expectedSubExecs {
					if i == 1 {
						By("failing on the second executable")
						mockRunner.EXPECT().IsCompatible(testRunner.ExecWithRef(e.Ref())).
							Return(true).Times(1)
						mockRunner.EXPECT().Exec(ctx, testRunner.ExecWithRef(e.Ref()), expectedPromptEnv).
							Return(errors.New("error")).Times(1)
						ctx.Logger.EXPECT().Error(gomock.Any(), gomock.Any()).Times(1)
					} else {
						mockRunner.EXPECT().IsCompatible(testRunner.ExecWithRef(e.Ref())).
							Return(true).Times(1)
						mockRunner.EXPECT().Exec(ctx, testRunner.ExecWithRef(e.Ref()), expectedPromptEnv).
							Return(nil).Times(1)
					}
				}

				Expect(serialRnr.Exec(ctx.Ctx, subject, make(map[string]string))).ToNot(Succeed())
			})
		})
	})
})
