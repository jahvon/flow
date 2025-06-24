import { Badge, Card, Code, Group, Stack, Title } from "@mantine/core";
import { IconTerminal } from "@tabler/icons-react";
import { useSettings } from "../../../../hooks/useSettings";
import { EnrichedExecutable } from "../../../../types/executable";
import { CodeHighlighter } from "../../../CodeHighlighter";

export type ExecutableExecDetailsProps = {
  executable: EnrichedExecutable;
};

export function ExecutableExecDetails({
  executable,
}: ExecutableExecDetailsProps) {
  const { settings } = useSettings();

  return (
    <Card withBorder>
      <Stack gap="sm">
        <Title order={4}>
          <Group gap="xs">
            <IconTerminal size={16} />
            Execution Details
          </Group>
        </Title>
        <Stack gap="xs">
          {executable.exec?.cmd && (
            <div>
              <Title order={5}>Command:</Title>
              <CodeHighlighter theme={settings.theme} copyButton={false}>
                {executable.exec.cmd}
              </CodeHighlighter>
            </div>
          )}
          {executable.exec?.file && (
            <div>
              <Title order={5}>File:</Title>
              <Code>{executable.exec.file}</Code>
            </div>
          )}
          {executable.exec?.dir && (
            <div>
              <Title order={5}>Directory:</Title>
              <Code>{executable.exec.dir}</Code>
            </div>
          )}
          {executable.exec?.logMode && (
            <div>
              <Title order={5}>Log Mode:</Title>
              <Badge variant="light">{executable.exec.logMode}</Badge>
            </div>
          )}
        </Stack>
      </Stack>
    </Card>
  );
}
