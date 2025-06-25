import { Badge, Card, Code, Group, Stack, Text, Title } from "@mantine/core";
import { IconExternalLink } from "@tabler/icons-react";
import { EnrichedExecutable } from "../../../types/executable";

export type ExecutableLaunchDetailsProps = {
  executable: EnrichedExecutable;
};

export function ExecutableLaunchDetails({
  executable,
}: ExecutableLaunchDetailsProps) {
  return (
    <Card withBorder>
      <Stack gap="sm">
        <Title order={4}>
          <Group gap="xs">
            <IconExternalLink size={16} />
            Launch Configuration
          </Group>
        </Title>
        <Stack gap="xs">
          {executable.launch?.app && (
            <div>
              <Title order={5}>App:</Title>
              <Code>{executable.launch.app}</Code>
            </div>
          )}
          {executable.launch?.uri && (
            <div>
              <Title order={5}>URI:</Title>
              <Text
                component="a"
                href={executable.launch.uri}
                target="_blank"
                rel="noopener noreferrer"
              >
                {executable.launch.uri}
              </Text>
            </div>
          )}
          {executable.launch?.wait && (
            <div>
              <Title order={5}>Wait:</Title>
              <Badge variant="light" color="blue">
                enabled
              </Badge>
            </div>
          )}
        </Stack>
      </Stack>
    </Card>
  );
}
