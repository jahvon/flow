package mcp

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func addServerPrompts(s *server.MCPServer) {
	validateFlowfile := mcp.NewPrompt("validate_flowfile",
		mcp.WithPromptDescription("Validate a FlowFile configuration"),
		mcp.WithArgument("flowfile_content",
			mcp.ArgumentDescription("FlowFile YAML content to validate"),
			mcp.RequiredArgument(),
		),
		mcp.WithArgument("flowfile_path",
			mcp.ArgumentDescription("Path to the FlowFile (optional)"),
		),
		mcp.WithArgument("strict_mode",
			mcp.ArgumentDescription("Enable strict validation with best practices"),
		),
	)
	s.AddPrompt(validateFlowfile, validateFlowfilePrompt)

	generateExecutable := mcp.NewPrompt("generate_executable",
		mcp.WithPromptDescription("Generate a Flow executable configuration"),
		mcp.WithArgument("verb",
			mcp.ArgumentDescription("Action verb (build, test, deploy, etc.)"),
			mcp.RequiredArgument(),
		),
		mcp.WithArgument("name", mcp.ArgumentDescription("Executable name")),
		mcp.WithArgument("type", mcp.ArgumentDescription("Executable type")),
		mcp.WithArgument("description", mcp.ArgumentDescription("What this executable does")),
		mcp.WithArgument("command", mcp.ArgumentDescription("Command to execute (for exec type)")),
	)
	s.AddPrompt(generateExecutable, generateExecutablePrompt)

	createWorkspace := mcp.NewPrompt("create_workspace",
		mcp.WithPromptDescription("Generate workspace setup and initial flow files"),
		mcp.WithArgument("name",
			mcp.ArgumentDescription("Workspace name"),
			mcp.RequiredArgument(),
		),
		mcp.WithArgument("type",
			mcp.ArgumentDescription("Project type (web, api, mobile, cli, etc.)"),
		),
		mcp.WithArgument("tech_stack",
			mcp.ArgumentDescription("Technology stack (node, go, python, etc.)"),
		),
	)
	s.AddPrompt(createWorkspace, createWorkspacePrompt)

	debugExecutable := mcp.NewPrompt("debug_executable",
		mcp.WithPromptDescription("Help debug Flow executable issues"),
		mcp.WithArgument("executable_ref",
			mcp.ArgumentDescription("Executable reference that's having issues"),
			mcp.RequiredArgument(),
		),
		mcp.WithArgument("error_message",
			mcp.ArgumentDescription("Error message or description of the problem"),
		),
		mcp.WithArgument("expected_behavior",
			mcp.ArgumentDescription("What you expected to happen"),
		),
	)
	s.AddPrompt(debugExecutable, debugExecutablePrompt)

	designWorkflow := mcp.NewPrompt("design_workflow",
		mcp.WithPromptDescription("Design a multi-step workflow using serial/parallel executables"),
		mcp.WithArgument("goal",
			mcp.ArgumentDescription("What the workflow should accomplish"),
			mcp.RequiredArgument(),
		),
		mcp.WithArgument("steps",
			mcp.ArgumentDescription("List of steps or tasks to include"),
			mcp.RequiredArgument(),
		),
		mcp.WithArgument("execution_style", mcp.ArgumentDescription("Execution style preference")),
	)
	s.AddPrompt(designWorkflow, designWorkflowPrompt)

	migrateWorkflows := mcp.NewPrompt("migrate_scripts",
		mcp.WithPromptDescription("Help migrate existing scripts/tools to Flow executables"),
		mcp.WithArgument("current_approach",
			mcp.ArgumentDescription("How you currently run these tasks"),
			mcp.RequiredArgument(),
		),
		mcp.WithArgument("scripts_description",
			mcp.ArgumentDescription("Description of existing scripts or commands"),
		),
	)
	s.AddPrompt(migrateWorkflows, migrateScriptsPrompt)
}

