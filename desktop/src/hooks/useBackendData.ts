import { useQuery, useQueryClient } from "@tanstack/react-query";
import { invoke } from "@tauri-apps/api/core";
import React from "react";
import { EnrichedExecutable } from "../types/executable";
import { Config } from "../types/generated/config";
import { EnrichedWorkspace } from "../types/workspace";

export function useConfig() {
  const queryClient = useQueryClient();

  const {
    data: config,
    isLoading: isConfigLoading,
    error: configError,
  } = useQuery({
    queryKey: ["config"],
    queryFn: async () => {
      return await invoke<Config>("get_config");
    },
  });

  const refreshConfig = () => {
    queryClient.invalidateQueries({ queryKey: ["config"] });
  };

  return {
    config,
    isConfigLoading,
    configError,
    refreshConfig,
  };
}

export function useWorkspaces() {
  const queryClient = useQueryClient();

  const {
    data: workspaces,
    isLoading: isWorkspacesLoading,
    error: workspacesError,
  } = useQuery({
    queryKey: ["workspaces"],
    queryFn: async () => {
      const response = await invoke<EnrichedWorkspace[]>("list_workspaces");
      return response;
    },
  });

  const refreshWorkspaces = () => {
    queryClient.invalidateQueries({ queryKey: ["workspaces"] });
  };

  return {
    workspaces,
    isWorkspacesLoading,
    workspacesError,
    refreshWorkspaces,
  };
}

export function useExecutable(executableRef: string) {
  const queryClient = useQueryClient();
  const [currentExecutable, setCurrentExecutable] =
    React.useState<EnrichedExecutable | null>(null);

  const {
    data: executable,
    isLoading: isExecutableLoading,
    error: executableError,
  } = useQuery({
    queryKey: ["executable", executableRef],
    queryFn: async () => {
      if (!executableRef) return null;
      const response = await invoke<EnrichedExecutable>("get_executable", {
        executableRef: executableRef,
      });
      return response;
    },
    enabled: !!executableRef,
  });

  // Update current executable when we have new data
  React.useEffect(() => {
    if (executable) {
      setCurrentExecutable(executable);
    }
  }, [executable]);

  const refreshExecutable = () => {
    if (executableRef) {
      queryClient.invalidateQueries({
        queryKey: ["executable", executableRef],
      });
    }
  };

  return {
    executable: currentExecutable,
    isExecutableLoading,
    executableError,
    refreshExecutable,
  };
}

export function useExecutables(selectedWorkspace: string | null) {
  const queryClient = useQueryClient();

  const {
    data: executables,
    isLoading: isExecutablesLoading,
    error: executablesError,
  } = useQuery({
    queryKey: ["executables", selectedWorkspace],
    queryFn: async () => {
      if (!selectedWorkspace) return [];
      const response = await invoke<EnrichedExecutable[]>("list_executables", {
        workspace: selectedWorkspace,
      });
      return response;
    },
    enabled: !!selectedWorkspace, // Only run when workspace is selected
  });

  const refreshExecutables = () => {
    if (selectedWorkspace) {
      queryClient.invalidateQueries({
        queryKey: ["executables", selectedWorkspace],
      });
    }
  };

  return {
    executables: executables || [],
    isExecutablesLoading,
    executablesError,
    refreshExecutables,
  };
}

// Composite hook that combines all data sources
export function useWorkspaceData(selectedWorkspace: string | null) {
  const { config, isConfigLoading, configError, refreshConfig } = useConfig();
  const {
    workspaces,
    isWorkspacesLoading,
    workspacesError,
    refreshWorkspaces,
  } = useWorkspaces();
  const {
    executables,
    isExecutablesLoading,
    executablesError,
    refreshExecutables,
  } = useExecutables(selectedWorkspace);

  const isLoading =
    isConfigLoading || isWorkspacesLoading || isExecutablesLoading;
  const hasError = configError || workspacesError || executablesError;
  if (hasError) {
    console.error("Error", hasError);
  }

  const refreshAll = () => {
    refreshConfig();
    refreshWorkspaces();
    refreshExecutables();
  };

  return {
    config,
    workspaces,
    executables,
    isLoading,
    hasError,
    refreshAll,
    refreshConfig,
    refreshWorkspaces,
    refreshExecutables,
  };
}
