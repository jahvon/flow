import { Badge, Card, Code, Group, Stack, Text, Title } from "@mantine/core";
import { IconTerminal } from "@tabler/icons-react";
import { useSettings } from "../../../hooks/useSettings";
import { EnrichedExecutable } from "../../../types/executable";
import { CodeHighlighter } from "../../../components/CodeHighlighter";

export type ExecutableParallelDetailsProps = {
  executable: EnrichedExecutable;
};

export function ExecutableParallelDetails({
  executable,
}: ExecutableParallelDetailsProps) {
  const { settings } = useSettings();

  return (
    <Card withBorder>
      <Stack gap="sm">
        <Title order={4}>
          <Group gap="xs">
            <IconTerminal size={16} />
            Parallel Configuration
          </Group>
        </Title>
        <Stack gap="xs">
          {executable.parallel?.maxThreads &&
            executable.parallel.maxThreads > 0 && (
              <div>
                <Title order={5}>Max Threads:</Title>
                <Code>{executable.parallel.maxThreads}</Code>
              </div>
            )}
          {executable.parallel?.failFast !== undefined && (
            <div>
              <Title order={5}>Fail Fast:</Title>
              <Badge
                variant="light"
                color={executable.parallel.failFast ? "red" : "green"}
              >
                {executable.parallel.failFast ? "enabled" : "disabled"}
              </Badge>
            </div>
          )}
          {executable.parallel?.execs &&
            executable.parallel.execs.length > 0 && (
              <div>
                <Title order={5}>Executables:</Title>
                <Stack gap="md">
                  {executable.parallel.execs.map((exec, index) => (
                    <div key={index}>
                      <Text fw={500}>
                        {index + 1}. {exec.ref ? `ref: ${exec.ref}` : "cmd:"}
                      </Text>
                      {exec.cmd && (
                        <CodeHighlighter
                          theme={settings.theme}
                          copyButton={false}
                        >
                          {exec.cmd}
                        </CodeHighlighter>
                      )}
                      {exec.retries !== undefined && exec.retries > 0 && (
                        <div>
                          <Text size="sm" c="dimmed">
                            • Retries: {exec.retries}
                          </Text>
                        </div>
                      )}
                      {exec.args && exec.args.length > 0 && (
                        <div>
                          <Text size="sm" c="dimmed">
                            • Arguments:
                          </Text>
                          <Stack gap="xs" ml="md">
                            {exec.args.map((arg, argIndex) => (
                              <Text key={argIndex} size="sm" c="dimmed">
                                - {arg}
                              </Text>
                            ))}
                          </Stack>
                        </div>
                      )}
                    </div>
                  ))}
                </Stack>
              </div>
            )}
        </Stack>
      </Stack>
    </Card>
  );
}
