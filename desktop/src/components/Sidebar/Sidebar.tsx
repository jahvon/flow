import { Group, Image, NavLink, Stack } from "@mantine/core";
import type { EnrichedExecutable } from "../../types/executable";
import type { Workspace } from "../../types/generated/workspace";
import { View, ViewLinks } from "../Viewer/Viewer";
import { ExecutableTree } from "./ExecutableTree/ExecutableTree";
import styles from "./Sidebar.module.css";
import { WorkspaceSelector } from "./WorkspaceSelector/WorkspaceSelector";
import iconImage from "/logo-dark.png";

interface SidebarProps {
  currentView: View;
  setCurrentView: (view: View) => void;
  workspaces: Record<string, Workspace>;
  selectedWorkspace: string | null;
  onSelectWorkspace: (workspaceId: string) => void;
  onClickWorkspaceInfo: () => void;
  visibleExecutables: EnrichedExecutable[];
  onSelectExecutable: (executableId: string) => void;
  onLogoClick: () => void;
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
  onLogoClick,
}: SidebarProps) {
  return (
    <div className={styles.sidebar}>
      <div className={styles.sidebar__logo}>
        <Image
          src={iconImage}
          alt="flow"
          fit="contain"
          onClick={onLogoClick}
          style={{ cursor: "pointer" }}
        />
      </div>
      <Stack gap="xs">
        <WorkspaceSelector
          workspaces={workspaces}
          selectedWorkspace={selectedWorkspace}
          onSelectWorkspace={onSelectWorkspace}
          onClickWorkspaceInfo={onClickWorkspaceInfo}
        />

        <Group gap="xs" mt="md">
          {ViewLinks.map((link) => (
            <NavLink
              key={link.view}
              label={link.label}
              leftSection={<link.icon size={16} />}
              active={currentView === link.view}
              onClick={() => setCurrentView(link.view)}
              variant="filled"
            />
          ))}
        </Group>

        {visibleExecutables && visibleExecutables.length > 0 && (
          <ExecutableTree
            visibleExecutables={visibleExecutables}
            onSelectExecutable={onSelectExecutable}
          />
        )}
      </Stack>
    </div>
  );
}
