import { Group, Loader, Text, UnstyledButton, Stack, rem, Select } from '@mantine/core';
import { Workspace } from '../../types/workspace';
import { IconHome, IconFolder, IconSettings } from '@tabler/icons-react';
import { useState, useEffect } from 'react';
import { Config } from '../../types/config';
import { invoke } from '@tauri-apps/api/core';
import { Executable } from '../../types/flowfile';

interface WorkspaceMap {
  [key: string]: Workspace;
}

interface SidebarProps {
  config: Config | null;
  workspaces: WorkspaceMap;
  selectedWorkspace: Workspace | null;
  onSelectWorkspace: (workspace: Workspace) => void;
  isLoading: boolean;
}

enum NavSection {
  CurrentWorkspace = 'current-workspace',
  AllWorkspaces = 'all-workspaces',
  Settings = 'settings',
}

const mainLinks = [
  { icon: IconHome, label: 'Current Workspace', section: NavSection.CurrentWorkspace },
  { icon: IconFolder, label: 'All Workspaces', section: NavSection.AllWorkspaces },
  { icon: IconSettings, label: 'Settings', section: NavSection.Settings },
];

export function Sidebar({ config, workspaces, selectedWorkspace, onSelectWorkspace, isLoading }: SidebarProps) {
  const [activeSection, setActiveSection] = useState<NavSection>(NavSection.CurrentWorkspace);
  const [executables, setExecutables] = useState<Executable[]>([]);
  const [isLoadingExecutables, setIsLoadingExecutables] = useState(false);
  const [executableError, setExecutableError] = useState<string | null>(null);

  useEffect(() => {
    if (config?.currentWorkspace) {
      setIsLoadingExecutables(true);
      setExecutableError(null);

      // fetch executables for the current workspace
      invoke<Executable[]>("list_executables", { workspace: config.currentWorkspace })
        .then(result => {
            setExecutables(result);
        })
        .catch(err => {
          console.error('Failed to fetch executables:', err);
          setExecutableError('Failed to load executables');
          setExecutables([]);
        })
        .finally(() => {
          setIsLoadingExecutables(false);
        });
    } else {
      setExecutables([]);
      setExecutableError(null);
    }
  }, [config?.currentWorkspace]);

  if (isLoading) {
    return (
      <Group justify="center" h="100%">
        <Loader size="sm" />
      </Group>
    );
  }

  const renderSecondaryNav = () => {
    switch (activeSection) {
      case NavSection.CurrentWorkspace:
        return (
          <>
            <Text size="xs" fw={700} c="dimmed" mb="md">CURRENT WORKSPACE</Text>
            <Stack>
              <Select
                value={config?.currentWorkspace || ''}
                onChange={(value) => {
                  if (value && workspaces[value]) {
                    onSelectWorkspace(workspaces[value]);
                  }
                }}
                data={Object.entries(workspaces).map(([name, workspace]) => ({
                  value: name,
                  label: workspace.displayName || name,
                }))}
                placeholder="Select a workspace"
                searchable
                size="sm"
                styles={{
                  input: {
                    backgroundColor: 'var(--mantine-color-dark-6)',
                    borderColor: 'var(--mantine-color-dark-4)',
                    color: 'var(--mantine-color-white)',
                  },
                  dropdown: {
                    backgroundColor: 'var(--mantine-color-dark-6)',
                    borderColor: 'var(--mantine-color-dark-4)',
                  },
                  option: {
                    color: 'var(--mantine-color-white)',
                    '&[data-selected]': {
                      backgroundColor: 'var(--mantine-color-dark-5)',
                    },
                    '&[data-hovered]': {
                      backgroundColor: 'var(--mantine-color-dark-5)',
                    },
                  },
                }}
              />
              <Text size="xs" fw={700} c="dimmed" mt="md">EXECUTABLES</Text>
              {isLoadingExecutables ? (
                <Group justify="center">
                  <Loader size="xs" />
                </Group>
              ) : executableError ? (
                <Text size="xs" c="red">{executableError}</Text>
              ) : executables.length === 0 ? (
                <Text size="xs" c="dimmed">No executables found</Text>
              ) : (
                executables.map((exec: Executable) => (
                  <Text size="xs" c="dimmed" key={exec.name}>{exec.verb} {exec.name}</Text>
                ))
              )}
            </Stack>
          </>
        );

      case NavSection.AllWorkspaces:
        return (
          <>
            <Text size="xs" fw={700} c="dimmed" mb="md">ALL WORKSPACES</Text>
            {Object.entries(workspaces).map(([name, workspace]: [string, Workspace]) => (
              <UnstyledButton
                key={name}
                onClick={() => onSelectWorkspace(workspace)}
                style={{
                  display: 'block',
                  width: '100%',
                  padding: '8px 12px',
                  borderRadius: '4px',
                  backgroundColor: selectedWorkspace === workspace
                    ? 'var(--mantine-color-dark-5)'
                    : 'transparent',
                  '&:hover': {
                    backgroundColor: 'var(--mantine-color-dark-5)',
                  },
                }}
              >
                <Text size="sm" fw={500} c="white">{workspace.displayName}</Text>
                {workspace.description && (
                  <Text size="xs" c="dimmed">{workspace.description}</Text>
                )}
              </UnstyledButton>
            ))}
          </>
        );

      case NavSection.Settings:
        return (
          <>
            <Text size="xs" fw={700} c="dimmed" mb="md">SETTINGS</Text>
            <Stack>
              <Text size="xs" c="dimmed">Settings content will go here</Text>
            </Stack>
          </>
        );
    }
  };

  return (
    <Group h="100%" gap={0}>
      {/* Main Navigation */}
      <Stack
        justify="flex-start"
        style={{
          width: rem(60),
          height: '100%',
          backgroundColor: 'var(--mantine-color-dark-6)',
        }}
      >
        {mainLinks.map((link) => (
          <UnstyledButton
            key={link.label}
            onClick={() => setActiveSection(link.section)}
            style={{
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              width: rem(60),
              height: rem(60),
              borderRadius: 'var(--mantine-radius-sm)',
              color: 'var(--mantine-color-white)',
              backgroundColor: activeSection === link.section
                ? 'var(--mantine-color-dark-5)'
                : 'transparent',
              '&:hover': {
                backgroundColor: 'var(--mantine-color-dark-5)',
              },
            }}
          >
            <link.icon style={{ width: rem(22), height: rem(22) }} stroke={1.5} />
          </UnstyledButton>
        ))}
      </Stack>

      {/* Sub Navigation */}
      <Stack
        style={{
          flex: 1,
          height: '100%',
          backgroundColor: 'var(--mantine-color-dark-7)',
          padding: 'var(--mantine-spacing-md)',
        }}
      >
        {renderSecondaryNav()}
      </Stack>
    </Group>
  );
}