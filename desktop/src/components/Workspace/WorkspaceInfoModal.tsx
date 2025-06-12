import { Button, Group, Modal, Stack, Text, Title } from "@mantine/core";
import { Workspace } from "../../types/generated/workspace";

interface WorkspaceInfoModalProps {
  workspace: Workspace | null;
  workspaceId: string | null;
  isOpen: boolean;
  onClose: () => void;
}

export function WorkspaceInfoModal({
  workspace,
  workspaceId,
  isOpen,
  onClose,
}: WorkspaceInfoModalProps) {
  if (!workspace) {
    return null;
  }

  return (
    <Modal
      opened={isOpen}
      onClose={onClose}
      title={<Title order={3}>{workspace.displayName || "Workspace"}</Title>}
      size="md"
    >
      <Stack>
        {workspace.description && <Text>{workspace.description}</Text>}

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

        <Button onClick={onClose} fullWidth mt="md">
          Close
        </Button>
      </Stack>
    </Modal>
  );
}
