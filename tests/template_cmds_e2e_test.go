package tests_test

import (
	stdCtx "context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/tests/utils"
	"github.com/jahvon/flow/types/executable"
)

var _ = Describe("flowfile template commands e2e", Ordered, func() {
	var (
		ctx              *context.Context
		run              *utils.CommandRunner
		template         *executable.Template
		expectedFlowFile *executable.FlowFile
	)

	BeforeAll(func() {
		ctx = utils.NewContext(stdCtx.Background(), GinkgoT())
		run = utils.NewE2ECommandRunner()
		workDir, err := os.MkdirTemp("", "flowfile-template-e2e")
		Expect(err).NotTo(HaveOccurred())
		tmpl := executable.FlowFile{
			Namespace:   "test",
			Description: "Template test flowfile",
			Tags:        []string{"test"},
			Executables: []*executable.Executable{
				{
					Verb: "exec",
					Name: "{{ .Name }}",
					Exec: &executable.ExecExecutableType{Cmd: fmt.Sprintf("echo '%s'", "{{ .Msg }}")}},
			},
		}
		tmplStr, err := tmpl.YAML()
		Expect(err).NotTo(HaveOccurred())
		template = &executable.Template{
			Template: tmplStr,
			Form: executable.FormFields{
				&executable.Field{
					Key:     "Name",
					Prompt:  "Enter a name",
					Default: "test",
				},
				&executable.Field{
					Key:     "Msg",
					Prompt:  "Enter a message",
					Default: "Hello, world!",
				},
			},
			Artifacts: []executable.Artifact{
				{SrcName: "artifact1"},
				{SrcName: "artifact2", DstName: "artifact2-renamed"},
			},
			PostRun: []executable.TemplatePostRunElem{
				{
					Cmd: "touch {{ .Name }}",
				},
			},
		}
		template.SetContext("e2e", filepath.Join(workDir, "flowfile.tmpl.flow"))
		data, err := template.YAML()
		Expect(err).NotTo(HaveOccurred())
		Expect(os.WriteFile(filepath.Join(workDir, "flowfile.tmpl.flow"), []byte(data), 0644)).To(Succeed())

		expectedFlowFile = &executable.FlowFile{
			Namespace:   "test",
			Description: "Template test flowfile",
			Tags:        []string{"test"},
			Executables: []*executable.Executable{
				{
					Verb: "exec",
					Name: "test",
					Exec: &executable.ExecExecutableType{Cmd: "echo 'Hello, world!'"}},
			},
		}
		expectedFlowFile.SetContext(
			ctx.CurrentWorkspace.AssignedName(),
			ctx.CurrentWorkspace.Location(),
			filepath.Join(workDir, "flowfile.flow"),
		)
	})

	BeforeEach(func() {
		utils.ResetTestContext(ctx, GinkgoT())
	})

	AfterEach(func() {
		ctx.Finalize()
	})

	When("registering a new template (flow set template)", func() {
		It("should complete successfully", func() {
			stdOut := ctx.StdOut()
			err := run.Run(ctx, "set", "template", "--verbosity", "-1", template.Name(), template.Location())
			Expect(err).ToNot(HaveOccurred())
			out, err := readFileContent(stdOut)
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(ContainSubstring(fmt.Sprintf("Template %s set", template.Name())))
		})
	})

	When("getting a registered template (flow get template)", func() {
		It("should return the template", func() {
			stdOut := ctx.StdOut()
			err := run.Run(ctx, "get", "template", "-t", template.Name(), "-o", "yaml")
			Expect(err).ToNot(HaveOccurred())
			out, err := readFileContent(stdOut)
			Expect(err).NotTo(HaveOccurred())
			str, err := template.YAML()
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(ContainSubstring(str))
		})
	})

	When("getting a template by path (flow get template)", func() {
		It("should return the template", func() {
			stdOut := ctx.StdOut()
			err := run.Run(ctx, "get", "template", "-f", template.Location(), "-o", "yaml")
			Expect(err).ToNot(HaveOccurred())
			out, err := readFileContent(stdOut)
			Expect(err).NotTo(HaveOccurred())
			str, err := template.YAML()
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(ContainSubstring(str))
		})
	})

	When("Listing all registered templates (flow list templates)", func() {
		It("should return the list of templates", func() {
			stdOut := ctx.StdOut()
			Expect(run.Run(ctx, "list", "templates", "-o", "yaml")).To(Succeed())
			out, err := readFileContent(stdOut)
			Expect(err).NotTo(HaveOccurred())
			// tabs may be present so instead of checking for exact match, we check for length
			str, err := template.YAML()
			Expect(err).NotTo(HaveOccurred())
			Expect(len(out)).To(BeNumerically(">", len(str)))
		})
	})

	When("Rendering a template (flow init template)", func() {
		It("should process the template options and render the flowfile", func() {
			stdIn := ctx.StdIn()
			stdOut := ctx.StdOut()
			name := "test"
			outputDir := filepath.Join(ctx.CurrentWorkspace.Location(), "output")
			writeInput := func() {
				defer GinkgoRecover()
				t := time.Tick(time.Millisecond)
				deadline := time.Now().Add(time.Second * 3)
				for range t {
					if ctx.TUIContainer != nil && ctx.TUIContainer.Ready() {
						break
					} else if time.Now().After(deadline) {
						Fail("timed out waiting for interactive container to be ready")
					}
				}
				Expect(writeUserInput(stdIn, "test")).To(Succeed())
				Expect(writeUserInput(stdIn, "hello")).To(Succeed())
				Expect(rewindFile(stdIn)).To(Succeed())
			}
			go writeInput()
			Eventually(
				run.Run(ctx, "init", "template", name, "-t", template.Name(), "-o", outputDir),
			).Within(time.Second * 3).Should(Succeed())
			out, err := readFileContent(stdOut)
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(ContainSubstring(fmt.Sprintf("Template '%s' rendered successfully", name)))
		})
	})
})
