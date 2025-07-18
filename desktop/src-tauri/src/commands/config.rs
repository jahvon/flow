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
        self.executor.execute(&["config", "get", "--output", "json"]).await
    }

    pub async fn set_theme(&self, theme: &str) -> CommandResult<()> {
        self.executor.execute::<()>(&["config", "set", "theme", theme]).await
    }

    pub async fn set_workspace_mode(&self, mode: &str) -> CommandResult<()> {
        self.executor.execute::<()>(&["config", "set", "workspace-mode", mode]).await
    }

    pub async fn set_log_mode(&self, mode: &str) -> CommandResult<()> {
        self.executor.execute::<()>(&["config", "set", "log-mode", mode]).await
    }

    pub async fn set_namespace(&self, namespace: &str) -> CommandResult<()> {
        self.executor.execute::<()>(&["config", "set", "namespace", namespace]).await
    }

    pub async fn reset(&self) -> CommandResult<()> {
        self.executor.execute::<()>(&["config", "reset"]).await
    }
}