use async_trait::async_trait;
use serde::de::DeserializeOwned;
use crate::commands::core::{CommandResult};

#[async_trait]
pub trait CommandExecutor: Send + Sync {
    async fn execute<T: DeserializeOwned + Send>(&self, args: &[&str]) -> CommandResult<T>;
}

#[async_trait]
pub trait ExecutableExecutor: Send + Sync {
    async fn execute<T: DeserializeOwned + Send>(
        &self,
        app: tauri::AppHandle,
        verb: &str,
        executable_id: &str,
        args: &[&str],
        params: Option<std::collections::HashMap<String, String>>,
    ) -> CommandResult<T>;
}
