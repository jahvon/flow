import { ActionIcon, AppShell, Group, Loader, Text } from "@mantine/core";
import "@mantine/core/styles.css";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { useEffect, useState } from "react";
import "./App.css";
import { Sidebar } from "./components/Sidebar/Sidebar";
import { View } from "./components/Viewer/Viewer";
import { WorkspaceInfoModal } from "./components/Workspace/WorkspaceInfoModal";
import { useWorkspaceData } from "./hooks/useWorkspaceData";

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: 1,
    },
  },
});

function AppContent() {
  const [currentView, setCurrentView] = useState<View>(View.Workspaces);
  const [selectedWorkspace, setSelectedWorkspace] = useState<string | null>(
    null
  );
  const [isWorkspaceInfoOpen, setIsWorkspaceInfoOpen] = useState(false);

  const { config, workspaces, executables, isLoading, hasError, refreshAll } =
    useWorkspaceData(selectedWorkspace);

  // Set initial workspace from config when it loads
  useEffect(() => {
    if (config?.currentWorkspace && Object.keys(workspaces || {}).length > 0) {
      setSelectedWorkspace(config.currentWorkspace);
    }
  }, [config, workspaces]);

  const handleWorkspaceInfoClick = () => {
    if (selectedWorkspace) {
      setIsWorkspaceInfoOpen(true);
    }
  };

  return (
    <AppShell
      header={{ height: 60 }}
      navbar={{ width: 360, breakpoint: "sm" }}
      padding="md"
    >
      <AppShell.Header>
        <Group h="100%" px="md" justify="space-between">
          <img
            src={`/logo-${window.matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light"}.png`}
            alt="Flow"
            height="32"
          />
          <ActionIcon
            color="palette"
            onClick={refreshAll}
            loading={isLoading}
            title="Refresh data"
            size="lg"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="16"
              height="16"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
            >
              <path d="M21 2v6h-6" />
              <path d="M3 12a9 9 0 0 1 15-6.7L21 8" />
              <path d="M3 22v-6h6" />
              <path d="M21 12a9 9 0 0 1-15 6.7L3 16" />
            </svg>
          </ActionIcon>
        </Group>
      </AppShell.Header>

      <AppShell.Navbar p="md">
        <Sidebar
          currentView={currentView}
          setCurrentView={setCurrentView}
          workspaces={workspaces || {}}
          selectedWorkspace={selectedWorkspace}
          onSelectWorkspace={setSelectedWorkspace}
          onClickWorkspaceInfo={handleWorkspaceInfoClick}
          visibleExecutables={executables}
          onSelectExecutable={() => {}}
        />
      </AppShell.Navbar>

      <AppShell.Main>
        {isLoading ? (
          <Group justify="center" h="100%">
            <Loader size="md" />
          </Group>
        ) : hasError ? (
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
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="16"
                height="16"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
              >
                <path d="M21 2v6h-6" />
                <path d="M3 12a9 9 0 0 1 15-6.7L21 8" />
                <path d="M3 22v-6h6" />
                <path d="M21 12a9 9 0 0 1-15 6.7L3 16" />
              </svg>
            </ActionIcon>
          </div>
        ) : selectedWorkspace && workspaces ? (
          <div>
            <Text size="xl" fw={700} mb="md">
              {workspaces[selectedWorkspace].displayName}
            </Text>
            {workspaces[selectedWorkspace].description && (
              <Text mb="lg">{workspaces[selectedWorkspace].description}</Text>
            )}
            <div
              style={{
                padding: "20px",
                backgroundColor: "var(--mantine-color-gray-0)",
                borderRadius: "8px",
                minHeight: "400px",
              }}
            >
              <Text c="dimmed">Select a workspace to view its contents</Text>
            </div>
          </div>
        ) : (
          <div
            style={{
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              height: "100%",
            }}
          >
            <Text c="dimmed">No workspaces available</Text>
          </div>
        )}
      </AppShell.Main>

      {selectedWorkspace && workspaces && (
        <WorkspaceInfoModal
          workspace={workspaces[selectedWorkspace]}
          workspaceId={selectedWorkspace}
          isOpen={isWorkspaceInfoOpen}
          onClose={() => setIsWorkspaceInfoOpen(false)}
        />
      )}
    </AppShell>
  );
}

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <AppContent />
    </QueryClientProvider>
  );
}

export default App;
