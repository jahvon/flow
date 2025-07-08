//go:build e2e

package tests_test

import (
	stdCtx "context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/flowexec/flow/tests/utils"
)

var _ = Describe("exec e2e", func() {
	var (
		ctx *utils.Context
	)

	BeforeEach(func() {
		ctx = utils.NewContext(stdCtx.Background(), GinkgoT())
	})

	AfterEach(func() {
		ctx.Finalize()
	})

	DescribeTable("with dir example executables", func(ref string) {
		runner := utils.NewE2ECommandRunner()
		stdOut := ctx.StdOut()
		Expect(runner.Run(ctx.Context, "exec", ref, "--log-level", "debug")).To(Succeed())
		Expect(readFileContent(stdOut)).To(ContainSubstring("flow completed"))
	},
		Entry("print example", "examples:simple-print"),
		Entry("tmp dir example", "examples:with-tmp-dir"),
		Entry("nameless example", ""),
		Entry("request with transformation", "examples:request-with-transform"),
	)

	When("param overrides are provided", func() {
		It("should run the executable with the provided overrides", func() {
			runner := utils.NewE2ECommandRunner()
			stdOut := ctx.StdOut()
			Expect(runner.Run(
				ctx.Context, "exec", "examples:with-params",
				"--param", "PARAM1=value1", "--param", "PARAM2=value2", "--param", "PARAM3=value3",
			)).To(Succeed())
			out, _ := readFileContent(stdOut)
			Expect(out).To(ContainSubstring("value1"))
			Expect(out).To(ContainSubstring("value2"))
			Expect(out).To(ContainSubstring("value3"))
		})
	})
})
