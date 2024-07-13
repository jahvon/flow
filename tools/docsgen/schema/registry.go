package schema

import (
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

func RegisteredSchemaMap() map[FileName]*JSONSchema {
	schemas := make(map[FileName]*JSONSchema)
	for _, file := range SchemaFilesForDocs {
		s := &JSONSchema{}
		p := filepath.Clean(filepath.Join(rootDir(), TypesRootDir, file.String()))
		f, err := os.ReadFile(p)
		if err != nil {
			panic(err)
		}
		if err := yaml.Unmarshal(f, s); err != nil {
			panic(err)
		}
		schemas[file] = s
	}
	return schemas
}

func rootDir() string {
	_, filename, _, _ := runtime.Caller(0)
	// ./tools/docsgen/schema.go -> ./
	return filepath.Dir(filepath.Dir(filepath.Dir(filepath.Base(filename))))
}
