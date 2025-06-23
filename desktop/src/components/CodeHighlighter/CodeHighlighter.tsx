import { Box, Code } from "@mantine/core";
import Prism from "prismjs";
import "prismjs/components/prism-bash";
import "prismjs/components/prism-shell-session";
import "prismjs/themes/prism-dark.css";
import "prismjs/themes/prism.css";
import { useEffect, useRef } from "react";

import { ThemeName } from "../../theme/types";
import { getConfigForTheme } from "./config";

interface CodeHighlighterProps {
  children: string;
  className?: string;
  copyButton?: boolean;
  theme?: ThemeName;
}

export function CodeHighlighter({
  children,
  className,
  copyButton,
  theme,
}: CodeHighlighterProps) {
  const codeRef = useRef<HTMLElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);

  const config = getConfigForTheme(theme);
  const finalCopyButton =
    copyButton !== undefined ? copyButton : config.defaultCopyButton;
  const language = "bash";

  useEffect(() => {
    if (codeRef.current) {
      Prism.highlightElement(codeRef.current);
    }
  }, [children]);

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(children);
      // TODO: Add a toast notification here
      console.log("Code copied to clipboard");
    } catch (error) {
      console.error("Failed to copy code:", error);
    }
  };

  return (
    <Box className={className}>
      <Box
        ref={containerRef}
        pos="relative"
        style={{
          borderRadius: config.styling.borderRadius,
          overflow: "hidden",
        }}
      >
        {finalCopyButton && (
          <Box pos="absolute" top={8} right={8} style={{ zIndex: 10 }}>
            <Code
              component="button"
              onClick={handleCopy}
              style={{
                cursor: "pointer",
                ...config.styling.copyButtonStyle,
                color: "inherit",
                fontFamily: "inherit",
              }}
            >
              Copy
            </Code>
          </Box>
        )}
        <pre
          style={{
            margin: 0,
            padding: config.styling.padding,
            background: config.styling.backgroundColor,
            borderRadius: config.styling.borderRadius,
            overflow: "auto",
            fontSize: config.styling.fontSize,
            lineHeight: config.styling.lineHeight,
          }}
        >
          <code
            ref={codeRef}
            className={`language-${language}`}
            style={{
              fontFamily: "var(--mantine-font-family-monospace)",
            }}
          >
            {children}
          </code>
        </pre>
      </Box>
    </Box>
  );
}
