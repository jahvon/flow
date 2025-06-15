import {
  ActionIcon,
  ComboboxItem,
  Group,
  OptionsFilter,
  Select,
} from "@mantine/core";
import { IconInfoCircle } from "@tabler/icons-react";
import type { Workspace } from "../../../types/generated/workspace";

interface WorkspaceSelectorProps {
  workspaces: Record<string, Workspace>;
  selectedWorkspace: string | null;
  onSelectWorkspace: (workspaceId: string) => void;
  onClickWorkspaceInfo: () => void;
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
      <Group gap="xs" align="center" wrap="nowrap">
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
        <ActionIcon
          size="md"
          aria-label="View workspace documentation"
          onClick={onClickWorkspaceInfo}
          variant="transparent"
          color="var(--mantine-color-dark-4)"
        >
          <IconInfoCircle />
        </ActionIcon>
      </Group>
    </>
  );
}
