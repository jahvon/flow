import type { Meta, StoryObj } from "@storybook/react";
import { MarkdownRenderer } from "./MarkdownRenderer";

const meta = {
  title: "Components/MarkdownRenderer",
  component: MarkdownRenderer,
  parameters: {
    layout: "padded",
  },
  tags: ["autodocs"],
} satisfies Meta<typeof MarkdownRenderer>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    children: `# Heading 1
## Heading 2
### Heading 3
#### Heading 4
##### Heading 5
###### Heading 6

This is a paragraph with **bold text** and *italic text*.

- List item 1
- List item 2
- List item 3

1. Numbered item 1
2. Numbered item 2
3. Numbered item 3

\`\`\`javascript
// Code block
function hello() {
  console.log("Hello, world!");
}
\`\`\`

[Link to example](https://example.com)`,
  },
};

export const SimpleText: Story = {
  args: {
    children: `This is a simple paragraph with some **bold text** and *italic text*.

You can also include [links](https://example.com) and \`inline code\`.`,
  },
};

export const HeadingsOnly: Story = {
  args: {
    children: `# Main Title
## Section Title
### Subsection
#### Sub-subsection
##### Small heading
###### Tiny heading`,
  },
};
