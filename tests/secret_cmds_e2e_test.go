package tests_test

import (
	stdCtx "context"
	"os"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/vault"
	"github.com/jahvon/flow/tests/utils"
)

var _ = Describe("vault/secrets e2e", Ordered, func() {
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

	When("creating a new vault (flow secret vault create)", func() {
		It("should return the generated key", func() {
			stdOut := ctx.StdOut()
			Expect(run.Run(ctx, "secret", "vault", "create", "--verbosity", "-1")).To(Succeed())
			out, err := readFileContent(stdOut)
			Expect(err).NotTo(HaveOccurred())

			parts := strings.Split(strings.TrimSpace(out), ":")
			Expect(parts).To(HaveLen(2))
			encryptionKey := strings.TrimSpace(parts[1])
			Expect(os.Setenv(vault.EncryptionKeyEnvVar, encryptionKey)).To(Succeed())
		})
	})

	When("setting a secret (flow set secret)", func() {
		It("should save into the vault", func() {
			Expect(run.Run(ctx, "secret", "set", "my-secret", "my-value")).To(Succeed())
			out, err := readFileContent(ctx.StdOut())
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(ContainSubstring("Secret my-secret set in vault"))
		})
	})

	When("getting a secret (flow secret view)", func() {
		It("should return the secret value", func() {
			stdOut := ctx.StdOut()
			Expect(run.Run(ctx, "secret", "view", "my-secret", "--plainText")).To(Succeed())
			out, err := readFileContent(stdOut)
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(ContainSubstring("my-value"))
		})
	})

	When("listing secrets (flow secret list)", func() {
		It("should return the list of secrets", func() {
			stdOut := ctx.StdOut()
			Expect(run.Run(ctx, "secret", "list")).To(Succeed())
			out, err := readFileContent(stdOut)
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(ContainSubstring("my-secret"))
		})
	})

	When("deleting a secret (flow secret delete)", func() {
		It("should remove the secret from the vault", func() {
			reader, writer, err := os.Pipe()
			Expect(err).NotTo(HaveOccurred())
			_, err = writer.Write([]byte("yes\n"))
			Expect(err).ToNot(HaveOccurred())

			ctx.SetIO(reader, ctx.StdOut())
			Expect(run.Run(ctx, "secret", "delete", "my-secret")).To(Succeed())
			out, err := readFileContent(ctx.StdOut())
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(ContainSubstring("Secret 'my-secret' deleted from vault"))
		})
	})
})
