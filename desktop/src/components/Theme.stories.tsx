import {
  ActionIcon,
  Box,
  Button,
  Card,
  Group,
  Paper,
  Stack,
  Text,
  Title,
} from "@mantine/core";
import type { Meta, StoryObj } from "@storybook/react";
import { mantineTheme as theme } from "../theme/mantineTheme";

const meta = {
  title: "Design System/Theme",
  parameters: {
    layout: "padded",
  },
} satisfies Meta;

export default meta;
type Story = StoryObj;

const ColorSwatch = ({ color, name }: { color: string; name: string }) => (
  <Box>
    <Paper
      style={{
        width: 100,
        height: 100,
        backgroundColor: color,
        borderRadius: theme.defaultRadius,
      }}
    />
    <Text size="sm" mt="xs">
      {name}
    </Text>
  </Box>
);

export const Colors: Story = {
  render: () => (
    <Stack gap="xl">
      <Box>
        <Title order={3} mb="md">
          Palette Colors
        </Title>
        <Group>
          {theme.colors?.palette?.map((color, index) => (
            <ColorSwatch key={index} color={color} name={`Palette ${index}`} />
          ))}
        </Group>
      </Box>

      <Box>
        <Title order={3} mb="md">
          Dark Colors
        </Title>
        <Group>
          {theme.colors?.dark?.map((color, index) => (
            <ColorSwatch key={index} color={color} name={`Dark ${index}`} />
          ))}
        </Group>
      </Box>

      <Box>
        <Title order={3} mb="md">
          Accent Colors
        </Title>
        <Group>
          {theme.colors?.accent?.map((color, index) => (
            <ColorSwatch key={index} color={color} name={`Accent ${index}`} />
          ))}
        </Group>
      </Box>
    </Stack>
  ),
};

export const Typography: Story = {
  render: () => (
    <Stack gap="md">
      <Title order={1}>Heading 1</Title>
      <Title order={2}>Heading 2</Title>
      <Title order={3}>Heading 3</Title>
      <Title order={4}>Heading 4</Title>
      <Title order={5}>Heading 5</Title>
      <Title order={6}>Heading 6</Title>
      <Text size="xl">Extra Large Text</Text>
      <Text size="lg">Large Text</Text>
      <Text size="md">Medium Text</Text>
      <Text size="sm">Small Text</Text>
      <Text size="xs">Extra Small Text</Text>
    </Stack>
  ),
};

export const Components: Story = {
  render: () => (
    <Stack gap="xl">
      <Box>
        <Title order={3} mb="md">
          Buttons
        </Title>
        <Group>
          <Button>Default Button</Button>
          <Button variant="light">Light Button</Button>
          <Button variant="outline">Outline Button</Button>
          <Button variant="subtle">Subtle Button</Button>
        </Group>
      </Box>

      <Box>
        <Title order={3} mb="md">
          Action Icons
        </Title>
        <Group>
          <ActionIcon size="lg" variant="light">
            <span>🔍</span>
          </ActionIcon>
          <ActionIcon size="lg" variant="light">
            <span>⚙️</span>
          </ActionIcon>
          <ActionIcon size="lg" variant="light">
            <span>📝</span>
          </ActionIcon>
        </Group>
      </Box>

      <Box>
        <Title order={3} mb="md">
          Cards
        </Title>
        <Card style={{ maxWidth: 300 }}>
          <Text>
            This is a card component with the default radius of{" "}
            {theme.defaultRadius}
          </Text>
        </Card>
      </Box>
    </Stack>
  ),
};
