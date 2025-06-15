import { MantineProvider, createTheme } from "@mantine/core";
import { useLocalStorage } from "@mantine/hooks";
import { themes } from "./themes";
import { ThemeName } from "./types";

interface ThemeProviderProps {
  children: React.ReactNode;
}

const createColorArray = (color: string) =>
  [
    color,
    color,
    color,
    color,
    color,
    color,
    color,
    color,
    color,
    color,
  ] as const;

export function ThemeProvider({ children }: ThemeProviderProps) {
  const [themeName, setThemeName] = useLocalStorage<ThemeName>({
    key: "theme",
    defaultValue: "everforest",
  });

  const currentTheme = themes[themeName];

  const theme = createTheme({
    colors: {
      primary: createColorArray(currentTheme.colors.primary),
      secondary: createColorArray(currentTheme.colors.secondary),
      tertiary: createColorArray(currentTheme.colors.tertiary),
      info: createColorArray(currentTheme.colors.info),
      body: createColorArray(currentTheme.colors.body),
      border: createColorArray(currentTheme.colors.border),
      emphasis: createColorArray(currentTheme.colors.emphasis),
    },
    primaryColor: "primary",
    primaryShade: 0,
    defaultRadius: "md",
    fontFamily:
      'system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif',
    components: {
      Button: {
        defaultProps: {
          radius: "md",
        },
      },
      Card: {
        defaultProps: {
          radius: "md",
        },
      },
      ActionIcon: {
        defaultProps: {
          variant: "light",
          color: "primary",
        },
      },
      AppShell: {
        defaultProps: {
          styles: {
            main: {
              backgroundColor: currentTheme.colors.background.main,
              color: currentTheme.colors.body,
            },
            header: {
              backgroundColor: currentTheme.colors.background.header,
              color: currentTheme.colors.bodyLight || currentTheme.colors.body,
            },
            navbar: {
              backgroundColor: currentTheme.colors.background.sidebar,
              color: currentTheme.colors.bodyLight || currentTheme.colors.body,
            },
          },
        },
      },
    },
    other: {
      backgroundColors: {
        main: currentTheme.colors.background.main,
        sidebar: currentTheme.colors.background.sidebar,
        header: currentTheme.colors.background.header,
        card: currentTheme.colors.background.card,
      },
      colors: {
        body: currentTheme.colors.body,
        bodyLight: currentTheme.colors.bodyLight,
        border: currentTheme.colors.border,
      },
    },
  });

  return (
    <MantineProvider theme={theme} defaultColorScheme="dark">
      {children}
    </MantineProvider>
  );
}
