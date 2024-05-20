package runner_test

import (
	"os"

	"github.com/jahvon/tuikit/io/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/runner"
)

var _ = Describe("Env", func() {
	var (
		ctrl   *gomock.Controller
		logger *mocks.MockLogger
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		logger = mocks.NewMockLogger(ctrl)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("SetEnv", func() {
		It("should set environment variables correctly", func() {
			exec := &config.ExecutableEnvironment{
				Parameters: []config.Parameter{
					{
						EnvKey: "TEST_KEY",
						Text:   "test",
					},
				},
			}
			promptedEnv := make(map[string]string)
			err := runner.SetEnv(logger, exec, promptedEnv)
			Expect(err).ToNot(HaveOccurred())
			val, exists := os.LookupEnv("TEST_KEY")
			Expect(exists).To(BeTrue())
			Expect(val).To(Equal("test"))
		})
	})

	Describe("ResolveParameterValue", func() {
		It("should return empty string when all parameter fields are empty", func() {
			param := config.Parameter{}
			promptedEnv := make(map[string]string)
			val, err := runner.ResolveParameterValue(logger, param, promptedEnv)
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal(""))
		})

		It("should return text when text field is not empty", func() {
			param := config.Parameter{
				Text: "test",
			}
			promptedEnv := make(map[string]string)
			val, err := runner.ResolveParameterValue(logger, param, promptedEnv)
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal("test"))
		})

		It("should return error when prompt field is not empty but env key is not in promptedEnv", func() {
			param := config.Parameter{
				Prompt: "test",
				EnvKey: "TEST_KEY",
			}
			promptedEnv := make(map[string]string)
			_, err := runner.ResolveParameterValue(logger, param, promptedEnv)
			Expect(err).To(HaveOccurred())
		})

		It("should return value from promptedEnv when prompt field is not empty and env key is in promptedEnv", func() {
			param := config.Parameter{
				Prompt: "test",
				EnvKey: "TEST_KEY",
			}
			promptedEnv := map[string]string{
				"TEST_KEY": "test",
			}
			val, err := runner.ResolveParameterValue(logger, param, promptedEnv)
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal("test"))
		})

		// TODO: Add test cases for SecretRef
	})

	Describe("BuildEnvList", func() {
		It("should build the environment list correctly", func() {
			exec := &config.ExecutableEnvironment{
				Parameters: []config.Parameter{
					{
						EnvKey: "TEST_KEY",
						Text:   "test",
					},
					{
						EnvKey: "TEST_KEY_2",
						Text:   "test2",
					},
				},
			}
			inputEnv := make(map[string]string)
			defaultEnv := make(map[string]string)
			envList, err := runner.BuildEnvList(logger, exec, inputEnv, defaultEnv)
			Expect(err).ToNot(HaveOccurred())
			Expect(envList).To(Equal([]string{"TEST_KEY=test", "TEST_KEY_2=test2"}))
		})
	})

	Describe("BuildEnvMap", func() {
		It("should build environment map correctly", func() {
			exec := &config.ExecutableEnvironment{
				Parameters: []config.Parameter{
					{
						EnvKey: "TEST_KEY",
						Text:   "test",
					},
					{
						EnvKey: "TEST_KEY_2",
						Text:   "test2",
					},
				},
			}
			inputEnv := make(map[string]string)
			defaultEnv := make(map[string]string)
			envMap, err := runner.BuildEnvMap(logger, exec, inputEnv, defaultEnv)
			Expect(err).ToNot(HaveOccurred())
			Expect(envMap).To(Equal(map[string]string{"TEST_KEY": "test", "TEST_KEY_2": "test2"}))
		})
	})

	Describe("DefaultEnv", func() {
		It("should return default environment correctly", func() {
			wsName := "test-workspace"
			wsLocation := "test-location"
			nsName := "test-namespace"
			execName := "test-executable"
			execDefinitionLocation := "test-definition-location"

			ws := config.WorkspaceConfig{}
			ws.SetContext(wsName, wsLocation)
			ctx := &context.Context{
				CurrentWorkspace: &ws,
				UserConfig: &config.UserConfig{
					CurrentNamespace: nsName,
				},
			}
			exec := config.Executable{Name: execName}
			exec.SetContext(wsName, wsLocation, nsName, execDefinitionLocation)
			envMap := runner.DefaultEnv(ctx, &exec)
			Expect(envMap["FLOW_RUNNER"]).To(Equal("true"))
			Expect(envMap["FLOW_CURRENT_WORKSPACE"]).To(Equal(wsName))
			Expect(envMap["FLOW_CURRENT_NAMESPACE"]).To(Equal(nsName))
			Expect(envMap["FLOW_EXECUTABLE_NAME"]).To(Equal(execName))
			Expect(envMap["DISABLE_FLOW_INTERACTIVE"]).To(Equal("true"))
			// TODO: Add more assertions for other keys in the environment map
		})
	})
})
