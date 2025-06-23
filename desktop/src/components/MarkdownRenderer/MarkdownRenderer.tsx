import { Text, Title } from "@mantine/core";
import { ComponentPropsWithoutRef } from "react";
import ReactMarkdown from "react-markdown";

interface MarkdownRendererProps {
  children: string;
  className?: string;
}

export function MarkdownRenderer({
  children,
  className,
}: MarkdownRendererProps) {
  return (
    <div className={className}>
      <ReactMarkdown
        components={{
          h1: (props: ComponentPropsWithoutRef<typeof Title>) => (
            <Title order={3} {...props} />
          ),
          h2: (props: ComponentPropsWithoutRef<typeof Title>) => (
            <Title order={4} {...props} />
          ),
          h3: (props: ComponentPropsWithoutRef<typeof Title>) => (
            <Title order={5} {...props} />
          ),
          h4: (props: ComponentPropsWithoutRef<typeof Title>) => (
            <Title order={6} {...props} />
          ),
          h5: (props: ComponentPropsWithoutRef<typeof Text>) => (
            <Text size="sm" fw={700} {...props} />
          ),
          h6: (props: ComponentPropsWithoutRef<typeof Text>) => (
            <Text size="xs" fw={500} {...props} />
          ),
          p: (props: ComponentPropsWithoutRef<typeof Text>) => (
            <Text size="sm" {...props} />
          ),
        }}
      >
        {children}
      </ReactMarkdown>
    </div>
  );
}
