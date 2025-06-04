import { Button as MantineButton } from '@mantine/core';
import { type ButtonProps } from '@mantine/core';

export function Button(props: ButtonProps) {
  return <MantineButton {...props} />;
}