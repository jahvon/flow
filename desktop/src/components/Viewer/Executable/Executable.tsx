import {
  ActionIcon,
  Badge,
  Button,
  Card,
  Code,
  Divider,
  Drawer,
  Grid,
  Group,
  ScrollArea,
  Stack,
  Table,
  Text,
  ThemeIcon,
  Title,
  Tooltip,
} from "@mantine/core";
import {
  IconClock,
  IconExternalLink,
  IconEye,
  IconFile,
  IconFlag,
  IconInfoCircle,
  IconKey,
  IconPlayerPlay,
  IconTag,
  IconTerminal,
} from "@tabler/icons-react";
import { invoke } from "@tauri-apps/api/core";
import { listen } from "@tauri-apps/api/event";
import { openPath } from "@tauri-apps/plugin-opener";
import { useEffect, useRef, useState } from "react";
import { useNotifier } from "../../../hooks/useNotifier";
import { useSettings } from "../../../hooks/useSettings";
import { EnrichedExecutable } from "../../../types/executable";
import {
  ExecutableArgument,
  ExecutableParameter,
} from "../../../types/generated/flowfile";
import { NotificationType } from "../../../types/notification";
import { CodeHighlighter } from "../../CodeHighlighter";
import { MarkdownRenderer } from "../../MarkdownRenderer";

export type ExecutableInfoProps = {
  executable: EnrichedExecutable;
};

function getExecutableTypeInfo(executable: EnrichedExecutable) {
  if (executable.exec)
    return {
      type: "exec",
      icon: IconTerminal,
      description: "Command execution",
    };
  if (executable.serial)
    return {
      type: "serial",
      icon: IconTerminal,
      description: "Sequential execution",
    };
  if (executable.parallel)
    return {
      type: "parallel",
      icon: IconTerminal,
      description: "Parallel execution",
    };
  if (executable.launch)
    return {
      type: "launch",
      icon: IconExternalLink,
      description: "Launch application/URI",
    };
  if (executable.request)
    return {
      type: "request",
      icon: IconExternalLink,
      description: "HTTP request",
    };
  if (executable.render)
    return {
      type: "render",
      icon: IconTerminal,
      description: "Render template",
    };
  return { type: "unknown", icon: IconTerminal, description: "Unknown type" };
}

function getVisibilityColor(visibility?: string) {
  switch (visibility) {
    case "public":
      return "green";
    case "private":
      return "blue";
    case "internal":
      return "orange";
    case "hidden":
      return "red";
    default:
      return "gray";
  }
}

type OutputLine = {
  type: "stdout" | "stderr";
  line: string;
};

