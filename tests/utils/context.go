package utils

import (
	stdCtx "context"
	"os"
	"path/filepath"

	tuikitIO "github.com/jahvon/tuikit/io"
	tuikitIOMocks "github.com/jahvon/tuikit/io/mocks"
	"github.com/onsi/ginkgo/v2"
	"go.uber.org/mock/gomock"
	"gopkg.in/yaml.v3"

	"github.com/jahvon/flow/internal/cache"
	cacheMocks "github.com/jahvon/flow/internal/cache/mocks"
	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/filesystem"
	"github.com/jahvon/flow/internal/io"
	"github.com/jahvon/flow/internal/runner/mocks"
	"github.com/jahvon/flow/tools/builder"
	"github.com/jahvon/flow/types/config"
	"github.com/jahvon/flow/types/workspace"
)

const (
	TestWorkspaceName        = "default"
	TestWorkspaceDisplayName = "Default Workspace"

	userConfigSubdir = "config"
	cacheSubdir      = "cache"
)

// NewContext creates a new context for testing runners. It initializes the context with
// a real logger that writes it's output to a temporary file.
// It also creates a temporary testing directory for the test workspace, user configs, and caches.
// Test environment variables are set the config and cache directories override paths.
func NewContext(ctx stdCtx.Context, t ginkgo.FullGinkgoTInterface) *context.Context {
	stdOut, stdIn := createTempIOFiles(t)
	logger := tuikitIO.NewLogger(stdOut, io.Theme(), tuikitIO.Text, "")
	ctxx := newTestContext(ctx, t, logger, stdIn, stdOut)
	return ctxx
}

// NewTestContextWithMockLogger creates a new context for testing runners. It initializes the context with
// a mock logger.
// It also creates a temporary testing directory for the test workspace, user configs, and caches.
// Test environment variables are set the config and cache directories override paths.
func NewTestContextWithMockLogger(
	ctx stdCtx.Context,
	t ginkgo.FullGinkgoTInterface,
	ctrl *gomock.Controller,
) (*context.Context, *tuikitIOMocks.MockLogger) {
	stdOut, stdIn := createTempIOFiles(t)
	logger := tuikitIOMocks.NewMockLogger(ctrl)
	expectInternalMockLoggerCalls(logger)
	ctxx := newTestContext(ctx, t, logger, stdIn, stdOut)
	return ctxx, logger
}

func ResetTestContext(ctx *context.Context, t ginkgo.FullGinkgoTInterface) {
	ctx.Ctx = stdCtx.Background()
	stdIn, stdOut := createTempIOFiles(t)
	ctx.SetIO(stdIn, stdOut)
	logger := tuikitIO.NewLogger(stdOut, io.Theme(), tuikitIO.Text, "")
	ctx.Logger = logger
}

func createTempIOFiles(t ginkgo.FullGinkgoTInterface) (stdIn *os.File, stdOut *os.File) {
	var err error
	stdOut, err = os.CreateTemp("", "flow-test-out")
	if err != nil {
		t.Fatalf("unable to create temp file: %v", err)
	}
	stdIn, err = os.CreateTemp("", "flow-test-in")
	if err != nil {
		t.Fatalf("unable to create temp file: %v", err)
	}
	return
}

func newTestContext(
	ctx stdCtx.Context,
	t ginkgo.FullGinkgoTInterface,
	logger tuikitIO.Logger,
	stdIn, stdOut *os.File,
) *context.Context {
	configDir, cacheDir, wsDir := initTestDirectories(t)
	setTestEnv(t, configDir, cacheDir)

	testWsCfg, err := testWsConfig(wsDir)
	if err != nil {
		t.Fatalf("unable to create workspace config: %v", err)
	}
	testCfg, err := testConfig(wsDir)
	if err != nil {
		t.Fatalf("unable to create user config: %v", err)
	}

	wsCache, execCache := testCaches(t, logger)

	cancel := func() {
		<-ctx.Done()
	}

	ctxx := &context.Context{
		Ctx:              ctx,
		CancelFunc:       cancel,
		Logger:           logger,
		Config:           testCfg,
		CurrentWorkspace: testWsCfg,
		WorkspacesCache:  wsCache,
		ExecutableCache:  execCache,
	}
	ctxx.SetIO(stdIn, stdOut)
	return ctxx
}

type ContextWithMocks struct {
	Ctx             *context.Context
	Logger          *tuikitIOMocks.MockLogger
	ExecutableCache *cacheMocks.MockExecutableCache
	WorkspaceCache  *cacheMocks.MockWorkspaceCache
	RunnerMock      *mocks.MockRunner
}

