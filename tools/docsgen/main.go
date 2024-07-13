package main

import (
	stdCtx "context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra/doc"

	"github.com/jahvon/flow/cmd"
	"github.com/jahvon/flow/internal/context"
)

const (
	DocsDir = "docs"
	cliDir  = "cli"
)

func main() {
	fmt.Println("Generating CLI docs...")
	ctx := context.NewContext(stdCtx.Background(), os.Stdin, os.Stdout)
	defer ctx.Finalize()

	rootCmd := cmd.NewRootCmd(ctx)
	cmd.RegisterSubCommands(ctx, rootCmd)
	rootCmd.DisableAutoGenTag = true
	if err := doc.GenMarkdownTree(rootCmd, filepath.Join(rootDir(), DocsDir, cliDir)); err != nil {
		panic(err)
	}

	fmt.Println("Generating markdown docs...")
	generateMarkdownDocs()

	fmt.Println("Generating schema files...")
	generateJSONSchemas()
}

func rootDir() string {
	_, filename, _, _ := runtime.Caller(0)
	// ./tools/docsgen/schema.go -> ./
	return filepath.Dir(filepath.Dir(filepath.Dir(filepath.Base(filename))))
}
