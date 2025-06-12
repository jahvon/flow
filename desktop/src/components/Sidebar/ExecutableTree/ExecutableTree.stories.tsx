import type { Meta, StoryObj } from "@storybook/react";
import { ExecutableTree } from "./ExecutableTree";
import { EnrichedExecutable } from "../../../types/executable";

const meta = {
  title: "Components/Sidebar/ExecutableTree",
  component: ExecutableTree,
  parameters: {
    layout: "centered",
  },
  tags: ["autodocs"],
} satisfies Meta<typeof ExecutableTree>;

export default meta;
type Story = StoryObj<typeof meta>;

const sampleExecutables: EnrichedExecutable[] = [
  {
    id: "run script",
    name: "script",
    namespace: null,
    verb: "run",
    workspace: "default",
    flowfile: "exec.flow",
  },
  {
    id: "validate deps",
    name: "deps",
    namespace: null,
    verb: "validate",
    workspace: "default",
    flowfile: "exec.flow",
  },
  {
    id: "exec devserver",
    name: "devserver",
    namespace: "frontend",
    workspace: "default",
    flowfile: "frontend.flow",
    verb: "exec",
  },
  {
    id: "test frontend",
    name: "frontend",
    namespace: "frontend",
    workspace: "default",
    flowfile: "frontend.flow",
    verb: "test",
  },
  {
    id: "build container",
    name: "container",
    namespace: "backend",
    workspace: "default",
    flowfile: "backend.flow",
    verb: "build",
  },
];

export const Default: Story = {
  args: {
    visibleExecutables: sampleExecutables,
    onSelectExecutable: (id) => console.log("Selected executable:", id),
  },
};

export const Empty: Story = {
  args: {
    visibleExecutables: [],
    onSelectExecutable: (id) => console.log("Selected executable:", id),
  },
};
