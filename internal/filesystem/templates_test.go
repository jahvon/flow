package filesystem_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/filesystem"
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

	Describe("RenderAndWriteExecutablesTemplate", func() {
		It("renders and writes the template correctly", func() {
			definitionTemplate := &config.ExecutableDefinitionTemplate{
				ExecutableDefinition: &config.ExecutableDefinition{
					Namespace: "test",
					Executables: config.ExecutableList{
						{Verb: "run", Name: "test-executable", Description: "{{ .key }}"},
					},
				},
				Data: config.TemplateData{{Key: "key", Default: "value"}},
			}
			definitionTemplate.SetContext(filepath.Join(tmpDir, "test.tmpl"+filesystem.ExecutableDefinitionExt))

			workspaceConfig := config.DefaultWorkspaceConfig("test")
			workspaceConfig.SetContext("test", tmpDir)

			err := filesystem.RenderAndWriteExecutablesTemplate(definitionTemplate, workspaceConfig, "test", "")
			Expect(err).ToNot(HaveOccurred())
			_, err = os.Stat(filepath.Join(tmpDir, "test"+filesystem.ExecutableDefinitionExt))
			Expect(err).NotTo(HaveOccurred())
		})

		It("renders and writes the template with artifacts correctly", func() {
			definitionTemplate := &config.ExecutableDefinitionTemplate{
				ExecutableDefinition: &config.ExecutableDefinition{
					Namespace: "test",
					Executables: config.ExecutableList{
						{Verb: "run", Name: "test-executable", Description: "{{ .key }}"},
					},
				},
				Data:      config.TemplateData{{Key: "key", Default: "value"}},
				Artifacts: []string{"subpath/test-artifact"},
			}
			definitionTemplate.SetContext(filepath.Join(tmpDir, "test.tmpl"+filesystem.ExecutableDefinitionExt))

			err := os.MkdirAll(filepath.Join(tmpDir, "subpath"), 0750)
			Expect(err).NotTo(HaveOccurred())
			_, err = os.Create(filepath.Join(tmpDir, "subpath", "test-artifact"))
			Expect(err).NotTo(HaveOccurred())

			workspaceConfig := config.DefaultWorkspaceConfig("test")
			workspaceConfig.SetContext("test", tmpDir)

			Expect(filesystem.RenderAndWriteExecutablesTemplate(definitionTemplate, workspaceConfig, "test", "")).To(Succeed())
			_, err = os.Stat(definitionTemplate.Location())
			Expect(err).NotTo(HaveOccurred())
			_, err = os.Stat(filepath.Join(tmpDir, "subpath", "test-artifact"))
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("LoadExecutableDefinitionTemplate", func() {
		It("loads the template correctly", func() {
			definitionTemplate := &config.ExecutableDefinitionTemplate{
				ExecutableDefinition: &config.ExecutableDefinition{
					Namespace: "test",
					Executables: config.ExecutableList{
						{Verb: "exec", Name: "test-executable", Description: "{{ .key }}"},
					},
				},
				Data: config.TemplateData{{Key: "key", Default: "value"}},
			}
			definitionTemplate.SetContext(filepath.Join(tmpDir, "test.tmpl"+filesystem.ExecutableDefinitionExt))

			Expect(filesystem.WriteExecutableDefinitionTemplate(definitionTemplate.Location(), definitionTemplate)).To(Succeed())

			readTemplate, err := filesystem.LoadExecutableDefinitionTemplate(definitionTemplate.Location())
			Expect(err).NotTo(HaveOccurred())
			Expect(readTemplate).To(Equal(definitionTemplate))
		})
	})
})
