package runner

import (
	stdCtx "context"
	"os"
	"path/filepath"
	"runtime"

	tuikitIO "github.com/jahvon/tuikit/io"
	tuikitIOMocks "github.com/jahvon/tuikit/io/mocks"
	"github.com/onsi/ginkgo/v2"
	"github.com/otiai10/copy"
	"go.uber.org/mock/gomock"
	"gopkg.in/yaml.v3"

	"github.com/jahvon/flow/config"
	"github.com/jahvon/flow/config/cache"
	"github.com/jahvon/flow/config/file"
	"github.com/jahvon/flow/internal/context"
	"github.com/jahvon/flow/internal/io"
	examplestest "github.com/jahvon/flow/tests/examples"
)

const (
	TestWorkspaceName        = "default"
	TestWorkspaceDisplayName = "Default Workspace"

	userConfigSubdir = "config"
	cacheSubdir      = "cache"
)

// NewTestContext creates a new context for testing runners. It initializes the context with
// a real logger that writes it's output to a temporary file.
// It also creates a temporary testing directory for the test workspace, user configs, and caches.
// Test environment variables are set the config and cache directories override paths.
func NewTestContext(
	ctx stdCtx.Context,
	t ginkgo.FullGinkgoTInterface,
) *context.Context {
	stdOut, err := os.CreateTemp("", "flow-test-out")
	if err != nil {
		t.Fatalf("unable to create temp file: %v", err)
	}
	stdIn, err := os.CreateTemp("", "flow-test-in")
	if err != nil {
		t.Fatalf("unable to create temp file: %v", err)
	}

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
	stdOut, err := os.CreateTemp("", "flow-test-out")
	if err != nil {
		t.Fatalf("unable to create temp file: %v", err)
	}
	stdIn, err := os.CreateTemp("", "flow-test-in")
	if err != nil {
		t.Fatalf("unable to create temp file: %v", err)
	}

	logger := tuikitIOMocks.NewMockLogger(ctrl)
	expectInternalMockLoggerCalls(logger)
	ctxx := newTestContext(ctx, t, logger, stdIn, stdOut)
	return ctxx, logger
}

func newTestContext(
	ctx stdCtx.Context,
	t ginkgo.FullGinkgoTInterface,
	logger tuikitIO.Logger,
	stdIn, stdOut *os.File,
) *context.Context {
	examples := examplesDir()
	configDir, cacheDir, wsDir := initTestDirectories(t, examples)
	setTestEnv(t, configDir, cacheDir)

	testWsCfg, err := testWsConfig(wsDir)
	if err != nil {
		t.Fatalf("unable to create workspace config: %v", err)
	}
	testUserCfg, err := testUserConfig(wsDir)
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
		UserConfig:       testUserCfg,
		CurrentWorkspace: testWsCfg,
		WorkspacesCache:  wsCache,
		ExecutableCache:  execCache,
	}
	ctxx.SetIO(stdIn, stdOut)
	return ctxx
}

func initTestDirectories(t ginkgo.FullGinkgoTInterface, srcWsDir string) (string, string, string) {
	tmpDir, err := os.MkdirTemp("", "flow-test")
	if err != nil {
		t.Fatalf("unable to create temp dir: %v", err)
	}

	tmpWsDir := filepath.Join(tmpDir, TestWorkspaceName)
	if err := os.Mkdir(filepath.Join(tmpDir, TestWorkspaceName), 0750); err != nil {
		t.Fatalf("unable to create workspace directory: %v", err)
	}
	if err := copy.Copy(srcWsDir, tmpWsDir); err != nil {
		t.Fatalf("unable to copy workspace directory: %v", err)
	}
	testData := examplestest.TestExecutableDefinition
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

func testUserConfig(wsDir string) (*config.UserConfig, error) {
	if err := file.InitUserConfig(); err != nil {
		return nil, err
	}
	userCfg, err := file.LoadUserConfig()
	if err != nil {
		return nil, err
	}
	userCfg.DefaultLogMode = tuikitIO.Text
	userCfg.CurrentWorkspace = TestWorkspaceName
	userCfg.Workspaces = map[string]string{
		TestWorkspaceName: wsDir,
	}
	userCfg.Interactive = &config.InteractiveConfig{Enabled: false}
	if err = file.WriteUserConfig(userCfg); err != nil {
		return nil, err
	}

	return userCfg, nil
}

func testWsConfig(wsDir string) (*config.WorkspaceConfig, error) {
	if err := file.InitWorkspaceConfig(TestWorkspaceName, wsDir); err != nil {
		return nil, err
	}
	wsCfg, err := file.LoadWorkspaceConfig(TestWorkspaceName, wsDir)
	if err != nil {
		return nil, err
	}
	wsCfg.DisplayName = TestWorkspaceDisplayName
	if err = file.WriteWorkspaceConfig(wsDir, wsCfg); err != nil {
		return nil, err
	}
	return wsCfg, nil
}

// testCaches must be called after the user and workspace configs have been created.
func testCaches(t ginkgo.FullGinkgoTInterface, logger tuikitIO.Logger) (*cache.WorkspaceCache, *cache.ExecutableCache) {
	wsCache := cache.NewWorkspaceCache()
	execCache := cache.NewExecutableCache()

	if err := wsCache.Update(logger); err != nil {
		t.Fatalf("unable to update cache: %v", err)
	}
	if err := execCache.Update(logger); err != nil {
		t.Fatalf("unable to update cache: %v", err)
	}
	return wsCache, execCache
}

func examplesDir() string {
	_, curFile, _, _ := runtime.Caller(0)
	// tests/runner/context.go -> tests/runner -> tests -> examples
	return filepath.Join(filepath.Dir(filepath.Dir(filepath.Dir(curFile))), "examples")
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

	t.Setenv(file.FlowConfigDirEnvVar, configDir)
	t.Setenv(file.FlowCacheDirEnvVar, cacheDir)
}

func expectInternalMockLoggerCalls(logger *tuikitIOMocks.MockLogger) {
	logger.EXPECT().Debugf(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Debugx(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().LogMode().AnyTimes()
	logger.EXPECT().SetMode(gomock.Any()).AnyTimes()
}
