import { Card, Code, Group, Stack, Title } from "@mantine/core";
import { IconTerminal } from "@tabler/icons-react";
import { EnrichedExecutable } from "../../../types/executable";

export type ExecutableRenderDetailsProps = {
  executable: EnrichedExecutable;
};

export function ExecutableRenderDetails({
  executable,
}: ExecutableRenderDetailsProps) {
  return (
    <Card withBorder>
      <Stack gap="sm">
        <Title order={4}>
          <Group gap="xs">
            <IconTerminal size={16} />
            Render Configuration
          </Group>
        </Title>
        <Stack gap="xs">
          {executable.render?.dir && (
            <div>
              <Title order={5}>Executed from:</Title>
              <Code>{executable.render.dir}</Code>
            </div>
          )}
          {executable.render?.templateFile && (
            <div>
              <Title order={5}>Template File:</Title>
              <Code>{executable.render.templateFile}</Code>
            </div>
          )}
          {executable.render?.templateDataFile && (
            <div>
              <Title order={5}>Template Store File:</Title>
              <Code>{executable.render.templateDataFile}</Code>
            </div>
          )}
        </Stack>
      </Stack>
    </Card>
  );
}
