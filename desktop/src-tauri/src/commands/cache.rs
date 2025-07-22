use super::command_executor::CommandExecutor;
use super::core::CommandResult;

pub struct CacheCommands<E: CommandExecutor + 'static> {
    executor: std::sync::Arc<E>,
}

impl<E: CommandExecutor + 'static> CacheCommands<E> {
    pub fn new(executor: std::sync::Arc<E>) -> Self {
        Self { executor }
    }

    pub async fn clear(&self) -> CommandResult<String> {
        self.executor.execute::<()>(&["cache", "clear"]).await
    }

    pub async fn set(&self, key: &str, value: &str) -> CommandResult<String> {
        self.executor
            .execute::<()>(&["cache", "set", key, value])
            .await
    }

    pub async fn get(&self, key: &str) -> CommandResult<String> {
        self.executor.execute::<()>(&["cache", "get", key]).await
    }

    pub async fn remove(&self, key: &str) -> CommandResult<String> {
        self.executor.execute::<()>(&["cache", "remove", key]).await
    }

    pub async fn list(&self) -> CommandResult<Vec<String>> {
        self.executor.execute_json(&["cache", "list"]).await
    }
}