func validateFlowfilePrompt(_ context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	args := req.Params.Arguments
	flowfileContent := getArgOrDefault(args, "flowfile_content", "")
	filePath := getArgOrDefault(args, "file_path", "unknown.flow")
	strictMode := getArgOrDefault(args, "strict_mode", "true")

	prompt := fmt.Sprintf(`Please validate this Flow file configuration:

**File Path**: %s
**Strict Validation**: %s

**Flow File Content:**
"""yaml
%s
"""

Please perform a comprehensive validation and provide:

1. **Syntax Validation**:
   - YAML syntax correctness
   - Required fields present
   - Field types and formats
   - Schema compliance

2. **Semantic Validation**:
   - Executable references are valid
   - Verb and name combinations are unique
   - Parameter and argument configurations make sense
   - Conditional expressions are valid (if present)

3. **Best Practice Review** (if strict_mode enabled):
   - Naming conventions
   - Description quality and completeness
   - Proper use of namespaces and tags
   - Security considerations (no hardcoded secrets)
   - Performance implications (timeouts, retries)

4. **Potential Issues**:
   - Common configuration mistakes
   - Missing recommended fields
   - Conflicting settings
   - Environment-specific problems

5. **Recommendations**:
   - Suggested improvements
   - Alternative approaches
   - Additional fields that would be helpful
   - Integration opportunities

6. **Corrected Configuration** (if issues found):
   - Provide a corrected version of the YAML
   - Highlight changes made
   - Explain reasoning for corrections

If the configuration is valid, confirm this and highlight any optional improvements. If invalid, provide specific error locations and clear fix instructions.`,
		filePath, strictMode, flowfileContent)

	return mcp.NewGetPromptResult(
		"Validate Flow File",
		[]mcp.PromptMessage{
			mcp.NewPromptMessage(
				mcp.RoleUser,
				mcp.NewTextContent(prompt),
			),
		},
	), nil
}

func generateExecutablePrompt(_ context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	args := request.Params.Arguments
	verb := getArgOrDefault(args, "verb", "exec")
	name := getArgOrDefault(args, "name", "my-task")
	execType := getArgOrDefault(args, "type", "exec")
	description := getArgOrDefault(args, "description", "")
	command := getArgOrDefault(args, "command", "")

	prompt := fmt.Sprintf(`I want to create a Flow executable with these requirements:

**Executable Details:**
- Verb: %s
- Name: %s (exclude the name if the verb alone is a sufficient identifier)
- Type: %s
- Description: %s
- Command: %s

Please generate a complete Flow executable configuration in YAML format. Include:

1. **Basic Structure**: Valid FlowFile YAML syntax with verb, name, and type
2. **Parameters**: Any environment variables or secrets that might be needed
3. **Arguments**: Command-line arguments if applicable
4. **Error Handling**: Appropriate timeout and retry configuration (for serial/parallel)
5. **Documentation**: Clear description and any usage notes
6. **Best Practices**: Follow Flow conventions and patterns

The executable should be production-ready and include proper error handling. If you need more information about the specific use case, please ask clarifying questions.`,
		verb, name, execType, description, command)

	return mcp.NewGetPromptResult(
		"Generate Flow Executable",
		[]mcp.PromptMessage{
			mcp.NewPromptMessage(
				mcp.RoleUser,
				mcp.NewTextContent(prompt),
			),
		},
	), nil
}

func createWorkspacePrompt(_ context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	args := request.Params.Arguments
	name := getArgOrDefault(args, "name", "my-project")
	projectType := getArgOrDefault(args, "type", "web")
	techStack := getArgOrDefault(args, "tech_stack", "")
	prompt := fmt.Sprintf(`I want to set up a new Flow workspace for a %s project:

**Project Details:**
- Workspace Name: %s
- Project Type: %s
- Technology Stack: %s

Please help me create:

1. **Workspace Configuration** (flow.yaml):
   - Appropriate workspace settings
   - Description and display name
   - Any necessary executable filters

2. **Initial Flow Files**:
   - Common executables for this project type
   - Proper organization with namespaces if needed
   - Development, testing, and deployment workflows

3. **Directory Structure**:
   - Recommended file organization
   - Where to place different types of executables
   - Integration with existing project structure

4. **Getting Started Guide**:
   - How to register the workspace with Flow
   - Essential first commands to run
   - Next steps for customization

Focus on creating a practical, immediately useful setup that follows Flow and coding best practices for %s projects.`,
		projectType, name, projectType, techStack, projectType)

	return mcp.NewGetPromptResult(
		"Create Flow Workspace",
		[]mcp.PromptMessage{
			mcp.NewPromptMessage(
				mcp.RoleUser,
				mcp.NewTextContent(prompt),
			),
		},
	), nil
}

