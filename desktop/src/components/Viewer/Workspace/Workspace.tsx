import { Button, Group, Stack, Text, Title } from "@mantine/core";
import type { Workspace } from "../../../types/generated/workspace";

interface WorkspaceProps {
  workspace: Workspace | null;
  workspaceId: string | null;
  onClose: () => void;
}

export function Workspace({ workspace, workspaceId, onClose }: WorkspaceProps) {
  if (!workspace) {
    return null;
  }

  return (
    <Stack p="xl">
      <Group justify="space-between" align="center">
        <Title order={2}>{workspace.displayName || "Workspace"}</Title>
        <Button variant="subtle" onClick={onClose}>
          Back
        </Button>
      </Group>

      {workspace.description && (
        <Text size="lg" c="dimmed">
          {workspace.description}
        </Text>
      )}

      {workspaceId && (
        <Group>
          <Text fw={500}>ID:</Text>
          <Text>{workspaceId}</Text>
        </Group>
      )}

      {workspace.tags && workspace.tags.length > 0 && (
        <Group>
          <Text fw={500}>Tags:</Text>
          <Text>{workspace.tags.join(", ")}</Text>
        </Group>
      )}
    </Stack>
  );
}
