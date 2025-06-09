import { useState, useEffect } from "react";
import { AppShell, Text, Group, Loader, ActionIcon } from '@mantine/core';
import { invoke } from "@tauri-apps/api/core";
import '@mantine/core/styles.css';
import "./App.css";
import { Workspace } from './types/workspace';
import { Sidebar } from './components/Sidebar/Sidebar';
import { Config } from "./types/config";

interface WorkspaceMap {
  [key: string]: Workspace;
}

function App() {
  const [config, setConfig] = useState<Config | null>(null);
  const [workspaces, setWorkspaces] = useState<WorkspaceMap>({});
  const [selectedWorkspace, setSelectedWorkspace] = useState<Workspace | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const loadWorkspaces = async () => {
    setIsLoading(true);
    try {
      const cfg = await invoke<Config>("get_config");
      setConfig(cfg);
    } catch (error) {
      console.error('Failed to load config:', error);
    }

    try {
      const response = await invoke<{ workspaces: WorkspaceMap }>("list_workspaces");
      setWorkspaces(response.workspaces);
      if (!selectedWorkspace && Object.keys(response.workspaces).length > 0) {
        setSelectedWorkspace(Object.values(response.workspaces)[0]);
      }
    } catch (error) {
      console.error('Failed to load workspaces:', error);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    loadWorkspaces();
  }, []);

  return (
    <AppShell
      header={{ height: 60 }}
      navbar={{ width: 360, breakpoint: 'sm' }}
      padding="md"
    >
      <AppShell.Header>
        <Group h="100%" px="md" justify="space-between">
          <img
            src={`/logo-${window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'}.png`}
            alt="Flow"
            height="32"
          />
          <ActionIcon
            color="palette"
            onClick={loadWorkspaces}
            loading={isLoading}
            title="Refresh workspaces"
            size="lg"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
              <path d="M21 2v6h-6"/>
              <path d="M3 12a9 9 0 0 1 15-6.7L21 8"/>
              <path d="M3 22v-6h6"/>
              <path d="M21 12a9 9 0 0 1-15 6.7L3 16"/>
            </svg>
          </ActionIcon>
        </Group>
      </AppShell.Header>

      <AppShell.Navbar p="md">
        <Sidebar
          config={config}
          workspaces={workspaces}
          selectedWorkspace={selectedWorkspace}
          onSelectWorkspace={setSelectedWorkspace}
          isLoading={isLoading}
        />
      </AppShell.Navbar>

      <AppShell.Main>
        {isLoading ? (
          <Group justify="center" h="100%">
            <Loader size="md" />
          </Group>
        ) : selectedWorkspace ? (
          <div>
            <Text size="xl" fw={700} mb="md">{selectedWorkspace.displayName}</Text>
            {selectedWorkspace.description && (
              <Text mb="lg">{selectedWorkspace.description}</Text>
            )}
            <div style={{
              padding: '20px',
              backgroundColor: 'var(--mantine-color-gray-0)',
              borderRadius: '8px',
              minHeight: '400px'
            }}>
              <Text c="dimmed">Select a workspace to view its contents</Text>
            </div>
          </div>
        ) : (
          <div style={{
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            height: '100%'
          }}>
            <Text c="dimmed">No workspaces available</Text>
          </div>
        )}
      </AppShell.Main>
    </AppShell>
  );
}

export default App;
