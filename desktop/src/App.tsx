import "@mantine/core/styles.css";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { useEffect, useState } from "react";
import "./App.css";
import { useBackendData, useExecutable } from "./hooks/useBackendData";
import { NotifierProvider, useNotifier } from "./hooks/useNotifier";
import { SettingsProvider } from "./hooks/useSettings";
import { AppShell, View, Viewer } from "./layout";
import { ThemeProvider } from "./theme/ThemeProvider";
import { NotificationType } from "./types/notification";

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
  const { notification, setNotification } = useNotifier();

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
      currentView={currentView}
      setCurrentView={(view: string) => setCurrentView(view as View)}
      workspaces={workspaces || []}
      selectedWorkspace={selectedWorkspace}
      onSelectWorkspace={(workspaceName) => {
        setSelectedWorkspace(workspaceName);
        setCurrentView(View.Workspace);
      }}
      visibleExecutables={executables || []}
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
      hasError={hasError}
      isLoading={isLoading}
      refreshAll={refreshAll}
      notification={notification}
      setNotification={setNotification}
    >
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
    </AppShell>
  );
}

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <NotifierProvider>
        <SettingsProvider>
          <ThemeProvider>
            <AppContent />
          </ThemeProvider>
        </SettingsProvider>
      </NotifierProvider>
    </QueryClientProvider>
  );
}
export default App;
