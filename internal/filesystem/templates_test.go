package filesystem_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/jahvon/flow/internal/filesystem"
	"github.com/jahvon/flow/types/executable"
	"github.com/jahvon/flow/types/workspace"
)

var _ = Describe("Templates", func() {
	var (
		tmpDir string
	)

	BeforeEach(func() {
		var err error
		tmpDir, err = os.MkdirTemp("", "flow-templates-test")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		Expect(os.RemoveAll(tmpDir)).To(Succeed())
	})

	Describe("WriteFlowFileFromTemplate", func() {
		It("renders and writes the template correctly", func() {
			definitionTemplate := &executable.FlowFileTemplate{
				FlowFile: &executable.FlowFile{
					Namespace: "test",
					Executables: executable.ExecutableList{
						{Verb: "run", Name: "test-executable", Description: "{{ .key }}"},
					},
				},
				Data: executable.TemplateData{{Key: "key", Default: "value"}},
			}
			definitionTemplate.SetContext(filepath.Join(tmpDir, "test.tmpl"+filesystem.FlowFileExt))

			workspaceConfig := workspace.DefaultWorkspaceConfig("test")
			workspaceConfig.SetContext("test", tmpDir)

			err := filesystem.WriteFlowFileFromTemplate(definitionTemplate, workspaceConfig, "test", "")
			Expect(err).ToNot(HaveOccurred())
			_, err = os.Stat(filepath.Join(tmpDir, "test"+filesystem.FlowFileExt))
			Expect(err).NotTo(HaveOccurred())
		})

		It("renders and writes the template with artifacts correctly", func() {
			definitionTemplate := &executable.FlowFileTemplate{
				FlowFile: &executable.FlowFile{
					Namespace: "test",
					Executables: executable.ExecutableList{
						{Verb: "run", Name: "test-executable", Description: "{{ .key }}"},
					},
				},
				Data:      executable.TemplateData{{Key: "key", Default: "value"}},
				Artifacts: []string{"subpath/test-artifact"},
			}
			definitionTemplate.SetContext(filepath.Join(tmpDir, "test.tmpl"+filesystem.FlowFileExt))

			err := os.MkdirAll(filepath.Join(tmpDir, "subpath"), 0750)
			Expect(err).NotTo(HaveOccurred())
			_, err = os.Create(filepath.Join(tmpDir, "subpath", "test-artifact"))
			Expect(err).NotTo(HaveOccurred())

			workspaceConfig := workspace.DefaultWorkspaceConfig("test")
			workspaceConfig.SetContext("test", tmpDir)

			Expect(filesystem.WriteFlowFileFromTemplate(definitionTemplate, workspaceConfig, "test", "")).To(Succeed())
			_, err = os.Stat(filepath.Join(tmpDir, "test"+filesystem.FlowFileExt))
			Expect(err).NotTo(HaveOccurred())
			_, err = os.Stat(filepath.Join(tmpDir, "test-artifact"))
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("LoadFlowFileTemplate", func() {
		It("loads the template correctly", func() {
			definitionTemplate := &executable.FlowFileTemplate{
				FlowFile: &executable.FlowFile{
					Namespace: "test",
					Executables: executable.ExecutableList{
						{Verb: "exec", Name: "test-executable", Description: "{{ .key }}"},
					},
				},
				Data: executable.TemplateData{{Key: "key", Default: "value"}},
			}
			definitionTemplate.SetContext(filepath.Join(tmpDir, "test.tmpl"+filesystem.FlowFileExt))

			Expect(filesystem.WriteFlowFileTemplate(definitionTemplate.Location(), definitionTemplate)).To(Succeed())

			readTemplate, err := filesystem.LoadFlowFileTemplate(definitionTemplate.Location())
			Expect(err).NotTo(HaveOccurred())
			Expect(readTemplate).To(Equal(definitionTemplate))
		})
	})
})
