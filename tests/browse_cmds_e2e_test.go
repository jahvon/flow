package tests_test

import (
	stdCtx "context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/tests/utils"
)

var _ = Describe("browse e2e", Ordered, func() {
	var (
		ctx *context.Context
		run *utils.CommandRunner
	)

	BeforeAll(func() {
		ctx = utils.NewContext(stdCtx.Background(), GinkgoT())
		run = utils.NewE2ECommandRunner()
	})

	BeforeEach(func() {
		utils.ResetTestContext(ctx, GinkgoT())
	})

	AfterEach(func() {
		ctx.Finalize()
	})

	DescribeTable("browse list with various filters produces YAML output",
		func(args []string) {
			stdOut := ctx.StdOut()
			cmdArgs := append([]string{"browse", "--list"}, args...)
			Expect(run.Run(ctx, cmdArgs...)).To(Succeed())
			out, err := readFileContent(stdOut)
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(ContainSubstring("executables:"))
		},
		Entry("no filter", []string{}),
		Entry("workspace filter", []string{"--workspace", "."}),
		Entry("namespace filter", []string{"--namespace", "."}),
		Entry("all namespaces", []string{"--all"}),
		Entry("verb filter", []string{"--verb", "exec"}),
		Entry("tag filter", []string{"--tag", "test"}),
		Entry("substring filter", []string{"--filter", "print"}),
		Entry("multiple filters", []string{"--verb", "exec", "--workspace", ".", "--namespace", "."}),
	)

	It("should show executable details by verb and name", func() {
		stdOut := ctx.StdOut()
		Expect(run.Run(ctx, "browse", "exec", "examples:simple-print")).To(Succeed())
		out, err := readFileContent(stdOut)
		Expect(err).NotTo(HaveOccurred())
		Expect(out).To(ContainSubstring("name: simple-print"))
	})
})