func debugExecutablePrompt(_ context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	args := req.Params.Arguments
	execRef := getArgOrDefault(args, "executable_ref", "")
	errorMessage := getArgOrDefault(args, "error_message", "")
	expectedBehavior := getArgOrDefault(args, "expected_behavior", "")

	prompt := fmt.Sprintf(`I'm having trouble with a Flow executable and need debugging help:

**Executable Reference**: %s
**Error/Issue**: %s
**Expected Behavior**: %s

Please help me troubleshoot this issue by:

1. **Analyzing the Problem**:
   - Identify potential causes based on the error
   - Consider common Flow executable issues
   - Check for configuration problems

2. **Diagnostic Steps**:
   - Commands to run to gather more information (including the flow logs command)
   - How to check executable configuration
   - Ways to test components individually

3. **Solution Recommendations**:
   - Specific fixes for the identified issues
   - Configuration adjustments needed
   - Best practices to prevent similar issues

4. **Verification**:
   - How to test that the fix works
   - Commands to validate the executable
   - Signs that everything is working correctly

Include specific Flow CLI commands and configuration examples in your response.`,
		execRef, errorMessage, expectedBehavior)

	return mcp.NewGetPromptResult(
		"Debug Flow Executable",
		[]mcp.PromptMessage{
			mcp.NewPromptMessage(
				mcp.RoleUser,
				mcp.NewTextContent(prompt),
			),
		},
	), nil
}

func designWorkflowPrompt(_ context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	args := req.Params.Arguments
	goal := getArgOrDefault(args, "goal", "")
	steps := getArgOrDefault(args, "steps", "")
	executionStyle := getArgOrDefault(args, "execution_style", "serial")

	prompt := fmt.Sprintf(`I need to design a Flow workflow with these requirements:

**Workflow Goal**: %s
**Required Steps**: %s
**Preferred Execution**: %s

Please design a complete Flow workflow that:

1. **Workflow Structure**:
   - Use appropriate serial/parallel execution
   - Proper step ordering and dependencies
   - Error handling and fail-fast configuration

2. **Step Implementation**:
   - Break down complex steps into manageable executables
   - Define clear inputs/outputs for each step
   - Include necessary parameters and environment variables

3. **Integration Points**:
   - How steps share data and state
   - Conditional execution based on previous results
   - Rollback or cleanup procedures if needed

4. **Configuration**:
   - Complete YAML configuration
   - Parameter passing between steps
   - Timeout and retry strategies

5. **Usage Examples**:
   - How to run the complete workflow
   - How to run individual steps for testing
   - Common variations and options

Make the workflow robust, maintainable, and easy to understand.`,
		goal, steps, executionStyle)

	return mcp.NewGetPromptResult(
		"Design Flow Workflow",
		[]mcp.PromptMessage{
			mcp.NewPromptMessage(
				mcp.RoleUser,
				mcp.NewTextContent(prompt),
			),
		},
	), nil
}

func migrateScriptsPrompt(_ context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	args := req.Params.Arguments
	currentApproach := getArgOrDefault(args, "current_approach", "")
	scriptsDescription := getArgOrDefault(args, "scripts_description", "")

	prompt := fmt.Sprintf(`I want to migrate my existing automation to Flow:

**Current Approach**: %s
**Scripts/Tools Description**: %s

Please help me migrate to Flow by providing:

1. **Migration Strategy**:
   - How to organize existing scripts in Flow
   - Workspace and namespace structure recommendations
   - Migration timeline and approach

2. **Flow Equivalents**:
   - Convert existing scripts to Flow executables
   - Preserve existing functionality and behavior
   - Improve error handling and logging

3. **Enhanced Features**:
   - Take advantage of Flow's parameter system
   - Add proper documentation and descriptions
   - Implement conditional logic where beneficial

4. **Integration Plan**:
   - How to gradually migrate from old to new system
   - Maintain compatibility during transition
   - Testing strategy for migrated workflows

5. **Best Practices**:
   - Flow-specific improvements over current approach
   - Better organization and discoverability
   - Maintenance and sharing considerations

Provide concrete examples and step-by-step migration instructions.`,
		currentApproach, scriptsDescription)

	return mcp.NewGetPromptResult(
		"Migrate Scripts to Flow",
		[]mcp.PromptMessage{
			mcp.NewPromptMessage(
				mcp.RoleUser,
				mcp.NewTextContent(prompt),
			),
		},
	), nil
}

func getArgOrDefault(args map[string]string, key, defaultVal string) string {
	if val, exists := args[key]; exists {
		return val
	}
	return defaultVal
}
