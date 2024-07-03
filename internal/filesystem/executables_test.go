package filesystem_test

import (
	"os"
	"path/filepath"

	"github.com/jahvon/tuikit/io/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/filesystem"
)

var _ = Describe("Executables", func() {
	var (
		tmpDir string
	)

	BeforeEach(func() {
		var err error
		tmpDir, err = os.MkdirTemp("", "flow-executables-test")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		Expect(os.RemoveAll(tmpDir)).To(Succeed())
	})

	Describe("EnsureExecutableDir", func() {
		It("creates the directory if it does not exist", func() {
			Expect(filesystem.EnsureExecutableDir(tmpDir, "subPath")).To(Succeed())
			_, err := os.Stat(filepath.Join(tmpDir, "subPath"))
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("WriteExecutableDefinition and LoadExecutableDefinition", func() {
		It("writes and reads executable definition correctly", func() {
			executableDefinition := &config.ExecutableDefinition{
				Namespace: "test",
				Executables: config.ExecutableList{
					{
						Verb: "exec",
						Name: "test-executable",
					},
				},
			}

			definitionFile := filepath.Join(tmpDir, "test"+filesystem.ExecutableDefinitionExt)
			Expect(filesystem.WriteExecutableDefinition(definitionFile, executableDefinition)).To(Succeed())

			readDefinition, err := filesystem.LoadExecutableDefinition(definitionFile)
			Expect(err).NotTo(HaveOccurred())
			Expect(readDefinition).To(Equal(executableDefinition))
		})
	})

	Describe("LoadWorkspaceExecutableDefinitions", func() {
		It("loads all executable definitions if no paths are set", func() {
			executableDefinition := &config.ExecutableDefinition{
				Namespace: "test",
				Executables: config.ExecutableList{
					{
						Verb: "exec",
						Name: "test-executable",
					},
				},
			}

			definitionFile := filepath.Join(tmpDir, "test"+filesystem.ExecutableDefinitionExt)
			Expect(filesystem.WriteExecutableDefinition(definitionFile, executableDefinition)).To(Succeed())

			workspaceCfg := &config.WorkspaceConfig{}
			workspaceCfg.SetContext("test", tmpDir)

			ctrl := gomock.NewController(GinkgoT())
			logger := mocks.NewMockLogger(ctrl)
			logger.EXPECT().Debugx(gomock.Any(), gomock.Any()).AnyTimes()

			definitions, err := filesystem.LoadWorkspaceExecutableDefinitions(logger, workspaceCfg)
			Expect(err).NotTo(HaveOccurred())
			Expect(definitions).To(HaveLen(1))
			Expect(definitions[0].Namespace).To(Equal(executableDefinition.Namespace))
		})
		It("loads executable definitions from the included path", func() {
			executableDefinition := &config.ExecutableDefinition{
				Namespace: "test",
				Executables: config.ExecutableList{
					{
						Verb: "exec",
						Name: "test-executable",
					},
				},
			}

			definitionFile := filepath.Join(tmpDir, "test"+filesystem.ExecutableDefinitionExt)
			Expect(filesystem.WriteExecutableDefinition(definitionFile, executableDefinition)).To(Succeed())

			workspaceCfg := &config.WorkspaceConfig{
				Executables: &config.ExecutableLocationConfig{
					Included: []string{tmpDir},
					Excluded: []string{filepath.Join(tmpDir, "excluded")},
				},
			}
			workspaceCfg.SetContext("test", tmpDir)

			ctrl := gomock.NewController(GinkgoT())
			logger := mocks.NewMockLogger(ctrl)
			logger.EXPECT().Debugx(gomock.Any(), gomock.Any()).AnyTimes()

			definitions, err := filesystem.LoadWorkspaceExecutableDefinitions(logger, workspaceCfg)
			Expect(err).NotTo(HaveOccurred())
			Expect(definitions).To(HaveLen(1))
			Expect(definitions[0].Namespace).To(Equal(executableDefinition.Namespace))
		})

		It("does not load executable definitions from excluded paths", func() {
			executableDefinition := &config.ExecutableDefinition{
				Namespace: "test",
				Executables: config.ExecutableList{
					{
						Verb: "exec",
						Name: "test-executable",
					},
				},
			}

			excludedDir, err := os.MkdirTemp(tmpDir, "excluded")
			Expect(err).NotTo(HaveOccurred())

			definitionFile := filepath.Join(excludedDir, "test"+filesystem.ExecutableDefinitionExt)
			Expect(filesystem.WriteExecutableDefinition(definitionFile, executableDefinition)).To(Succeed())

			workspaceCfg := &config.WorkspaceConfig{
				Executables: &config.ExecutableLocationConfig{
					Included: []string{tmpDir},
					Excluded: []string{excludedDir},
				},
			}
			workspaceCfg.SetContext("test", tmpDir)

			ctrl := gomock.NewController(GinkgoT())
			logger := mocks.NewMockLogger(ctrl)
			logger.EXPECT().Debugx(gomock.Any(), gomock.Any()).AnyTimes()

			definitions, err := filesystem.LoadWorkspaceExecutableDefinitions(logger, workspaceCfg)
			Expect(err).NotTo(HaveOccurred())
			Expect(definitions).To(BeEmpty())
		})
	})
})
