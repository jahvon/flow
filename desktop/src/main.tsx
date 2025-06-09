import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App";
import { MantineProvider, ColorSchemeScript } from '@mantine/core';
import { theme } from './theme';

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
  <React.StrictMode>
    <MantineProvider theme={theme} defaultColorScheme="auto">
      <ColorSchemeScript />
      <App />
    </MantineProvider>
  </React.StrictMode>,
);
