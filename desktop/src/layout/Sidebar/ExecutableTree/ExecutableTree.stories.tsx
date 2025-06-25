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
    fullDescription: "Run a script",
  },
  {
    id: "deps",
    ref: "validate deps",
    name: "deps",
    namespace: null,
    verb: "validate",
    workspace: "default",
    flowfile: "exec.flow",
    fullDescription: "Validate dependencies",
  },
  {
    id: "devserver",
    ref: "exec frontend/devserver",
    name: "devserver",
    namespace: "frontend",
    workspace: "default",
    flowfile: "frontend.flow",
    verb: "exec",
    fullDescription: "Execute a development server",
  },
  {
    id: "frontend",
    ref: "test frontend",
    name: "",
    namespace: "frontend",
    workspace: "default",
    flowfile: "frontend.flow",
    verb: "test",
    fullDescription: "Test a frontend",
  },
  {
    id: "container",
    ref: "build backend/container",
    name: "container",
    namespace: "backend",
    workspace: "default",
    flowfile: "backend.flow",
    verb: "build",
    fullDescription: "Build a container",
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
