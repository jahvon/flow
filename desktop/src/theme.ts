import { createTheme } from '@mantine/core';

export const theme = createTheme({
  colors: {
    palette: [
      '#343F44',
      '#5C6A72',
      '#7FBBB3',
      '#83C092',
      '#A7C080',
      '#DBBC7F',
      '#D699B6',
      '#E67E80',
      '#F85552',
      '#DFDDC8',
    ],
    dark: [
      '#C1C2C5',
      '#A6A7AB',
      '#909296',
      '#5C5F66',
      '#373A40',
      '#2C2E33',
      '#25262B',
      '#1A1B1E',
      '#141517',
      '#101113',
    ],
    accent: ['#3A94C5', '#35A77C', '#8DA101', '#DFA000', '#DF69BA', '#D3C6AA', '#D3C6AA', '#D3C6AA', '#D3C6AA', '#D3C6AA'],
    red: ['#F85552', '#F85552', '#F85552', '#F85552', '#F85552', '#F85552', '#F85552', '#F85552', '#F85552', '#F85552'],
  },
  primaryColor: 'palette',
  primaryShade: 0,
  defaultRadius: 'md',
  fontFamily: 'system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif',
  components: {
    Button: {
      defaultProps: {
        radius: 'md',
      },
    },
    Card: {
      defaultProps: {
        radius: 'md',
      },
    },
    ActionIcon: {
      defaultProps: {
        variant: 'light',
        color: 'palette',
      },
      styles: {
        root: {
          '&[data-variant="light"]': {
            backgroundColor: 'var(--mantine-color-palette-0)',
            color: 'var(--mantine-color-palette-9)',
            '&:hover': {
              backgroundColor: 'var(--mantine-color-palette-1)',
            },
          },
        },
      },
    },
  },
  other: {
    darkModeColors: {
      background: 'var(--mantine-color-dark-7)',
      text: 'var(--mantine-color-dark-0)',
      border: 'var(--mantine-color-dark-4)',
    },
  },
});