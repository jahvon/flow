import { Text } from "@mantine/core";
import { IconFolders, IconLogs, IconSettings } from "@tabler/icons-react";
import type { EnrichedExecutable } from "../../types/executable";
import { EnrichedWorkspace } from "../../types/workspace";
import ExecutableInfo from "./Executable/Executable";
import { Settings } from "./Settings/Settings";
import { Welcome } from "./Welcome/Welcome";
import { Workspace as WorkspaceView } from "./Workspace/Workspace";

export enum View {
  Welcome = "welcome",
  Workspace = "workspace",
  Executable = "executable",
  Logs = "logs",
  Settings = "settings",
}

export const ViewLinks = [
  { icon: IconFolders, label: "Workspace", view: View.Workspace },
  { icon: IconLogs, label: "Logs", view: View.Logs },
  { icon: IconSettings, label: "Settings", view: View.Settings },
];

interface ViewerProps {
  currentView: View;
  selectedExecutable: EnrichedExecutable | null;
  isExecutableLoading: boolean;
  executableError: Error | null;
  welcomeMessage?: string;
  workspace: EnrichedWorkspace | null;
}

export function Viewer({
  currentView,
  selectedExecutable,
  isExecutableLoading,
  executableError,
  welcomeMessage,
  workspace,
}: ViewerProps) {
  const renderContent = () => {
    switch (currentView) {
      case View.Workspace:
        return <WorkspaceView workspace={workspace} />;
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
          return <ExecutableInfo executable={selectedExecutable} />;
        }
        return (
          <Welcome welcomeMessage="Select an executable to get started." />
        );
      case View.Welcome:
        return <Welcome welcomeMessage={welcomeMessage} />;
      case View.Logs:
        return <Text>Logs view coming soon...</Text>;
      case View.Settings:
        return <Settings />;
      default:
        return <Welcome welcomeMessage={welcomeMessage} />;
    }
  };

  return renderContent();
}
