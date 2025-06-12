import { AppShell, Group } from "@mantine/core";
import { ActionButtons } from "./ActionButtons/ActionButtons";

interface HeaderProps {
  onCreateWorkspace: () => void;
  onRefreshWorkspaces: () => void;
}

export function Header({ onCreateWorkspace, onRefreshWorkspaces }: HeaderProps) {
  return (
    <AppShell.Header>
      <Group h="100%" px="sm" justify="space-between">
        <img
          src={`/logo-${window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'}.png`}
          alt="Flow"
          height="24"
        />
        <ActionButtons onCreateWorkspace={onCreateWorkspace} onRefreshWorkspaces={onRefreshWorkspaces} />
      </Group>
    </AppShell.Header>

  )
}