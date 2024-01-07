package main

import (
	"path/filepath"
	"runtime"

	"github.com/jahvon/flow/cmd"
)

const DocsDir = "docs"

var (
	_, b, _, _ = runtime.Caller(0)
	Root       = filepath.Join(filepath.Dir(b), "../..")
)

func main() {
	if err := cmd.GenerateMarkdownTree(filepath.Join(Root, DocsDir, "cli")); err != nil {
		panic(err)
	}
}
