import {
  ActionIcon,
  ComboboxItem,
  Group,
  OptionsFilter,
  Select,
  Text,
} from "@mantine/core";
import { IconInfoCircle } from "@tabler/icons-react";
import { WorkspaceMap } from "../../../types/workspace";

interface WorkspaceSelectorProps {
  workspaces: WorkspaceMap;
  selectedWorkspace: string;
  onSelectWorkspace: (workspace: string) => void;
  onClickWorkspaceInfo: (workspace: string) => void;
}

const filter: OptionsFilter = ({ options, search }) => {
  const filtered = (options as ComboboxItem[]).filter((option) =>
    option.label.toLowerCase().trim().includes(search.toLowerCase().trim())
  );

  filtered.sort((a, b) => a.label.localeCompare(b.label));
  return options;
};

export function WorkspaceSelector({
  workspaces,
  selectedWorkspace,
  onSelectWorkspace,
  onClickWorkspaceInfo,
}: WorkspaceSelectorProps) {
  return (
    <>
      <Text size="xs" fw={700} c="dimmed" mb="sm">
        CURRENT WORKSPACE
      </Text>
      <Group gap="xs">
        <Select
          value={selectedWorkspace}
          onChange={(value) => {
            if (value && workspaces[value]) {
              onSelectWorkspace(value);
            }
          }}
          data={Object.entries(workspaces).map(([name, workspace]) => ({
            value: name,
            label: workspace.displayName || name,
          }))}
          filter={filter}
          placeholder="No workspace selected"
          searchable
          size="md"
          styles={{
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
        <ActionIcon
          size="xl"
          aria-label="View workspace documentation"
          onClick={() => onClickWorkspaceInfo(selectedWorkspace)}
          variant="transparent"
          color="var(--mantine-color-dark-4)"
        >
          <IconInfoCircle />
        </ActionIcon>
      </Group>
    </>
  );
}
