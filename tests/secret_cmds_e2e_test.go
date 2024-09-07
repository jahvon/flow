package tests_test

import (
	stdCtx "context"
	"os"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/vault"
	"github.com/jahvon/flow/tests/utils"
)

var _ = FDescribe("vault/secrets e2e", Ordered, func() {
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

	When("creating a new vault (flow init vault)", func() {
		It("should return the generated key", func() {
			stdOut := ctx.StdOut()
			Expect(run.Run(ctx, "init", "vault", "--verbosity", "-1")).To(Succeed())
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
			Expect(run.Run(ctx, "set", "secret", "my-secret", "my-value")).To(Succeed())
			out, err := readFileContent(ctx.StdOut())
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(ContainSubstring("Secret my-secret set in vault"))
		})
	})

	When("getting a secret (flow get secret)", func() {
		It("should return the secret value", func() {
			stdOut := ctx.StdOut()
			Expect(run.Run(ctx, "get", "secret", "my-secret", "--plainText")).To(Succeed())
			out, err := readFileContent(stdOut)
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(ContainSubstring("my-value"))
		})
	})

	When("listing secrets (flow list secrets)", func() {
		It("should return the list of secrets", func() {
			stdOut := ctx.StdOut()
			Expect(run.Run(ctx, "list", "secrets")).To(Succeed())
			out, err := readFileContent(stdOut)
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(ContainSubstring("my-secret"))
		})
	})

	When("deleting a secret (flow remove secret)", func() {
		It("should remove the secret from the vault", func() {
			go func() {
				defer GinkgoRecover()
				Expect(writeUserInput(ctx.StdIn(), "yes")).To(Succeed())
				Expect(rewindFile(ctx.StdIn())).To(Succeed())
			}()
			Eventually(run.Run(ctx, "remove", "secret", "my-secret")).Within(time.Second * 3).Should(Succeed())
			out, err := readFileContent(ctx.StdOut())
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(ContainSubstring("Secret my-secret removed from vault"))
		})
	})
})
