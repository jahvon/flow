import { ScrollArea, Text } from "@mantine/core";
import Data from "../../pages/Data/Data";
import Executable from "../../pages/Executable/Executable";
import { Settings } from "../../pages/Settings/Settings";
import { Welcome } from "../../pages/Welcome/Welcome";
import { Workspace } from "../../pages/Workspace/Workspace";
import type { EnrichedExecutable } from "../../types/executable";
import { EnrichedWorkspace } from "../../types/workspace";

export enum View {
  Welcome = "welcome",
  Workspace = "workspace",
  Executable = "executable",
  Logs = "logs",
  Data = "data",
  Settings = "settings",
}

interface ViewerProps {
  currentView: View;
  selectedExecutable: EnrichedExecutable | null;
  executableError: Error | null;
  welcomeMessage?: string;
  workspace: EnrichedWorkspace | null;
}

export function Viewer({
  currentView,
  selectedExecutable,
  executableError,
  welcomeMessage,
  workspace,
}: ViewerProps) {
  const renderContent = () => {
    switch (currentView) {
      case View.Workspace:
        return <Workspace workspace={workspace} />;
      case View.Executable:
        if (selectedExecutable) {
          if (executableError) {
            console.error(executableError);
            return (
              <Text c="red">
                Error loading executable: {executableError.message}
              </Text>
            );
          }
          return <Executable executable={selectedExecutable} />;
        }
        return (
          <Welcome welcomeMessage="Select an executable to get started." />
        );
      case View.Welcome:
        return <Welcome welcomeMessage={welcomeMessage} />;
      case View.Logs:
        return <Text>Logs view coming soon...</Text>;
      case View.Data:
        return <Data />;
      case View.Settings:
        return <Settings />;
      default:
        return <Welcome welcomeMessage={welcomeMessage} />;
    }
  };

  return (
    <ScrollArea
      h="calc(100vh - var(--app-header-height) - var(--app-shell-padding-total))"
      w="calc(100vw - var(--app-navbar-width) - var(--app-shell-padding-total))"
      type="auto"
      scrollbarSize={6}
      scrollHideDelay={100}
      offsetScrollbars
    >
      {renderContent()}
    </ScrollArea>
  );
}
