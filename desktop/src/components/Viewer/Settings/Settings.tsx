import { Card, Select, Stack, Text, TextInput, Title } from "@mantine/core";
import { useSettings } from "../../../hooks/useSettings";
import { ThemeName } from "../../../theme/types";
import styles from "./Settings.module.css";

const themeOptions = [
  { value: "everforest", label: "Everforest" },
  { value: "dark", label: "Dark" },
  { value: "dracula", label: "Dracula" },
  { value: "light", label: "Light" },
  { value: "tokyo-night", label: "Tokyo Night" },
];

export function Settings() {
  const { settings, updateWorkspaceApp, updateExecutableApp, updateTheme } =
    useSettings();

  return (
    <div className={styles.settings}>
      <Title order={2} mb="md">
        Settings
      </Title>

      <Stack gap="sm">
        <Card className={styles.settings__section}>
          <Text size="lg" fw={500} mb="md">
            Appearance
          </Text>
          <Select
            label="Theme"
            description="Choose your preferred theme"
            value={settings.theme}
            onChange={(value) => value && updateTheme(value as ThemeName)}
            data={themeOptions}
          />
        </Card>

        <Card className={styles.settings__section}>
          <Text size="lg" fw={500} mb="md">
            Application Configuration
          </Text>
          <Stack gap="md">
            <TextInput
              label="Workspace Command"
              description="Command to use when opening workspace directories. Leave empty to use system default."
              value={settings.workspaceApp}
              onChange={(event) =>
                updateWorkspaceApp(event.currentTarget.value)
              }
            />
            <TextInput
              label="Executable Command"
              description="Command to use when opening flow files. Leave empty to use system default."
              value={settings.executableApp}
              onChange={(event) =>
                updateExecutableApp(event.currentTarget.value)
              }
            />
          </Stack>
        </Card>
      </Stack>
    </div>
  );
}
