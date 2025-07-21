import { ComboboxItem, Group, OptionsFilter, Select } from "@mantine/core";
import { useConfig } from "../../../hooks/useBackendData";
import { useNotifier } from "../../../hooks/useNotifier";
import { EnrichedWorkspace } from "../../../types/workspace";
import { NotificationType } from "../../../types/notification";

interface WorkspaceSelectorProps {
  workspaces: EnrichedWorkspace[];
  selectedWorkspace: string | null;
  onSelectWorkspace: (workspaceName: string) => void;
}

const filter: OptionsFilter = ({ options, search }) => {
  const filtered = (options as ComboboxItem[]).filter((option) =>
    option.label.toLowerCase().trim().includes(search.toLowerCase().trim())
  );

  filtered.sort((a, b) => a.label.localeCompare(b.label));
  return filtered;
};

export function WorkspaceSelector({
  workspaces,
  selectedWorkspace,
  onSelectWorkspace,
}: WorkspaceSelectorProps) {
  const { config, updateCurrentWorkspace } = useConfig();
  const { setNotification } = useNotifier();

  const handleWorkspaceChange = async (workspaceName: string) => {
    onSelectWorkspace(workspaceName);

    if (config?.workspaceMode === 'dynamic') {
      try {
        await updateCurrentWorkspace(workspaceName);
        setNotification({
          type: NotificationType.Success,
          title: 'Workspace switched',
          message: `Switched to workspace: ${workspaceName}`,
          autoClose: true,
        });
      } catch (error) {
        setNotification({
          type: NotificationType.Error,
          title: 'Error switching workspace',
          message: error.message,
          autoClose: true,
        });
      }
    }
  };

  return (
    <>
      <Group gap="xs" align="center" wrap="nowrap">
        <Select
          value={selectedWorkspace}
          onChange={(value) => {
            if (value && workspaces.find((w) => w.name === value)) {
              handleWorkspaceChange(value);
            }
          }}
          data={workspaces.map((workspace) => ({
            value: workspace.name,
            label: workspace.displayName || workspace.name,
          }))}
          filter={filter}
          placeholder="No workspace selected"
          searchable
          size="sm"
          styles={{
            root: {
              flex: 1,
            },
            input: {
              backgroundColor: "var(--mantine-color-dark-6)",
              borderColor: "var(--mantine-color-dark-4)",
              color: "var(--mantine-color-white)",
            },
            dropdown: {
              backgroundColor: "var(--mantine-color-dark-6)",
              borderColor: "var(--mantine-color-dark-4)",
            },
            option: {
              color: "var(--mantine-color-white)",
              "&[data-selected]": {
                backgroundColor: "var(--mantine-color-dark-5)",
              },
              "&[data-hovered]": {
                backgroundColor: "var(--mantine-color-dark-5)",
              },
            },
          }}
        />
      </Group>
    </>
  );
}
