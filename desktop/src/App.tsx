import {
  ActionIcon,
  AppShell,
  Loader,
  Notification as MantineNotification,
  Text,
} from "@mantine/core";
import "@mantine/core/styles.css";
import { IconRefresh } from "@tabler/icons-react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { useEffect, useState } from "react";
import "./App.css";
import styles from "./components/AppShell.module.css";
import { Header } from "./components/Header/Header";
import { Sidebar } from "./components/Sidebar/Sidebar";
import { View, Viewer } from "./components/Viewer/Viewer";
import { useExecutable, useBackendData } from "./hooks/useBackendData";
import { SettingsProvider } from "./hooks/useSettings";
import { ThemeProvider } from "./theme/ThemeProvider";
import {
  colorFromType,
  Notification,
  NotificationType,
} from "./types/notification";

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: 1,
    },
  },
});

function AppContent() {
  const [currentView, setCurrentView] = useState<View>(View.Welcome);
  const [welcomeMessage, setWelcomeMessage] = useState<string>("");
  const [selectedExecutable, setSelectedExecutable] = useState<string | null>(
    null
  );
  const [selectedWorkspace, setSelectedWorkspace] = useState<string | null>(
    null
  );
  const [notification, setNotification] = useState<Notification | null>(null);

  const { config, workspaces, executables, isLoading, hasError, refreshAll } =
    useBackendData(selectedWorkspace);

  const { executable, executableError, isExecutableLoading } = useExecutable(
    selectedExecutable || ""
  );

  // Set initial workspace from config when it loads
  useEffect(() => {
    if (config?.currentWorkspace && workspaces && workspaces.length > 0) {
      // Only update if we don't have a selected workspace or if the config workspace is different
      if (!selectedWorkspace || config.currentWorkspace !== selectedWorkspace) {
        setSelectedWorkspace(config.currentWorkspace);
      }
    }
  }, [config, workspaces]);

  useEffect(() => {
    if (notification?.autoClose) {
      setTimeout(() => {
        setNotification(null);
      }, notification.autoCloseDelay || 6000);
    }
  }, [notification]);

  useEffect(() => {
    if (hasError) {
      setNotification({
        title: "Unexpected error",
        message: hasError.message || "An error occurred",
        type: NotificationType.Error,
        autoClose: true,
        autoCloseDelay: 6000,
      });
    }
  }, [hasError]);

  useEffect(() => {
    if (welcomeMessage === "" && executables?.length > 0) {
      setWelcomeMessage("Select an executable to get started.");
    }
  }, [executables, welcomeMessage]);

  const handleLogoClick = () => {
    setCurrentView(View.Welcome);
  };

  return (
    <AppShell
      header={{ height: 48 }}
      navbar={{ width: 300, breakpoint: "sm" }}
      padding="md"
      classNames={{
        root: styles.appShell,
        main: styles.main,
        header: styles.header,
        navbar: styles.navbar,
      }}
    >
      <AppShell.Header>
        <Header
          onCreateWorkspace={() => {}}
          onRefreshWorkspaces={() => {
            refreshAll();
            setNotification({
              title: "Refresh completed",
              message: "flow data has synced and refreshed successfully",
              type: NotificationType.Success,
              autoClose: true,
              autoCloseDelay: 3000,
            });
          }}
        />
      </AppShell.Header>

      <AppShell.Navbar>
        <Sidebar
          currentView={currentView}
          setCurrentView={setCurrentView}
          workspaces={workspaces || []}
          selectedWorkspace={selectedWorkspace}
          onSelectWorkspace={(workspaceName) => {
            setSelectedWorkspace(workspaceName);
            setCurrentView(View.Workspace);
          }}
          visibleExecutables={executables}
          onSelectExecutable={(executable) => {
            if (executable === selectedExecutable) {
              return;
            }
            setSelectedExecutable(executable);
            if (currentView !== View.Executable) {
              setCurrentView(View.Executable);
            }
          }}
          onLogoClick={handleLogoClick}
        />
      </AppShell.Navbar>

      <AppShell.Main>
        {hasError ? (
          <div
            style={{
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              height: "100%",
              flexDirection: "column",
              gap: "1rem",
            }}
          >
            <Text c="red">Error loading data</Text>
            <ActionIcon
              color="red"
              onClick={refreshAll}
              title="Try again"
              variant="light"
            >
              <IconRefresh size={16} />
            </ActionIcon>
          </div>
        ) : (
          <div style={{ position: "relative", height: "100%" }}>
            <Viewer
              currentView={currentView}
              selectedExecutable={executable}
              isExecutableLoading={isExecutableLoading}
              executableError={executableError}
              welcomeMessage={welcomeMessage}
              workspace={
                workspaces?.find((w) => w.name === selectedWorkspace) || null
              }
            />
            {isLoading && (
              <div
                style={{
                  position: "absolute",
                  top: 16,
                  right: 16,
                  zIndex: 1000,
                }}
              >
                <Loader size="sm" />
              </div>
            )}
          </div>
        )}
      </AppShell.Main>

      {notification && (
        <MantineNotification
          title={notification.title}
          onClose={() => setNotification(null)}
          color={colorFromType(notification.type)}
          style={{
            position: "fixed",
            bottom: 20,
            right: 20,
            zIndex: 1000,
          }}
        >
          {notification.message}
        </MantineNotification>
      )}
    </AppShell>
  );
}

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <SettingsProvider>
        <ThemeProvider>
          <AppContent />
        </ThemeProvider>
      </SettingsProvider>
    </QueryClientProvider>
  );
}
export default App;
