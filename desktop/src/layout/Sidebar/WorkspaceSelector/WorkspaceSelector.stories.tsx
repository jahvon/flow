import type { Meta, StoryObj } from "@storybook/react";
import { WorkspaceSelector } from "./WorkspaceSelector";
import { Workspace } from "../../../types/generated/workspace";
import { useState } from "react";

const meta: Meta<typeof WorkspaceSelector> = {
  title: "Components/Sidebar/WorkspaceSelector",
  component: WorkspaceSelector,
  parameters: {
    layout: "centered",
  },
  tags: ["autodocs"],
};

export default meta;
type Story = StoryObj<typeof WorkspaceSelector>;

const mockWorkspaces = {
  dotfiles: {
    displayName: "Dotfiles",
  } as Workspace,
  app: {
    displayName: "Web App",
  } as Workspace,
};

// Example of how to use the component with state management
const WorkspaceSelectorWithState = () => {
  const [selected, setSelected] = useState("dotfiles");

  return (
    <WorkspaceSelector
      workspaces={mockWorkspaces}
      selectedWorkspace={selected}
      onSelectWorkspace={setSelected}
      onClickWorkspaceInfo={(workspace) =>
        console.log("Clicked workspace info:", workspace)
      }
    />
  );
};

export const Default: Story = {
  render: () => <WorkspaceSelectorWithState />,
};

export const NoWorkspaces: Story = {
  args: {
    workspaces: {},
    selectedWorkspace: "",
    onSelectWorkspace: (workspace) =>
      console.log("Selected workspace:", workspace),
    onClickWorkspaceInfo: (workspace) =>
      console.log("Clicked workspace info:", workspace),
  },
};
