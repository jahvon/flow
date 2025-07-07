//go:build e2e

package tests_test

import (
	stdCtx "context"
	"os"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	internalUtils "github.com/jahvon/flow/internal/utils"
	"github.com/jahvon/flow/tests/utils"
)

var _ = Describe("vault/secrets e2e", Ordered, func() {
	var (
		ctx *utils.Context
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

	When("creating a new vault (flow vault create)", func() {
		It("should return the generated key", func() {
			stdOut := ctx.StdOut()
			keyEnv := "FLOW_TEST_VAULT_KEY"
			Expect(run.Run(ctx.Context, "vault", "create", "test", "--key-env", keyEnv)).To(Succeed())
			out, err := readFileContent(stdOut)
			Expect(err).NotTo(HaveOccurred())

			lines := strings.Split(strings.TrimSpace(out), "\n")
			Expect(lines).ToNot(BeEmpty())
			parts := strings.Split(strings.TrimSpace(lines[0]), ":")
			Expect(parts).To(HaveLen(2))
			encryptionKey := strings.TrimSpace(internalUtils.StripAnsi(parts[1]))
			Expect(os.Setenv(keyEnv, encryptionKey)).To(Succeed())
		})

		It("should create vault with custom path", func() {
			stdOut := ctx.StdOut()
			tmpdir, err := os.MkdirTemp("", "flow-vault-test")
			Expect(err).NotTo(HaveOccurred())
			defer os.RemoveAll(tmpdir)

			Expect(run.Run(ctx.Context, "vault", "create", "test2", "--type", "aes256", "--path", tmpdir)).To(Succeed())
			out, err := readFileContent(stdOut)
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(ContainSubstring("Vault 'test2' with AES256 encryption created successfully"))
		})
	})

	It("Should remove the created vault", func() {
		reader, writer, err := os.Pipe()
		Expect(err).NotTo(HaveOccurred())
		_, err = writer.Write([]byte("yes\n"))
		Expect(err).ToNot(HaveOccurred())

		ctx.SetIO(reader, ctx.StdOut())
		Expect(run.Run(ctx.Context, "vault", "remove", "test2")).To(Succeed())
		out, err := readFileContent(ctx.StdOut())
		Expect(err).NotTo(HaveOccurred())
		Expect(out).To(ContainSubstring("Vault 'test2' deleted"))
	})

	When("switching vaults (flow vault switch)", func() {
		It("should switch to demo vault successfully", func() {
			Expect(run.Run(ctx.Context, "vault", "switch", "demo")).To(Succeed())
			out, err := readFileContent(ctx.StdOut())
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(ContainSubstring("Vault set to demo"))
		})
	})

	When("getting vault information (flow vault get)", func() {
		It("should get demo vault in YAML format", func() {
			stdOut := ctx.StdOut()
			Expect(run.Run(ctx.Context, "vault", "get", "demo")).To(Succeed())
			out, err := readFileContent(stdOut)
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(ContainSubstring("name: demo"))
			Expect(out).To(ContainSubstring("type: demo"))
		})
	})

	When("listing vaults (flow vault list)", func() {
		It("should list vaults in YAML format", func() {
			stdOut := ctx.StdOut()
			Expect(run.Run(ctx.Context, "vault", "list")).To(Succeed())
			out, err := readFileContent(stdOut)
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(ContainSubstring("vaults:"))
		})
	})
	When("setting a secret (flow secret set)", func() {
		It("should save into the vault", func() {
			Expect(run.Run(ctx.Context, "secret", "set", "message", "my-value")).To(Succeed())
			out, err := readFileContent(ctx.StdOut())
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(ContainSubstring("Secret message set in vault"))
		})
	})

	When("getting a secret (flow secret get)", func() {
		It("should return the secret value", func() {
			stdOut := ctx.StdOut()
			Expect(run.Run(ctx.Context, "secret", "get", "message", "--plaintext")).To(Succeed())
			out, err := readFileContent(stdOut)
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(ContainSubstring("Thanks for trying flow!"))
		})
	})

	When("listing secrets (flow secret list)", func() {
		It("should return the list of secrets", func() {
			stdOut := ctx.StdOut()
			Expect(run.Run(ctx.Context, "secret", "list")).To(Succeed())
			out, err := readFileContent(stdOut)
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(ContainSubstring("message"))
		})
	})

	When("deleting a secret (flow secret remove)", func() {
		It("should remove the secret from the vault", func() {
			reader, writer, err := os.Pipe()
			Expect(err).NotTo(HaveOccurred())
			_, err = writer.Write([]byte("yes\n"))
			Expect(err).ToNot(HaveOccurred())

			ctx.SetIO(reader, ctx.StdOut())
			Expect(run.Run(ctx.Context, "secret", "remove", "message")).To(Succeed())
			out, err := readFileContent(ctx.StdOut())
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(ContainSubstring("Secret 'message' deleted from vault"))
		})
	})
})
