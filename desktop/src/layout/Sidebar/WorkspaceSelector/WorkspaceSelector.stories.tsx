import type { Meta, StoryObj } from "@storybook/react";
import { useState } from "react";
import { EnrichedWorkspace } from "../../../types/workspace";
import { WorkspaceSelector } from "./WorkspaceSelector";

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

const mockWorkspaces: EnrichedWorkspace[] = [
  {
    name: "Dotfiles",
    path: "/Users/flow/.dotfiles",
    fullDescription: "Dotfiles workspace",
  },
  {
    name: "Web App",
    path: "/Users/flow/WebApp",
    fullDescription: "Web App workspace",
  },
];

// Example of how to use the component with state management
const WorkspaceSelectorWithState = () => {
  const [selected, setSelected] = useState("dotfiles");

  return (
    <WorkspaceSelector
      workspaces={mockWorkspaces}
      selectedWorkspace={selected}
      onSelectWorkspace={setSelected}
    />
  );
};

export const Default: Story = {
  render: () => <WorkspaceSelectorWithState />,
};

export const NoWorkspaces: Story = {
  args: {
    workspaces: [],
    selectedWorkspace: "",
    onSelectWorkspace: (workspace) =>
      console.log("Selected workspace:", workspace),
  },
};
