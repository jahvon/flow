import type { Meta, StoryObj } from "@storybook/react";
import { EnrichedExecutable } from "../../../types/executable";
import { ExecutableTree } from "./ExecutableTree";

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
    id: "script",
    ref: "run script",
    name: "script",
    namespace: null,
    verb: "run",
    workspace: "default",
    flowfile: "exec.flow",
  },
  {
    id: "deps",
    ref: "validate deps",
    name: "deps",
    namespace: null,
    verb: "validate",
    workspace: "default",
    flowfile: "exec.flow",
  },
  {
    id: "devserver",
    ref: "exec frontend/devserver",
    name: "devserver",
    namespace: "frontend",
    workspace: "default",
    flowfile: "frontend.flow",
    verb: "exec",
  },
  {
    id: "frontend",
    ref: "test frontend",
    name: "",
    namespace: "frontend",
    workspace: "default",
    flowfile: "frontend.flow",
    verb: "test",
  },
  {
    id: "container",
    ref: "build backend/container",
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
    onSelectExecutable: (ref) => console.log("Selected executable:", ref),
  },
};

export const Empty: Story = {
  args: {
    visibleExecutables: [],
    onSelectExecutable: (ref) => console.log("Selected executable:", ref),
  },
};
