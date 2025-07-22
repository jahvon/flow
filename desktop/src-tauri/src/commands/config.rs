use super::command_executor::CommandExecutor;
use super::core::CommandResult;
use crate::types::generated::config::Config;

pub struct ConfigCommands<E: CommandExecutor + 'static> {
    executor: std::sync::Arc<E>,
}

impl<E: CommandExecutor + 'static> ConfigCommands<E> {
    pub fn new(executor: std::sync::Arc<E>) -> Self {
        Self { executor }
    }

    pub async fn get(&self) -> CommandResult<Config> {
        self.executor.execute_json(&["config", "get"]).await
    }

    pub async fn set_theme(&self, theme: &str) -> CommandResult<String> {
        self.executor
            .execute::<()>(&["config", "set", "theme", theme])
            .await
    }

    pub async fn set_workspace_mode(&self, mode: &str) -> CommandResult<String> {
        self.executor
            .execute::<()>(&["config", "set", "workspace-mode", mode])
            .await
    }

    pub async fn set_log_mode(&self, mode: &str) -> CommandResult<String> {
        self.executor
            .execute::<()>(&["config", "set", "log-mode", mode])
            .await
    }

    pub async fn set_namespace(&self, namespace: &str) -> CommandResult<String> {
        let namespace = if namespace.is_empty() {
            "\"\""
        } else {
            namespace
        };
        self.executor
            .execute::<()>(&["config", "set", "namespace", namespace])
            .await
    }

    pub async fn set_timeout(&self, timeout: &str) -> CommandResult<String> {
        let timeout = if timeout.is_empty() {
            "0"
        } else {
            timeout
        };
        self.executor
            .execute::<()>(&["config", "set", "timeout", timeout])
            .await
    }

    pub async fn reset(&self) -> CommandResult<String> {
        self.executor.execute::<()>(&["config", "reset"]).await
    }
}
