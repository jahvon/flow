import {
  IconDatabase,
  IconFolders,
  IconLogs,
  IconSettings,
} from "@tabler/icons-react";
import { View } from "./Viewer";

export const ViewLinks = [
  { icon: IconFolders, label: "Workspace", view: View.Workspace },
  { icon: IconLogs, label: "Logs", view: View.Logs },
  { icon: IconDatabase, label: "Data", view: View.Data },
  { icon: IconSettings, label: "Settings", view: View.Settings },
];