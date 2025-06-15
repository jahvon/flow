import { useCallback, useEffect, useState } from "react";

export type Theme = "everforest" | "dark" | "dracula" | "light" | "tokyo-night";

export function useTheme() {
  const [theme, setThemeState] = useState<Theme>(() => {
    const savedTheme = localStorage.getItem("theme");
    return (savedTheme as Theme) || "everforest";
  });

  const setTheme = useCallback((newTheme: Theme) => {
    setThemeState(newTheme);
    localStorage.setItem("theme", newTheme);
  }, []);

  useEffect(() => {
    document.documentElement.setAttribute("data-palette", theme);
  }, [theme]);

  return { theme, setTheme };
}