export default function ExecutableInfo({ executable }: ExecutableInfoProps) {
  const typeInfo = getExecutableTypeInfo(executable);
  const { settings } = useSettings();
  const { setNotification } = useNotifier();
  const [output, setOutput] = useState<OutputLine[]>([]);
  const outputListenerSetup = useRef(false);

  useEffect(() => {
    if (outputListenerSetup.current) {
      return;
    }

    outputListenerSetup.current = true;
    let unlistenOutput: (() => void) | undefined;
    let unlistenComplete: (() => void) | undefined;

    const setupListeners = async () => {
      console.log("Setting up listeners for executable:", executable.ref);
      unlistenOutput = await listen("command-output", (event) => {
        const payload = event.payload as OutputLine;
        console.log("Received output:", payload);
        setOutput((prev) => [...prev, payload]);
      });

      unlistenComplete = await listen("command-complete", (event) => {
        const payload = event.payload as {
          success: boolean;
          exit_code: number | null;
        };
        setNotification({
          title: payload.success ? "Execution completed" : "Execution failed",
          message: payload.success
            ? "Execution completed successfully"
            : "Execution failed",
          type: payload.success
            ? NotificationType.Success
            : NotificationType.Error,
        });
      });
    };

    setupListeners();

    return () => {
      outputListenerSetup.current = false;
      if (unlistenOutput) {
        unlistenOutput();
      }
      if (unlistenComplete) {
        unlistenComplete();
      }
    };
  }, [setNotification]);

  const onOpenFile = async () => {
    try {
      console.log(executable.flowfile);
      await openPath(executable.flowfile, settings.executableApp || undefined);
    } catch (error) {
      console.error(error);
    }
  };

  const onExecute = async () => {
    try {
      setOutput([]);

      setNotification({
        title: "Execution started",
        message: `Execution of ${executable.ref} started`,
        type: NotificationType.Success,
        autoClose: true,
        autoCloseDelay: 5000,
      });
      await invoke("execute", {
        verb: executable.verb,
        executableId: executable.id,
        args: [],
      });
    } catch (error) {
      console.error(error);
      setNotification({
        title: "Execution failed",
        message: `Execution of ${executable.ref} failed`,
        type: NotificationType.Error,
        autoClose: true,
        autoCloseDelay: 5000,
      });
    }
  };

  return (
    <Stack gap="lg">
      <Card withBorder>
        <Stack gap="md">
          <Group justify="space-between" align="flex-start">
            <Stack gap="xs">
              <Group gap="sm" align="center">
                <ThemeIcon variant="light" size="lg">
                  <typeInfo.icon size={20} />
                </ThemeIcon>
                <div>
                  <Title order={2}>{executable.ref}</Title>
                  <Text size="sm" c="dimmed">
                    {typeInfo.description}
                  </Text>
                </div>
              </Group>

              <Group gap="xs">
                <Tooltip label={`Defined in ${executable.flowfile}`}>
                  <Badge variant="light" size="sm">
                    <Group gap={4}>
                      <IconFile size={12} />
                      {executable.flowfile.split("/").pop() ||
                        executable.flowfile}
                    </Group>
                  </Badge>
                </Tooltip>
                <Badge
                  variant="light"
                  color={getVisibilityColor(executable.visibility)}
                >
                  <Group gap={4}>
                    <IconEye size={12} />
                    {executable.visibility || "public"}
                  </Group>
                </Badge>
                {executable.timeout && (
                  <Badge variant="light" color="gray">
                    <Group gap={4}>
                      <IconClock size={12} />
                      {executable.timeout}
                    </Group>
                  </Badge>
                )}
              </Group>
            </Stack>

            <Group gap="sm">
              <Button
                leftSection={<IconPlayerPlay size={16} />}
                onClick={onExecute}
                size="md"
              >
                Execute
              </Button>
              <Button
                variant="light"
                leftSection={<IconExternalLink size={16} />}
                onClick={onOpenFile}
                size="md"
              >
                Edit
              </Button>
            </Group>
          </Group>

          {executable.description && (
            <>
              <Divider />
              <MarkdownRenderer>{executable.description}</MarkdownRenderer>
            </>
          )}
        </Stack>
      </Card>

      <Grid>
        {executable.aliases && executable.aliases.length > 0 && (
          <Grid.Col span={6}>
            <Card withBorder>
              <Stack gap="sm">
                <Title order={4}>
                  <Group gap="xs">
                    <IconTag size={16} />
                    Aliases
                  </Group>
                </Title>
                <Group gap="xs">
                  {executable.aliases.map((alias, index) => (
                    <Code key={index}>{alias}</Code>
                  ))}
                </Group>
              </Stack>
            </Card>
          </Grid.Col>
        )}

        {executable.tags && executable.tags.length > 0 && (
          <Grid.Col span={6}>
            <Card withBorder>
              <Stack gap="sm">
                <Title order={4}>
                  <Group gap="xs">
                    <IconTag size={16} />
                    Tags
                  </Group>
                </Title>
                <Group gap="xs">
                  {executable.tags.map((tag, index) => (
                    <Badge key={index} variant="dot">
                      {tag}
                    </Badge>
                  ))}
                </Group>
              </Stack>
            </Card>
          </Grid.Col>
        )}

        {executable.exec?.params && executable.exec.params.length > 0 && (
          <Grid.Col span={12}>
            <Card withBorder>
              <Stack gap="sm">
                <Title order={4}>
                  <Group gap="xs">
                    <IconKey size={16} />
                    Environment Parameters
                  </Group>
                </Title>
                <Table>
                  <Table.Thead>
                    <Table.Tr>
                      <Table.Th>Variable</Table.Th>
                      <Table.Th>Type</Table.Th>
                      <Table.Th>Source</Table.Th>
                    </Table.Tr>
                  </Table.Thead>
                  <Table.Tbody>
                    {executable.exec.params.map(
                      (param: ExecutableParameter, index: number) => {
                        const type = param.text
                          ? "static"
                          : param.secretRef
                            ? "secret"
                            : "prompt";
                        const source =
                          param.text || param.secretRef || param.prompt;

                        return (
                          <Table.Tr key={index}>
                            <Table.Td>
                              <Code>{param.envKey}</Code>
                            </Table.Td>
                            <Table.Td>
                              <Badge
                                size="sm"
                                variant="light"
                                color={
                                  type === "secret"
                                    ? "red"
                                    : type === "prompt"
                                      ? "blue"
                                      : "gray"
                                }
                              >
                                {type}
                              </Badge>
                            </Table.Td>
                            <Table.Td>
                              <Text size="sm" style={{ maxWidth: 200 }}>
                                {source}
                              </Text>
                            </Table.Td>
                          </Table.Tr>
                        );
                      }
                    )}
                  </Table.Tbody>
                </Table>
              </Stack>
            </Card>
          </Grid.Col>
        )}

        {executable.exec?.args && executable.exec.args.length > 0 && (
          <Grid.Col span={12}>
            <Card withBorder>
              <Stack gap="sm">
                <Title order={4}>
                  <Group gap="xs">
                    <IconFlag size={16} />
                    Command Arguments
                  </Group>
                </Title>
                <Table>
                  <Table.Thead>
                    <Table.Tr>
                      <Table.Th>Variable</Table.Th>
                      <Table.Th>Input</Table.Th>
                      <Table.Th>Type</Table.Th>
                      <Table.Th>Required</Table.Th>
                    </Table.Tr>
                  </Table.Thead>
                  <Table.Tbody>
                    {executable.exec.args.map(
                      (arg: ExecutableArgument, index: number) => (
                        <Table.Tr key={index}>
                          <Table.Td>
                            <Code>{arg.envKey}</Code>
                          </Table.Td>
                          <Table.Td>
                            <Group gap="xs">
                              <Badge size="sm" variant="light">
                                {arg.pos ? `pos-${arg.pos}` : `--${arg.flag}`}
                              </Badge>
                              {arg.default && (
                                <Tooltip label={`Default: ${arg.default}`}>
                                  <ActionIcon variant="subtle" size="xs">
                                    <IconInfoCircle size={12} />
                                  </ActionIcon>
                                </Tooltip>
                              )}
                            </Group>
                          </Table.Td>
                          <Table.Td>
                            <Text size="sm">{arg.type || "string"}</Text>
                          </Table.Td>
                          <Table.Td>
                            <Badge
                              size="sm"
                              variant="light"
                              color={arg.required ? "red" : "green"}
                            >
                              {arg.required ? "Yes" : "No"}
                            </Badge>
                          </Table.Td>
                        </Table.Tr>
                      )
                    )}
                  </Table.Tbody>
                </Table>
              </Stack>
            </Card>
          </Grid.Col>
        )}

        {executable.exec && (
          <Grid.Col span={12}>
            <Card withBorder>
              <Stack gap="sm">
                <Title order={4}>
                  <Group gap="xs">
                    <IconTerminal size={16} />
                    Execution Details
                  </Group>
                </Title>
                <Stack gap="xs">
                  {executable.exec.cmd && (
                    <div>
                      <Title order={5}>Command:</Title>
                      <CodeHighlighter
                        theme={settings.theme}
                        copyButton={false}
                      >
                        {executable.exec.cmd}
                      </CodeHighlighter>
                    </div>
                  )}
                  {executable.exec.file && (
                    <div>
                      <Title order={5}>File:</Title>
                      <Code>{executable.exec.file}</Code>
                    </div>
                  )}
                  {executable.exec.dir && (
                    <div>
                      <Title order={5}>Directory:</Title>
                      <Code>{executable.exec.dir}</Code>
                    </div>
                  )}
                  {executable.exec.logMode && (
                    <div>
                      <Title order={5}>Log Mode:</Title>
                      <Badge variant="light">{executable.exec.logMode}</Badge>
                    </div>
                  )}
                </Stack>
              </Stack>
            </Card>
          </Grid.Col>
        )}
      </Grid>

      {output.length > 0 && (
        <Drawer
          opened={true}
          onClose={() => setOutput([])}
          title="Execution Output"
          scrollAreaComponent={ScrollArea.Autosize}
          position="bottom"
        >
          <Stack gap="sm">
            {output.map((line, index) => (
              <Text
                key={index}
                color={line.type === "stderr" ? "red" : "green"}
              >
                {line.line}
              </Text>
            ))}
          </Stack>
        </Drawer>
      )}
    </Stack>
  );
}
