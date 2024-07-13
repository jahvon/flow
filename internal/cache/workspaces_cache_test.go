package cache_test

import (
	"os"

	"github.com/jahvon/tuikit/io/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	"github.com/jahvon/flow/internal/cache"
	"github.com/jahvon/flow/internal/filesystem"
	"github.com/jahvon/flow/types/workspace"
)

var _ = Describe("WorkspaceCacheImpl", func() {
	var (
		logger    *mocks.MockLogger
		wsCache   *cache.WorkspaceCacheImpl
		cacheDir  string
		configDir string
	)

	BeforeEach(func() {
		ctrl := gomock.NewController(GinkgoT())
		logger = mocks.NewMockLogger(ctrl)
		wsCache = &cache.WorkspaceCacheImpl{
			Data: &cache.WorkspaceCacheData{
				Workspaces:         make(map[string]*workspace.Workspace),
				WorkspaceLocations: make(map[string]string),
			},
		}

		var err error
		cacheDir, err = os.MkdirTemp("", "flow-cache-test")
		Expect(err).NotTo(HaveOccurred())
		Expect(os.Setenv(filesystem.FlowCacheDirEnvVar, cacheDir)).To(Succeed())
		configDir, err = os.MkdirTemp("", "flow-config-test")
		Expect(err).NotTo(HaveOccurred())
		Expect(os.Setenv(filesystem.FlowConfigDirEnvVar, configDir)).To(Succeed())

		testWs := &workspace.Workspace{}
		testWs.SetContext("test", cacheDir)
		wsCache.Data.Workspaces["test"] = testWs
		wsCache.Data.WorkspaceLocations["test"] = cacheDir
	})

	AfterEach(func() {
		Expect(os.RemoveAll(cacheDir)).To(Succeed())
		Expect(os.Unsetenv(filesystem.FlowCacheDirEnvVar)).To(Succeed())
		Expect(os.RemoveAll(configDir)).To(Succeed())
		Expect(os.Unsetenv(filesystem.FlowConfigDirEnvVar)).To(Succeed())
	})

	Describe("Update and GetLatestData", func() {
		It("should update the workspace cache and retrieve the same data", func() {
			logger.EXPECT().Debugf(gomock.Any()).Times(1)
			logger.EXPECT().Debugx(gomock.Any(), "count", 2).Times(1)
			err := wsCache.Update(logger)
			Expect(err).ToNot(HaveOccurred())

			var readData *cache.WorkspaceCacheData
			readData, err = wsCache.GetLatestData(logger)
			Expect(err).ToNot(HaveOccurred())
			Expect(readData).To(Equal(wsCache.Data))
		})
	})

	Describe("GetWorkspaceConfigList", func() {
		It("returns the expected workspace configs", func() {
			list, err := wsCache.GetWorkspaceConfigList(logger)
			Expect(err).ToNot(HaveOccurred())
			Expect(list).To(HaveLen(1))
			Expect(list.FindByName("test")).To(Equal(wsCache.Data.Workspaces["test"]))
		})
	})
})
