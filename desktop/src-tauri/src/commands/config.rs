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

    pub async fn set_current_workspace(&self, workspace: &str) -> CommandResult<()> {
        self.executor.execute::<()>(&["config", "set", "current-workspace", workspace]).await
    }

    pub async fn set_current_vault(&self, vault: &str) -> CommandResult<()> {
        self.executor.execute::<()>(&["config", "set", "current-vault", vault]).await
    }

    pub async fn set_default_timeout(&self, timeout: &str) -> CommandResult<()> {
        self.executor.execute::<()>(&["config", "set", "default-timeout", timeout]).await
    }

    pub async fn add_workspace(&self, name: &str, path: &str) -> CommandResult<()> {
        self.executor.execute::<()>(&["config", "workspace", "add", name, path]).await
    }

    pub async fn remove_workspace(&self, name: &str) -> CommandResult<()> {
        self.executor.execute::<()>(&["config", "workspace", "remove", name]).await
    }

    pub async fn add_vault(&self, name: &str, path: &str) -> CommandResult<()> {
        self.executor.execute::<()>(&["config", "vault", "add", name, path]).await
    }

    pub async fn remove_vault(&self, name: &str) -> CommandResult<()> {
        self.executor.execute::<()>(&["config", "vault", "remove", name]).await
    }

    pub async fn add_template(&self, name: &str, path: &str) -> CommandResult<()> {
        self.executor.execute::<()>(&["config", "template", "add", name, path]).await
    }

    pub async fn remove_template(&self, name: &str) -> CommandResult<()> {
        self.executor.execute::<()>(&["config", "template", "remove", name]).await
    }

    pub async fn reset(&self) -> CommandResult<()> {
        self.executor.execute::<()>(&["config", "reset"]).await
    }
}