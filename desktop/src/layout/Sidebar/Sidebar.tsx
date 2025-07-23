import { Group, Image, NavLink, Stack } from "@mantine/core";
import { ExecutableTree } from "./ExecutableTree/ExecutableTree";
import styles from "./Sidebar.module.css";
import { WorkspaceSelector } from "./WorkspaceSelector/WorkspaceSelector";
import iconImage from "/logo-dark.png";
import {IconDatabase, IconFolders, IconLogs, IconSettings} from "@tabler/icons-react";
import {Link, useLocation} from "react-router-dom";
import {useAppContext} from "../../hooks/useAppContext.tsx";

export function Sidebar() {
  const location = useLocation();
  const { executables, selectedWorkspace } = useAppContext()

  return (
    <div className={styles.sidebar}>
      <Link to="/" className={styles.sidebar__logo}>
        <Image
          src={iconImage}
          alt="flow"
          fit="contain"
        />
      </Link>
      <Stack gap="xs">
        <WorkspaceSelector />

        <Group gap="xs" mt="md">
          <NavLink
              label="Workspace"
              leftSection={<IconFolders size={16} />}
              component={Link}
              to={`/workspace/${selectedWorkspace}`}
              active={location.pathname.startsWith('/workspace')}
              variant="filled"
          />
          <NavLink
              label="Logs"
              leftSection={<IconLogs size={16} />}
              component={Link}
              to={`/logs`}
              active={location.pathname.startsWith('/logs')}
              variant="filled"
          />
          <NavLink
              label="Data"
              leftSection={<IconDatabase size={16} />}
              variant="filled"
              childrenOffset={28}
          >
            <NavLink label="Cache" component={Link} to={`/cache`} variant="filled" active={location.pathname.startsWith('/cache')}/>
            <NavLink label="Vault" component={Link} to={`/vault`} variant="filled" active={location.pathname.startsWith('/vault')}/>
          </NavLink>
          <NavLink
              label="Settings"
              leftSection={<IconSettings size={16} />}
              component={Link}
              to={`/settings`}
              active={location.pathname.startsWith('/settings')}
              variant="filled"
          />
        </Group>

        {executables && executables.length > 0 && (
          <ExecutableTree />
        )}
      </Stack>
    </div>
  );
}
