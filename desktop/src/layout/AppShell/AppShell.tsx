import {
  ActionIcon,
  Loader,
  AppShell as MantineAppShell,
  Notification as MantineNotification,
  Text,
} from "@mantine/core";
import { IconRefresh } from "@tabler/icons-react";
import { ReactNode } from "react";
import { EnrichedExecutable } from "../../types/executable";
import {
  colorFromType,
  Notification,
  NotificationType,
} from "../../types/notification";
import { EnrichedWorkspace } from "../../types/workspace";
import { Header } from "../Header/Header";
import { Sidebar } from "../Sidebar/Sidebar";
import { View } from "../Viewer/Viewer";
import styles from "./AppShell.module.css";

interface AppShellProps {
  children: ReactNode;
  currentView: View;
  setCurrentView: (view: View) => void;
  workspaces: EnrichedWorkspace[];
  selectedWorkspace: string | null;
  onSelectWorkspace: (workspaceName: string) => void;
  visibleExecutables: EnrichedExecutable[];
  onSelectExecutable: (executable: string) => void;
  onLogoClick: () => void;
  hasError: Error | null;
  isLoading: boolean;
  refreshAll: () => void;
  notification: Notification | null;
  setNotification: (notification: Notification | null) => void;
}

export function AppShell({
  children,
  currentView,
  setCurrentView,
  workspaces,
  selectedWorkspace,
  onSelectWorkspace,
  visibleExecutables,
  onSelectExecutable,
  onLogoClick,
  hasError,
  isLoading,
  refreshAll,
  notification,
  setNotification,
}: AppShellProps) {
  return (
    <MantineAppShell
      header={{ height: "var(--app-header-height)" }}
      navbar={{ width: "var(--app-navbar-width)", breakpoint: "sm" }}
      padding="md"
      classNames={{
        root: styles.appShell,
        main: styles.main,
        header: styles.header,
        navbar: styles.navbar,
      }}
    >
      <MantineAppShell.Header>
        <Header
          onCreateWorkspace={() => {}}
          onRefreshWorkspaces={() => {
            refreshAll();
            setNotification({
              title: "Refresh completed",
              message: "flow data has synced and refreshed successfully",
              type: NotificationType.Success,
              autoClose: true,
              autoCloseDelay: 3000,
            });
          }}
        />
      </MantineAppShell.Header>

      <MantineAppShell.Navbar>
        <Sidebar
          currentView={currentView}
          setCurrentView={setCurrentView}
          workspaces={workspaces}
          selectedWorkspace={selectedWorkspace}
          onSelectWorkspace={onSelectWorkspace}
          visibleExecutables={visibleExecutables}
          onSelectExecutable={onSelectExecutable}
          onLogoClick={onLogoClick}
        />
      </MantineAppShell.Navbar>

      <MantineAppShell.Main>
        {hasError ? (
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
              <IconRefresh size={16} />
            </ActionIcon>
          </div>
        ) : (
          <div style={{ position: "relative", height: "100%" }}>
            {children}
            {isLoading && (
              <div
                style={{
                  position: "absolute",
                  top: 16,
                  right: 16,
                  zIndex: 1000,
                }}
              >
                <Loader size="sm" />
              </div>
            )}
          </div>
        )}
      </MantineAppShell.Main>

      {notification && (
        <MantineNotification
          title={notification.title}
          onClose={() => setNotification(null)}
          color={colorFromType(notification.type)}
          style={{
            position: "fixed",
            bottom: 20,
            right: 20,
            zIndex: 1000,
          }}
        >
          {notification.message}
        </MantineNotification>
      )}
    </MantineAppShell>
  );
}
