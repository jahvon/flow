import React from 'react';
import { MantineProvider } from '@mantine/core';
import { theme } from '../src/theme';

export const withMantine = (Story: React.ComponentType) => (
  <MantineProvider theme={theme}>
    <Story />
  </MantineProvider>
);