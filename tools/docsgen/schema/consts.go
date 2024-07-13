package schema

const (
	TypesRootDir = "types"

	CommonSchema     FileName = "common/schema.yaml"
	WorkspaceSchema  FileName = "workspace/schema.yaml"
	ConfigSchema     FileName = "config/schema.yaml"
	ExecutableSchema FileName = "executable/executable_schema.yaml"
	FlowfileSchema   FileName = "executable/flowfile_schema.yaml"

	CommonDefinitionTitle     = "Common"
	WorkspaceDefinitionTitle  = "Workspace"
	ConfigDefinitionTitle     = "Config"
	ExecutableDefinitionTitle = "Executable"
	FlowfileDefinitionTitle   = "FlowFile"
)

var (
	SchemaFilesForDocs = []FileName{
		CommonSchema,
		WorkspaceSchema,
		ConfigSchema,
		ExecutableSchema,
		FlowfileSchema,
	}
)
