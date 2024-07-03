package cache_test

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/jahvon/tuikit/io/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/internal/cache"
	cacheMocks "github.com/jahvon/flow/internal/cache/mocks"
	"github.com/jahvon/flow/internal/filesystem"
)

var _ = Describe("ExecutableCacheImpl", func() {
	var (
		logger              *mocks.MockLogger
		execCache           *cache.ExecutableCacheImpl
		wsCache             *cacheMocks.MockWorkspaceCache
		wsName, wsPath      string
		cacheDir, configDir string
	)

	BeforeEach(func() {
		var err error
		cacheDir, err = os.MkdirTemp("", "flow-cache-test")
		Expect(err).NotTo(HaveOccurred())
		Expect(os.Setenv(filesystem.FlowCacheDirEnvVar, cacheDir)).To(Succeed())
		configDir, err = os.MkdirTemp("", "flow-config-test")
		Expect(err).NotTo(HaveOccurred())
		Expect(os.Setenv(filesystem.FlowConfigDirEnvVar, configDir)).To(Succeed())

		wsName = "test"
		wsPath = filepath.Join(cacheDir, "workspace")
		err = filesystem.InitWorkspaceConfig(wsName, wsPath)
		Expect(err).NotTo(HaveOccurred())
		wsConfig, err := filesystem.LoadWorkspaceConfig(wsName, wsPath)
		Expect(err).NotTo(HaveOccurred())

		logger = mocks.NewMockLogger(gomock.NewController(GinkgoT()))
		execDefinition := &config.ExecutableDefinition{
			Namespace:  "testdata",
			Visibility: config.VisibilityPrivate,
			Executables: config.ExecutableList{
				{Verb: "run", Name: "exec"},
			},
		}
		execDefinition.SetContext(wsName, wsPath, filepath.Join(wsPath, "test"+filesystem.ExecutableDefinitionExt))
		err = filesystem.WriteExecutableDefinition(execDefinition.DefinitionPath(), execDefinition)
		Expect(err).NotTo(HaveOccurred())
		execCacheData := &cache.ExecutableCacheData{
			ExecutableMap: make(map[config.Ref]string),
			AliasMap:      make(map[config.Ref]config.Ref),
			DefinitionMap: make(map[string]cache.WorkspaceInfo),
		}
		wsCache = cacheMocks.NewMockWorkspaceCache(gomock.NewController(GinkgoT()))
		wsCache.EXPECT().GetLatestData(gomock.Any()).Return(&cache.WorkspaceCacheData{
			Workspaces:         map[string]*config.WorkspaceConfig{wsName: wsConfig},
			WorkspaceLocations: map[string]string{wsName: wsPath},
		}, nil).AnyTimes()
		execCache = &cache.ExecutableCacheImpl{
			Data:           execCacheData,
			WorkspaceCache: wsCache,
		}
	})

	AfterEach(func() {
		Expect(os.RemoveAll(cacheDir)).To(Succeed())
		Expect(os.Unsetenv(filesystem.FlowCacheDirEnvVar)).To(Succeed())
	})

	Describe("Update and GetExecutableList", func() {
		It("should update the executable cache from filesystem and retrieve the expected data", func() {
			logger.EXPECT().Debugf(gomock.Any()).Times(1)
			logger.EXPECT().Debugx(gomock.Any(), "workspace", wsName).Times(1)
			logger.EXPECT().Debugx(gomock.Any(), "count", 1).Times(1)
			err := execCache.Update(logger)
			Expect(err).ToNot(HaveOccurred())

			var readData config.ExecutableList
			readData, err = execCache.GetExecutableList(logger)
			Expect(err).ToNot(HaveOccurred())
			Expect(readData).ToNot(BeNil())
			execs := readData.FilterByWorkspace(wsName)
			Expect(execs).To(HaveLen(1))
		})

		When("generated executables are expected", func() {
			It("should generate the expected executables", func() {
				_, f, _, _ := runtime.Caller(0)
				err := filesystem.CopyFile(
					filepath.Join(filepath.Dir(f), "testdata", "from-file.sh"),
					filepath.Join(wsPath, "from-file.sh"),
				)
				Expect(err).NotTo(HaveOccurred())
				execDefinition := &config.ExecutableDefinition{
					Namespace:  "testdata",
					Visibility: config.VisibilityPrivate,
					FromFiles:  []string{"from-file.sh"},
				}
				execDefinition.SetContext(wsName, wsPath, filepath.Join(wsPath, "test"+filesystem.ExecutableDefinitionExt))
				err = filesystem.WriteExecutableDefinition(execDefinition.DefinitionPath(), execDefinition)
				Expect(err).NotTo(HaveOccurred())

				logger.EXPECT().Debugf(gomock.Any()).Times(1)
				logger.EXPECT().Debugx(gomock.Any(), "workspace", wsName).Times(1)
				logger.EXPECT().Debugx(gomock.Any(), "count", 1).Times(1)
				err = execCache.Update(logger)
				Expect(err).ToNot(HaveOccurred())

				var readData config.ExecutableList
				readData, err = execCache.GetExecutableList(logger)
				Expect(err).ToNot(HaveOccurred())
				Expect(readData).ToNot(BeNil())
				execs := readData.FilterByWorkspace(wsName)
				Expect(execs).To(HaveLen(1))
			})
		})
	})

	Describe("Update and GetExecutableList", func() {
		It("should update the executable cache from filesystem and retrieve the expected data", func() {
			logger.EXPECT().Debugf(gomock.Any()).Times(1)
			logger.EXPECT().Debugx(gomock.Any(), "workspace", wsName).Times(1)
			logger.EXPECT().Debugx(gomock.Any(), "count", 1).Times(1)
			err := execCache.Update(logger)
			Expect(err).ToNot(HaveOccurred())

			var readData *config.Executable
			ref := config.Ref("run test/testdata:exec")
			readData, err = execCache.GetExecutableByRef(logger, ref)
			Expect(err).ToNot(HaveOccurred())
			Expect(readData).ToNot(BeNil())
			Expect(readData.Ref()).To(Equal(ref))
		})
	})
})
