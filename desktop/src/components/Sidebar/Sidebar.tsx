import { ActionIcon, Group, rem, Stack, Text, Tooltip } from "@mantine/core";
import { EnrichedExecutable } from "../../types/executable";
import { WorkspaceMap } from "../../types/workspace";
import { View, ViewLinks } from "../Viewer/Viewer";
import { ExecutableTree } from "./ExecutableTree/ExecutableTree";
import { WorkspaceSelector } from "./WorkspaceSelector/WorkspaceSelector";

interface SidebarProps {
  currentView: View;
  setCurrentView: (view: View) => void;

  workspaces: WorkspaceMap;
  selectedWorkspace: string | null;
  onSelectWorkspace: (workspace: string) => void;
  onClickWorkspaceInfo: (workspace: string) => void;

  visibleExecutables: EnrichedExecutable[];
  onSelectExecutable: (executable: string) => void;
}

export function Sidebar({
  currentView,
  setCurrentView,
  workspaces,
  selectedWorkspace,
  onSelectWorkspace,
  onClickWorkspaceInfo,
  visibleExecutables,
  onSelectExecutable,
}: SidebarProps) {
  const renderSecondaryNav = () => {
    switch (currentView) {
      case View.Workspaces:
        return (
          <>
            <WorkspaceSelector
              workspaces={workspaces}
              selectedWorkspace={selectedWorkspace ?? ""}
              onSelectWorkspace={onSelectWorkspace}
              onClickWorkspaceInfo={onClickWorkspaceInfo}
            />
            <ExecutableTree
              visibleExecutables={visibleExecutables}
              onSelectExecutable={onSelectExecutable}
            />
          </>
        );

      case View.Logs:
        return (
          <>
            <Text size="xs" fw={700} c="dimmed" mb="md">
              LOGS
            </Text>
            <Text size="xs" fw={700} c="dimmed" mb="md">
              To Be Implemented
            </Text>
          </>
        );

      case View.Settings:
        return (
          <>
            <Text size="xs" fw={700} c="dimmed" mb="md">
              SETTINGS
            </Text>
            <Stack>
              <Text size="xs" c="dimmed">
                Settings content will go here
              </Text>
            </Stack>
          </>
        );
    }
  };

  return (
    <Group h="100%" gap={0}>
      <Stack
        justify="flex-start"
        style={{
          width: rem(35),
          height: "100%",
          alignItems: "center",
        }}
      >
        {ViewLinks.map((link) => (
          <Tooltip label={link.label} position="right">
            <ActionIcon
              onClick={() => setCurrentView(link.view)}
              variant={currentView === link.view ? "filled" : "transparent"}
              size="lg"
              title={link.label}
              style={{
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
              }}
            >
              <link.icon
                style={{ width: rem(20), height: rem(20) }}
                stroke={1.5}
              />
            </ActionIcon>
          </Tooltip>
        ))}
      </Stack>

      <Stack
        style={{
          flex: 1,
          height: "100%",
          backgroundColor: "var(--mantine-color-dark-7)",
          padding: "var(--mantine-spacing-md)",
        }}
      >
        {renderSecondaryNav()}
      </Stack>
    </Group>
  );
}
