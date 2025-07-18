use super::command_executor::CommandExecutor;
use super::core::CommandResult;
use crate::types::enriched::Workspace;
use serde::Deserialize;

#[derive(Deserialize, Debug)]
struct WorkspaceResponse {
    workspaces: Vec<Workspace>,
}

pub struct WorkspaceCommands<E: CommandExecutor + 'static> {
    executor: std::sync::Arc<E>,
}

impl<E: CommandExecutor + 'static> WorkspaceCommands<E> {
    pub fn new(executor: std::sync::Arc<E>) -> Self {
        Self { executor }
    }

    pub async fn list(&self) -> CommandResult<Vec<Workspace>> {
        let response: WorkspaceResponse = self.executor
            .execute(&["workspace", "list", "--output", "json"]).await?;
        Ok(response.workspaces)
    }

    pub async fn get(&self, workspace: &str) -> CommandResult<Workspace> {
        self.executor
            .execute(&["workspace", "get", workspace, "--output", "json"]).await
    }

    pub async fn add(&self, name: &str, path: &str, set_current: bool) -> CommandResult<()> {
        let mut args = vec!["workspace", "add", name, path];
        if set_current {
            args.push("--set");
        }
        let args_ref: Vec<&str> = args.iter().map(|s| *s).collect();
        self.executor.execute(&args_ref).await
    }

    pub async fn switch(&self, name: &str, fixed: bool) -> CommandResult<()> {
        let mut args = vec!["workspace", "switch", name];
        if fixed {
            args.push("--fixed");
        }
        let args_ref: Vec<&str> = args.iter().map(|s| *s).collect();
        self.executor.execute(&args_ref).await
    }

    pub async fn remove(&self, name: &str) -> CommandResult<()> {
        self.executor.execute(&["workspace", "remove", name]).await
    }
}