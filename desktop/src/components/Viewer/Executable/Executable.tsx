import {
  ActionIcon,
  Alert,
  Badge,
  Button,
  Card,
  Code,
  Divider,
  Grid,
  Group,
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
  IconFlag,
  IconInfoCircle,
  IconKey,
  IconPlayerPlay,
  IconTag,
  IconTerminal,
} from "@tabler/icons-react";
import { EnrichedExecutable } from "../../../types/executable";
import {
  ExecutableArgument,
  ExecutableParameter,
} from "../../../types/generated/flowfile";

export type ExecutableInfoProps = {
  executable: EnrichedExecutable;
  onExecute?: () => void;
  onOpenFile?: () => void;
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

export default function ExecutableInfo({
  executable,
  onExecute,
  onOpenFile,
}: ExecutableInfoProps) {
  const typeInfo = getExecutableTypeInfo(executable);

  return (
    <Stack gap="lg">
      {/* Header with Action Buttons */}
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
                <Badge
                  variant="light"
                  color={getVisibilityColor(executable.visibility)}
                >
                  <Group gap={4}>
                    <IconEye size={12} />
                    {executable.visibility || "public"}
                  </Group>
                </Badge>
                <Badge variant="light" color="gray">
                  <Group gap={4}>
                    <IconClock size={12} />
                    {executable.timeout || "30m"}
                  </Group>
                </Badge>
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
                Open File
              </Button>
            </Group>
          </Group>

          {executable.description && (
            <>
              <Divider />
              <Text>{executable.description}</Text>
            </>
          )}
        </Stack>
      </Card>

      <Grid>
        {/* Left Column */}
        <Grid.Col span={6}>
          <Stack gap="md">
            {/* Aliases */}
            {executable.aliases && executable.aliases.length > 0 && (
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
            )}

            {/* Tags */}
            {executable.tags && executable.tags.length > 0 && (
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
            )}

            {/* Execution Details */}
            {executable.exec && (
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
                        <Text size="sm" fw={500}>
                          Command:
                        </Text>
                        <Code block>{executable.exec.cmd}</Code>
                      </div>
                    )}
                    {executable.exec.file && (
                      <div>
                        <Text size="sm" fw={500}>
                          File:
                        </Text>
                        <Code>{executable.exec.file}</Code>
                      </div>
                    )}
                    {executable.exec.dir && (
                      <div>
                        <Text size="sm" fw={500}>
                          Directory:
                        </Text>
                        <Code>{executable.exec.dir}</Code>
                      </div>
                    )}
                    {executable.exec.logMode && (
                      <div>
                        <Text size="sm" fw={500}>
                          Log Mode:
                        </Text>
                        <Badge variant="light">{executable.exec.logMode}</Badge>
                      </div>
                    )}
                  </Stack>
                </Stack>
              </Card>
            )}
          </Stack>
        </Grid.Col>

        {/* Right Column */}
        <Grid.Col span={6}>
          <Stack gap="md">
            {/* Parameters */}
            {executable.exec?.params && executable.exec.params.length > 0 && (
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
            )}

            {/* Arguments */}
            {executable.exec?.args && executable.exec.args.length > 0 && (
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
            )}
          </Stack>
        </Grid.Col>
      </Grid>

      {/* Footer */}
      <Alert icon={<IconInfoCircle size={16} />} color="blue" variant="light">
        <Group justify="space-between">
          <Text size="sm">
            Defined in <Code>{executable.flowfile}</Code>
          </Text>
          <Button variant="subtle" size="xs" onClick={onOpenFile}>
            View Source
          </Button>
        </Group>
      </Alert>
    </Stack>
  );
}