// NewContextWithMocks creates a new context for testing runners. It initializes the context with
// a mock logger and mock caches. The mock logger is set to expect debug calls.
func NewContextWithMocks(ctx stdCtx.Context, t ginkgo.FullGinkgoTInterface) *ContextWithMocks {
	null := os.NewFile(0, os.DevNull)
	configDir, cacheDir, wsDir := initTestDirectories(t)
	setTestEnv(t, configDir, cacheDir)
	testWsCfg, err := testWsConfig(wsDir)
	if err != nil {
		t.Fatalf("unable to create workspace config: %v", err)
	}
	testUserCfg, err := testConfig(wsDir)
	if err != nil {
		t.Fatalf("unable to create config: %v", err)
	}
	cancel := func() {
		<-ctx.Done()
	}
	logger := tuikitIOMocks.NewMockLogger(gomock.NewController(t))
	expectInternalMockLoggerCalls(logger)
	wsCache := cacheMocks.NewMockWorkspaceCache(gomock.NewController(t))
	execCache := cacheMocks.NewMockExecutableCache(gomock.NewController(t))
	ctxx := &context.Context{
		Ctx:              ctx,
		CancelFunc:       cancel,
		Logger:           logger,
		Config:           testUserCfg,
		CurrentWorkspace: testWsCfg,
		WorkspacesCache:  wsCache,
		ExecutableCache:  execCache,
	}
	ctxx.SetIO(null, null)
	return &ContextWithMocks{
		Ctx:             ctxx,
		Logger:          logger,
		ExecutableCache: execCache,
		WorkspaceCache:  wsCache,
		RunnerMock:      mocks.NewMockRunner(gomock.NewController(t)),
	}
}

func initTestDirectories(t ginkgo.FullGinkgoTInterface) (string, string, string) {
	tmpDir, err := os.MkdirTemp("", "flow-test")
	if err != nil {
		t.Fatalf("unable to create temp dir: %v", err)
	}

	tmpWsDir := filepath.Join(tmpDir, TestWorkspaceName)
	if err := os.Mkdir(filepath.Join(tmpDir, TestWorkspaceName), 0750); err != nil {
		t.Fatalf("unable to create workspace directory: %v", err)
	}

	testData := builder.ExamplesExecFlowFile()
	execDef, err := yaml.Marshal(testData)
	if err != nil {
		t.Fatalf("unable to marshal test data: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpWsDir, "testdata.flow"), execDef, 0600); err != nil {
		t.Fatalf("unable to write test data: %v", err)
	}
	tmpConfigDir := filepath.Join(tmpDir, userConfigSubdir)
	tmpCacheDir := filepath.Join(tmpDir, cacheSubdir)

	return tmpConfigDir, tmpCacheDir, tmpWsDir
}

func testConfig(wsDir string) (*config.Config, error) {
	if err := filesystem.InitConfig(); err != nil {
		return nil, err
	}
	userCfg, err := filesystem.LoadConfig()
	if err != nil {
		return nil, err
	}
	userCfg.DefaultLogMode = tuikitIO.Text
	userCfg.CurrentWorkspace = TestWorkspaceName
	userCfg.Workspaces = map[string]string{
		TestWorkspaceName: wsDir,
	}
	userCfg.Interactive = &config.Interactive{Enabled: false}
	if err = filesystem.WriteConfig(userCfg); err != nil {
		return nil, err
	}

	return userCfg, nil
}

func testWsConfig(wsDir string) (*workspace.Workspace, error) {
	if err := filesystem.InitWorkspaceConfig(TestWorkspaceName, wsDir); err != nil {
		return nil, err
	}
	wsCfg, err := filesystem.LoadWorkspaceConfig(TestWorkspaceName, wsDir)
	if err != nil {
		return nil, err
	}
	wsCfg.DisplayName = TestWorkspaceDisplayName
	if err = filesystem.WriteWorkspaceConfig(wsDir, wsCfg); err != nil {
		return nil, err
	}
	return wsCfg, nil
}

// testCaches must be called after the user and workspace configs have been created.
func testCaches(t ginkgo.FullGinkgoTInterface, logger tuikitIO.Logger) (cache.WorkspaceCache, cache.ExecutableCache) {
	wsCache := cache.NewWorkspaceCache()
	execCache := cache.NewExecutableCache(wsCache)

	if err := wsCache.Update(logger); err != nil {
		t.Fatalf("unable to update cache: %v", err)
	}
	if err := execCache.Update(logger); err != nil {
		t.Fatalf("unable to update cache: %v", err)
	}
	return wsCache, execCache
}

func setTestEnv(t ginkgo.FullGinkgoTInterface, configDir, cacheDir string) {
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0750); err != nil {
			t.Fatalf("unable to create config directory: %v", err)
		}
	}
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		if err := os.MkdirAll(cacheDir, 0750); err != nil {
			t.Fatalf("unable to create cache directory: %v", err)
		}
	}

	t.Setenv(filesystem.FlowConfigDirEnvVar, configDir)
	t.Setenv(filesystem.FlowCacheDirEnvVar, cacheDir)
}

func expectInternalMockLoggerCalls(logger *tuikitIOMocks.MockLogger) {
	logger.EXPECT().Debugf(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Debugx(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().LogMode().AnyTimes()
	logger.EXPECT().SetMode(gomock.Any()).AnyTimes()
}
