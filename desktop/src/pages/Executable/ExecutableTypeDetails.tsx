import { Grid } from "@mantine/core";
import { EnrichedExecutable } from "../../types/executable";
import { ExecutableExecDetails } from "./types/ExecutableExecDetails";
import { ExecutableLaunchDetails } from "./types/ExecutableLaunchDetails";
import { ExecutableParallelDetails } from "./types/ExecutableParallelDetails";
import { ExecutableRenderDetails } from "./types/ExecutableRenderDetails";
import { ExecutableRequestDetails } from "./types/ExecutableRequestDetails";
import { ExecutableSerialDetails } from "./types/ExecutableSerialDetails";

export type ExecutableTypeDetailsProps = {
  executable: EnrichedExecutable;
};

export function ExecutableTypeDetails({
  executable,
}: ExecutableTypeDetailsProps) {
  return (
    <Grid>
      {executable.exec && (
        <Grid.Col span={12}>
          <ExecutableExecDetails executable={executable} />
        </Grid.Col>
      )}

      {executable.launch && (
        <Grid.Col span={12}>
          <ExecutableLaunchDetails executable={executable} />
        </Grid.Col>
      )}

      {executable.request && (
        <Grid.Col span={12}>
          <ExecutableRequestDetails executable={executable} />
        </Grid.Col>
      )}

      {executable.render && (
        <Grid.Col span={12}>
          <ExecutableRenderDetails executable={executable} />
        </Grid.Col>
      )}

      {executable.serial && (
        <Grid.Col span={12}>
          <ExecutableSerialDetails executable={executable} />
        </Grid.Col>
      )}

      {executable.parallel && (
        <Grid.Col span={12}>
          <ExecutableParallelDetails executable={executable} />
        </Grid.Col>
      )}
    </Grid>
  );
}
