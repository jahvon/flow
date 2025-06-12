import { IconFolders, IconLogs, IconSettings } from "@tabler/icons-react";

export enum View {
  Workspaces = "workspaces",
  Logs = "logs",
  Settings = "settings",
}

export const ViewLinks = [
  { icon: IconFolders, label: "Workspaces", view: View.Workspaces },
  { icon: IconLogs, label: "Logs", view: View.Logs },
  { icon: IconSettings, label: "Settings", view: View.Settings },
];
