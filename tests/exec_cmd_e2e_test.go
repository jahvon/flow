package tests_test

import (
	stdCtx "context"
	"fmt"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/tests/runner"
)

var _ = Describe("Examples End-to-end", func() {
	var (
		ctx *context.Context
	)

	BeforeEach(func() {
		ctx = runner.NewTestContext(stdCtx.Background(), GinkgoT())
	})

	AfterEach(func() {
		ctx.Finalize()
	})

	DescribeTable("with dir example executables", func(ref string) {
		runner := runner.NewE2ECommandRunner(ctx)
		stdOut := ctx.StdOut()
		Expect(runner.Run(ctx, "exec", ref)).To(Succeed())
		Expect(readFileContent(stdOut)).To(ContainSubstring("flow completed"))
	},
		Entry("tmp dir example", "examples:tmp-dir"),
	)
})

func readFileContent(f *os.File) (string, error) {
	out, err := os.ReadFile(f.Name())
	if err != nil {
		return "", err
	}
	outStr := string(out)
	if os.Getenv("SUPPRESS_OUTPUT") == "" {
		fmt.Println(outStr)
	}
	return outStr, nil
}
