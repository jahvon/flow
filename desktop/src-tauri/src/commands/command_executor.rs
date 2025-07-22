use crate::commands::core::CommandResult;
use async_trait::async_trait;
use serde::de::DeserializeOwned;

#[async_trait]
pub trait CommandExecutor: Send + Sync {
    async fn execute<T: DeserializeOwned + Send>(&self, args: &[&str]) -> CommandResult<String>;
    async fn execute_json<T: DeserializeOwned + Send>(&self, args: &[&str]) -> CommandResult<T>;
    async fn execute_executable<T: DeserializeOwned + Send>(
        &self,
        app: tauri::AppHandle,
        verb: &str,
        executable_id: &str,
        args: &[&str],
        params: Option<std::collections::HashMap<String, String>>,
    ) -> CommandResult<T>;
}
