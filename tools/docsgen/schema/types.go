package schema

import (
	"fmt"
	"strings"
)

type FileName string

func (s FileName) String() string {
	return string(s)
}

func (s FileName) MarkdownFile() string {
	return fmt.Sprintf("%s.md", strings.ToLower(s.Title()))
}

func (s FileName) JSONSchemaFile() string {
	return fmt.Sprintf("%s_schema.json", strings.ToLower(s.Title()))
}

func (s FileName) Title() string {
	parts := strings.Split(s.String(), "/")
	if len(parts) != 2 {
		fmt.Println("invalid schema file name", s.String())
		return ""
	}

	var title string
	if strings.HasSuffix(s.String(), "/schema.yaml") {
		title = parts[0]
	} else if strings.HasSuffix(s.String(), "_schema.yaml") {
		title = strings.TrimSuffix(parts[1], "_schema.yaml")
	}
	return TitleCase(title)
}

type Ref string

func (s Ref) String() string {
	return string(s)
}

func (s Ref) ExternalFile() FileName {
	if s.String() == "" || strings.HasPrefix(s.String(), "#") {
		return ""
	}
	parts := strings.Split(s.String(), "#")
	if len(parts) != 2 {
		fmt.Println("invalid schema ref", s.String())
		return ""
	}
	// remove leading '.' and `/` from file name if present to avoid relative paths
	return FileName(strings.Trim(parts[0], "./"))
}

func (s Ref) DefinitionPath() string {
	if s.String() == "" {
		return ""
	}

	return fmt.Sprintf("#/definitions/%s", s.Key())
}

func (s Ref) IsRoot() bool {
	if s.String() == "" {
		return false
	}

	var def string
	if strings.HasPrefix(s.String(), "#") {
		def = strings.Trim(s.String(), "./")
	} else {
		parts := strings.Split(s.String(), "#")
		if len(parts) != 2 {
			fmt.Println("invalid schema ref", s.String())
			return false
		}
		def = parts[1]
	}
	return def == "/"
}

func (s Ref) Key() FieldKey {
	if s.String() == "" {
		return ""
	}

	var def string
	if strings.HasPrefix(s.String(), "#") {
		def = strings.Trim(s.String(), "./")
	} else {
		parts := strings.Split(s.String(), "#")
		if len(parts) != 2 {
			fmt.Println("invalid schema ref", s.String())
			return ""
		}
		def = parts[1]
	}
	if def == "/" {
		switch {
		case strings.Contains(s.String(), WorkspaceSchema.String()):
			return WorkspaceDefinitionTitle
		case strings.Contains(s.String(), ConfigSchema.String()):
			return ConfigDefinitionTitle
		case strings.Contains(s.String(), FlowfileSchema.String()):
			return FlowfileDefinitionTitle
		case strings.Contains(s.String(), CommonSchema.String()):
			return CommonDefinitionTitle
		case strings.Contains(s.String(), TemplateSchema.String()):
			return TemplateDefinitionTitle
		case strings.Contains(s.String(), ExecutableSchema.String()):
			return ExecutableDefinitionTitle
		}
		fmt.Println("unknown schema ref; defaulting to the file's title", s.String())
		return FieldKey(s.ExternalFile().Title())
	}
	parts := strings.Split(strings.Trim(def, "./"), "/")
	if len(parts) < 2 {
		fmt.Println("invalid schema ref", s.String())
		return ""
	}
	return FieldKey(parts[len(parts)-1])
}

func convertToLocalSchemaRef(original Ref, file FileName) Ref {
	switch {
	case original.String() == "", original.ExternalFile() == file, original.ExternalFile() == "":
		// already a local reference
		return original
	case FieldKey(original.ExternalFile().Title()) == original.Key():
		// likely a root level reference
		return Ref(fmt.Sprintf("#/definitions/%s", original.ExternalFile().Title()))
	case strings.HasPrefix(string(original.Key()), original.ExternalFile().Title()):
		// external reference that already includes the file title
		return Ref(fmt.Sprintf("#/definitions/%s", original.Key()))
	default:
		// external reference should include the file title
		return Ref(fmt.Sprintf("#/definitions/%s%s", original.ExternalFile().Title(), original.Key()))
	}
}

func expandLocalSchemaRef(original Ref, file FileName) Ref {
	switch {
	case original.String() == "", original.ExternalFile() != "":
		return original
	default:
		return Ref(fmt.Sprintf("%s%s", file, original.DefinitionPath()))
	}
}
