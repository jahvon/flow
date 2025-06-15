import { Title } from "@mantine/core";
import { EnrichedExecutable } from "../../../types/executable";

export type ExecutableInfoProps = {
  executable: EnrichedExecutable;
};

export default function ExecutableInfo({ executable }: ExecutableInfoProps) {
  return (
    <>
      <Title order={2}>{executable.ref}</Title>
    </>
  );
}
