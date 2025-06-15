import { Card, Select, Stack, Text, Title } from "@mantine/core";
import { useLocalStorage } from "@mantine/hooks";
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
  const [theme, setTheme] = useLocalStorage<ThemeName>({
    key: "theme",
    defaultValue: "everforest",
  });

  return (
    <div className={styles.settings}>
      <Title order={2} mb="xl">
        Settings
      </Title>

      <Stack gap="md">
        <Card className={styles.settings__section}>
          <Text size="lg" fw={500} mb="md">
            Appearance
          </Text>
          <Select
            label="Theme"
            description="Choose your preferred theme"
            value={theme}
            onChange={(value) => value && setTheme(value as ThemeName)}
            data={themeOptions}
          />
        </Card>
      </Stack>
    </div>
  );
}
