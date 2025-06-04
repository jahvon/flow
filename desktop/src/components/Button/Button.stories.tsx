import type { Meta, StoryObj } from '@storybook/react';
import { Button } from './Button';

const meta: Meta<typeof Button> = {
  title: 'Components/Button',
  component: Button,
  parameters: {
    layout: 'centered',
  },
  tags: ['autodocs'],
};

export default meta;
type Story = StoryObj<typeof Button>;

export const Primary: Story = {
  args: {
    children: 'Button',
    variant: 'filled',
  },
};

export const Secondary: Story = {
  args: {
    children: 'Button',
    variant: 'light',
  },
};

export const Danger: Story = {
  args: {
    children: 'Button',
    variant: 'filled',
    color: 'red',
  },
};