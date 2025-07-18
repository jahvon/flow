use super::command_executor::CommandExecutor;
use super::core::CommandResult;

pub struct VaultCommands<E: CommandExecutor + 'static> {
    executor: std::sync::Arc<E>,
}

impl<E: CommandExecutor + 'static> VaultCommands<E> {
    pub fn new(executor: std::sync::Arc<E>) -> Self {
        Self { executor }
    }

    pub async fn create(&self, name: &str) -> CommandResult<()> {
        self.executor.execute::<()>(&["vault", "create", name]).await
    }

    pub async fn list(&self) -> CommandResult<Vec<String>> {
        self.executor.execute(&["vault", "list", "--output", "json"]).await
    }

    pub async fn get(&self, name: &str) -> CommandResult<String> {
        self.executor.execute(&["vault", "get", name]).await
    }

    pub async fn switch(&self, name: &str) -> CommandResult<()> {
        self.executor.execute::<()>(&["vault", "switch", name]).await
    }

    pub async fn remove(&self, name: &str) -> CommandResult<()> {
        self.executor.execute::<()>(&["vault", "remove", name]).await
    }
}