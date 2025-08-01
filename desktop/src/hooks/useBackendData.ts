import { useQuery, useQueryClient } from "@tanstack/react-query";
import { invoke } from "@tauri-apps/api/core";
import React from "react";
import { EnrichedExecutable } from "../types/executable";
import { Config } from "../types/generated/config";
import { EnrichedWorkspace } from "../types/workspace";

export function useConfig(enabled: boolean = true) {
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
    enabled,
  });

  const refreshConfig = () => {
    queryClient.invalidateQueries({ queryKey: ["config"] });
  };

  const updateTheme = async (theme: string) => {
    try {
      await invoke("set_config_theme", { theme });
      refreshConfig();
    } catch (error) {
      throw new Error(`Failed to update theme: ${error}`);
    }
  };

  const updateWorkspaceMode = async (mode: string) => {
    try {
      await invoke("set_config_workspace_mode", { mode });
      refreshConfig();
    } catch (error) {
      throw new Error(`Failed to update workspace mode: ${error}`);
    }
  };

  const updateLogMode = async (mode: string) => {
    try {
      await invoke("set_config_log_mode", { mode });
      refreshConfig();
    } catch (error) {
      throw new Error(`Failed to update log mode: ${error}`);
    }
  };

  const updateNamespace = async (namespace: string) => {
    try {
      await invoke("set_config_namespace", { namespace });
      refreshConfig();
    } catch (error) {
      throw new Error(`Failed to update namespace: ${error}`);
    }
  };

  const updateCurrentWorkspace = async (workspace: string) => {
    try {
      await invoke("set_workspace", { name: workspace, fixed: false });
      refreshConfig();
    } catch (error) {
      throw new Error(`Failed to update current workspace: ${error}`);
    }
  };

  const updateDefaultTimeout = async (timeout: string) => {
    try {
      await invoke("set_config_timeout", { timeout });
      refreshConfig();
    } catch (error) {
      throw new Error(`Failed to update default timeout: ${error}`);
    }
  };

  return {
    config,
    isConfigLoading,
    configError,
    refreshConfig,
    updateTheme,
    updateWorkspaceMode,
    updateLogMode,
    updateNamespace,
    updateCurrentWorkspace,
    updateDefaultTimeout,
  };
}

export function useWorkspaces(enabled: boolean = true) {
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
    enabled,
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

export function useExecutables(selectedWorkspace: string | null, enabled: boolean = true) {
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
    enabled: enabled && !!selectedWorkspace, // Only run when workspace is selected AND enabled
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

// Hook to check if flow binary is available
export function useFlowBinaryCheck() {
  const { data: isFlowBinaryAvailable, isLoading: isCheckingBinary, error: binaryCheckError } = useQuery({
    queryKey: ["flowBinaryCheck"],
    queryFn: async () => {
      try {
        await invoke("check_flow_binary");
        return true;
      } catch (error) {
        console.error(error);
        throw new Error("flow CLI not found or not executable");
      }
    },
    retry: false,
    refetchOnWindowFocus: false,
  });

  return {
    isFlowBinaryAvailable,
    isCheckingBinary,
    binaryCheckError,
  };
}

// Composite hook that combines all data sources
export function useBackendData(selectedWorkspace: string | null) {
  const { isCheckingBinary, binaryCheckError } = useFlowBinaryCheck();
  
  // Only enable other queries if flow binary is available
  const enabled = !binaryCheckError && !isCheckingBinary;
  
  const { config, isConfigLoading, configError, refreshConfig } = useConfig(enabled);
  const {
    workspaces,
    isWorkspacesLoading,
    workspacesError,
    refreshWorkspaces,
  } = useWorkspaces(enabled);
  const {
    executables,
    isExecutablesLoading,
    executablesError,
    refreshExecutables,
  } = useExecutables(selectedWorkspace, enabled);

  // If flow binary is not available, return early with error state
  if (binaryCheckError) {
    return {
      config: null,
      workspaces: [],
      executables: [],
      isLoading: false,
      hasError: binaryCheckError,
      refreshAll: () => {},
      refreshConfig: () => {},
      refreshWorkspaces: () => {},
      refreshExecutables: () => {},
    };
  }

  // If still checking binary, show loading state
  if (isCheckingBinary) {
    return {
      config: null,
      workspaces: [],
      executables: [],
      isLoading: true,
      hasError: null,
      refreshAll: () => {},
      refreshConfig: () => {},
      refreshWorkspaces: () => {},
      refreshExecutables: () => {},
    };
  }

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
