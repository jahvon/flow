package schema

import (
	"strings"
)

type FieldKey string

func (f FieldKey) String() string {
	return string(f)
}

func (f FieldKey) Lower() string {
	return FieldCase(string(f))
}

func (f FieldKey) Title() string {
	return TitleCase(string(f))
}

//nolint:lll
type JSONSchema struct {
	Schema               string                   `json:"$schema,omitempty"              yaml:"$schema,omitempty"`
	Ref                  Ref                      `json:"$ref,omitempty"                 yaml:"$ref,omitempty"`
	ID                   string                   `json:"$id,omitempty"                  yaml:"$id,omitempty"`
	Description          string                   `json:"description,omitempty"          yaml:"description,omitempty"`
	Type                 string                   `json:"type,omitempty"                 yaml:"type,omitempty"`
	Required             []string                 `json:"required,omitempty"             yaml:"required,omitempty"`
	Default              interface{}              `json:"default,omitempty"              yaml:"default,omitempty"`
	Enum                 []string                 `json:"enum,omitempty"                 yaml:"enum,omitempty"`
	Definitions          map[FieldKey]*JSONSchema `json:"definitions,omitempty"          yaml:"definitions,omitempty"`
	Properties           map[FieldKey]*JSONSchema `json:"properties,omitempty"           yaml:"properties,omitempty"`
	AdditionalProperties *JSONSchema              `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`
	Items                *JSONSchema              `json:"items,omitempty"                yaml:"items,omitempty"`
	Ext                  SchemaExt                `json:"-"                              yaml:"goJSONSchema,omitempty"`
}

type SchemaExt struct {
	Type       string `json:"-" yaml:"type"`
	Identifier string `json:"-" yaml:"identifier"`
}

func (e SchemaExt) IsExported() bool {
	return e.Identifier == "" || (e.Identifier[0] >= 'A' && e.Identifier[0] <= 'Z')
}

//nolint:gocognit,nestif
func MergeSchemas(dst, src *JSONSchema, dstFile FileName, schemaMap map[FileName]*JSONSchema) {
	if src.Items != nil {
		MergeSchemas(dst, src.Items, dstFile, schemaMap)
	}

	var match *JSONSchema
	switch {
	case src.Ref.String() == "":
		// the source is not a reference
		for key, value := range src.Properties {
			if !value.Ext.IsExported() {
				delete(src.Properties, key)
				continue
			}
			MergeSchemas(dst, value, dstFile, schemaMap)
		}
	case src.Ref.ExternalFile() == "" && dst.Definitions[src.Ref.Key()] == nil:
		// the ref is a local definition but doesn't exist in the destination schema
		// should never happen if defined correctly
		match = &JSONSchema{
			Type:        src.Type,
			Description: src.Description,
			Required:    src.Required,
			Default:     src.Default,
			Enum:        src.Enum,
		}
	case src.Ref.ExternalFile() == "":
		// the ref is a local definition and exists in the destination schema
		for key, value := range src.Properties {
			if !value.Ext.IsExported() {
				delete(src.Properties, key)
				continue
			}
			MergeSchemas(dst, value, dstFile, schemaMap)
		}
	default:
		// the ref is an external definition
		for fn, s := range schemaMap {
			if src.Ref.ExternalFile().Title() == fn.Title() {
				if FieldKey(fn.Title()) == src.Ref.Key() {
					// root level reference
					match = s
					match.ID = ""
				} else {
					def, found := s.Definitions[src.Ref.Key()]
					if !found {
						continue
					}
					match = def
				}
				if match.Items != nil {
					match.Items.Ref = expandLocalSchemaRef(match.Items.Ref, fn)
					MergeSchemas(dst, match.Items, dstFile, schemaMap)
					match.Items.Ref = convertToLocalSchemaRef(match.Items.Ref, dstFile)
				}
				for _, value := range match.Properties {
					if !value.Ext.IsExported() {
						continue
					}
					value.Ref = expandLocalSchemaRef(value.Ref, fn)
					MergeSchemas(dst, value, dstFile, schemaMap)
				}
				break
			}
		}
	}

	if match == nil {
		return
	}

	if src.Ref.String() != "" {
		src.Ref = convertToLocalSchemaRef(src.Ref, dstFile)
	}
	for key, value := range match.Definitions {
		MergeSchemas(dst, value, dstFile, schemaMap)
		delete(match.Definitions, key)
	}
	if _, found := dst.Definitions[src.Ref.Key()]; !found {
		dst.Definitions[src.Ref.Key()] = match
	}
	for key, value := range match.Properties {
		if !value.Ext.IsExported() {
			delete(match.Properties, key)
			continue
		}
		MergeSchemas(dst, value, dstFile, schemaMap)
	}
	for _, value := range match.Definitions {
		MergeSchemas(dst, value, dstFile, schemaMap)
	}
}

func FieldCase(s string) string {
	if s == "" {
		return ""
	}
	firstLetter := strings.ToLower(s[:1])
	return firstLetter + s[1:]
}

func TitleCase(s string) string {
	if s == "" {
		return ""
	} else if s == strings.ToLower(FlowfileDefinitionTitle) {
		return "FlowFile"
	}
	firstLetter := strings.ToUpper(s[:1])
	return firstLetter + s[1:]
}
