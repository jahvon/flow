import { ActionIcon } from "@mantine/core";
import { IconPlus, IconRefresh } from "@tabler/icons-react";
import styles from "./Header.module.css";

interface HeaderProps {
  onCreateWorkspace: () => void;
  onRefreshWorkspaces: () => void;
}

export function Header({
  onCreateWorkspace,
  onRefreshWorkspaces,
}: HeaderProps) {
  return (
    <div className={styles.header}>
      <div className={styles.header__actions}>
        <ActionIcon
          className={styles.header__button}
          onClick={onCreateWorkspace}
          title="Create workspace"
          variant="light"
        >
          <IconPlus size={16} />
        </ActionIcon>

        <ActionIcon
          className={styles.header__button}
          variant="light"
          onClick={onRefreshWorkspaces}
          title="Refresh workspaces"
        >
          <IconRefresh size={16} />
        </ActionIcon>
      </div>
    </div>
  );
}
