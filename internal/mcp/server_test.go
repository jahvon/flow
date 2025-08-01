package mcp_test

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	flowMcp "github.com/flowexec/flow/internal/mcp"
	"github.com/flowexec/flow/internal/mcp/mocks"
)

func TestServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MCP Server Suite")
}

var _ = Describe("MCP Server", func() {
	var (
		flowServer   *flowMcp.Server
		mockExecutor *mocks.MockCommandExecutor
		mcpClient    *client.Client
		ctx          context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		ctrl := gomock.NewController(GinkgoT())
		mockExecutor = mocks.NewMockCommandExecutor(ctrl)
		flowServer = flowMcp.NewServer(mockExecutor)

		var err error
		mcpClient, err = client.NewInProcessClient(flowServer.GetMCPServer())
		Expect(err).ToNot(HaveOccurred())

		// Initialize the client
		initRequest := mcp.InitializeRequest{
			Params: mcp.InitializeParams{
				ProtocolVersion: mcp.LATEST_PROTOCOL_VERSION,
				ClientInfo: mcp.Implementation{
					Name:    "flow-test-client",
					Version: "1.0.0",
				},
				Capabilities: mcp.ClientCapabilities{},
			},
		}

		_, err = mcpClient.Initialize(ctx, initRequest)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		if mcpClient != nil {
			mcpClient.Close()
		}
	})

	Describe("Server Initialization", func() {
		It("should create server successfully", func() {
			Expect(flowServer).ToNot(BeNil())
			Expect(mcpClient).ToNot(BeNil())
		})
	})

	Describe("Tool Registration", func() {
		It("should register all expected tools", func() {
			toolsResult, err := mcpClient.ListTools(ctx, mcp.ListToolsRequest{})
			Expect(err).ToNot(HaveOccurred())

			toolNames := make([]string, len(toolsResult.Tools))
			for i, tool := range toolsResult.Tools {
				toolNames[i] = tool.Name
			}

			expectedTools := []string{
				"get_info",
				"get_workspace",
				"list_workspaces",
				"switch_workspace",
				"get_executable",
				"list_executables",
				"execute_flow",
				"get_execution_logs",
				"sync_executables",
			}

			for _, expectedTool := range expectedTools {
				Expect(toolNames).To(ContainElement(expectedTool))
			}
		})
	})

	Describe("Prompt Registration", func() {
		It("should register all expected prompts", func() {
			promptsResult, err := mcpClient.ListPrompts(ctx, mcp.ListPromptsRequest{})
			Expect(err).ToNot(HaveOccurred())

			promptNames := make([]string, len(promptsResult.Prompts))
			for i, prompt := range promptsResult.Prompts {
				promptNames[i] = prompt.Name
			}

			expectedPrompts := []string{
				"validate_flowfile",
				"generate_executable",
				"create_workspace",
				"debug_executable",
				"design_workflow",
				"migrate_scripts",
			}

			for _, expectedPrompt := range expectedPrompts {
				Expect(promptNames).To(ContainElement(expectedPrompt))
			}
		})
	})
})
