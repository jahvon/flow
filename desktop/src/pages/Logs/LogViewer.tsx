import { Stack, Text } from "@mantine/core";
import { useMemo } from "react";

export interface LogLine {
  line: string;
  type?: "stdout" | "stderr";
}

export interface LogViewerProps {
  logs: string[] | LogLine[];
  formatted?: boolean;
  fontSize?: number;
}

interface ParsedLogLine {
  time?: string;
  level?: string;
  message: string;
  raw: string;
}

const LOG_LEVEL_COLORS = {
  INF: "blue",
  INFO: "blue",
  NTC: "yellow",
  NOTICE: "yellow",
  ERR: "red",
  ERROR: "red",
  DBG: "gray",
  DEBUG: "gray",
} as const;

function parseLogLine(line: string): ParsedLogLine {
  // Pattern: TIME LEVEL MESSAGE (e.g., "5:49PM INF building go binary...")
  const logPattern = /^(\d{1,2}:\d{2}[AP]M)\s+([A-Z]{3,6})\s+(.+)$/;
  const match = line.match(logPattern);

  if (match) {
    const [, time, level, message] = match;
    return { time, level, message, raw: line };
  }

  return { message: line, raw: line };
}

function LogLineComponent({
  parsedLine,
  formatted,
  fontSize,
}: {
  parsedLine: ParsedLogLine;
  formatted: boolean;
  fontSize: number;
}) {
  if (!formatted || !parsedLine.time || !parsedLine.level) {
    return (
      <Text
        component="div"
        size={`${fontSize}px`}
        ff="monospace"
        style={{ whiteSpace: "pre-wrap" }}
      >
        {parsedLine.raw}
      </Text>
    );
  }

  const levelColor =
    LOG_LEVEL_COLORS[parsedLine.level as keyof typeof LOG_LEVEL_COLORS] ||
    "gray";

  return (
    <Text
      component="div"
      size={`${fontSize}px`}
      ff="monospace"
      style={{ whiteSpace: "pre-wrap" }}
    >
      <Text component="span" c="dimmed">
        {parsedLine.time}
      </Text>{" "}
      <Text component="span" c={levelColor} fw={500}>
        {parsedLine.level}
      </Text>{" "}
      <Text component="span">{parsedLine.message}</Text>
    </Text>
  );
}

export function LogViewer({
  logs,
  formatted = false,
  fontSize = 13,
}: LogViewerProps) {
  const parsedLogs = useMemo(() => {
    return logs.map((log) => {
      const lineText = typeof log === "string" ? log : log.line;
      return parseLogLine(lineText);
    });
  }, [logs]);

  return (
    <Stack gap={0}>
      {parsedLogs.map((parsedLine, index) => (
        <LogLineComponent
          key={index}
          parsedLine={parsedLine}
          formatted={formatted}
          fontSize={fontSize}
        />
      ))}
    </Stack>
  );
}
