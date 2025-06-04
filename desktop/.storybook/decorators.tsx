import React from 'react';
import { MantineProvider } from '@mantine/core';

export const withMantine = (Story: React.ComponentType) => (
  <MantineProvider>
    <Story />
  </MantineProvider>
);