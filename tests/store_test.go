package tests_test

import (
	stdCtx "context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/tests/utils"
)

var _ = Describe("store e2e", Ordered, func() {
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

	When("setting a key (flow store set)", func() {
		It("should save the value into the store", func() {
			Expect(run.Run(ctx, "store", "set", "my-key", "my-value")).To(Succeed())
			out, err := readFileContent(ctx.StdOut())
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(ContainSubstring("Key \"my-key\" set in the store"))
		})
	})

	When("getting a value (flow store get)", func() {
		It("should return the secret value", func() {
			stdOut := ctx.StdOut()
			Expect(run.Run(ctx, "store", "get", "my-key")).To(Succeed())
			out, err := readFileContent(stdOut)
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(ContainSubstring("my-value"))
		})
	})

	When("clearing the store (flow store clear)", func() {
		It("should remove all set secrets", func() {
			stdOut := ctx.StdOut()
			Expect(run.Run(ctx, "store", "clear")).To(Succeed())
			out, err := readFileContent(stdOut)
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(ContainSubstring("Data store cleared"))
		})
	})
})
