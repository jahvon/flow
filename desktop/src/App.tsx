import { useState, useEffect } from "react";
import { AppShell, Text, UnstyledButton, Group, Loader, ActionIcon } from '@mantine/core';
import { invoke } from "@tauri-apps/api/core";
import '@mantine/core/styles.css';
import "./App.css";
import { Workspace } from './types/workspace';

function App() {
  const [workspaces, setWorkspaces] = useState<Workspace[]>([]);
  const [selectedWorkspace, setSelectedWorkspace] = useState<Workspace | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const loadWorkspaces = async () => {
    setIsLoading(true);
    try {
      const response = await invoke<{ workspaces: Workspace[] }>("workspaces");
      setWorkspaces(response.workspaces);
      // Select the first workspace by default if none is selected
      if (!selectedWorkspace && response.workspaces.length > 0) {
        setSelectedWorkspace(response.workspaces[0]);
      }
    } catch (error) {
      console.error('Failed to load workspaces:', error);
    } finally {
      setIsLoading(false);
    }
  };

  // Load workspaces when the component mounts
  useEffect(() => {
    loadWorkspaces();
  }, []);

  return (
    <AppShell
      header={{ height: 60 }}
      navbar={{ width: 300, breakpoint: 'sm' }}
      padding="md"
    >
      <AppShell.Header>
        <Group h="100%" px="md" justify="space-between">
          <Text size="lg" fw={500}>Flow</Text>
          <ActionIcon
            variant="subtle"
            onClick={loadWorkspaces}
            loading={isLoading}
            title="Refresh workspaces"
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
        {isLoading ? (
          <Group justify="center" h="100%">
            <Loader size="sm" />
          </Group>
        ) : (
          workspaces.map((workspace, index) => (
            <UnstyledButton
              key={index}
              onClick={() => setSelectedWorkspace(workspace)}
              style={{
                display: 'block',
                width: '100%',
                padding: '8px 12px',
                borderRadius: '4px',
                backgroundColor: selectedWorkspace === workspace
                  ? 'var(--mantine-color-blue-1)'
                  : 'transparent',
                '&:hover': {
                  backgroundColor: 'var(--mantine-color-gray-1)',
                },
              }}
            >
              <Text size="sm" fw={500}>{workspace.displayName}</Text>
              {workspace.description && (
                <Text size="xs" c="dimmed">{workspace.description}</Text>
              )}
            </UnstyledButton>
          ))
        )}
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
