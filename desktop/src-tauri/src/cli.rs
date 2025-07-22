use crate::commands::command_executor::CommandExecutor;
use crate::commands::CommandRunner;
use serde::de::DeserializeOwned;
use std::sync::Arc;

pub struct FlowCLI;

#[async_trait::async_trait]
impl CommandExecutor for FlowCLI {
    async fn execute<T: DeserializeOwned + Send>(
        &self,
        args: &[&str],
    ) -> crate::commands::core::CommandResult<String> {
        let cli = crate::commands::core::CliCommand::new();
        cli.execute(args).await
    }
    async fn execute_json<T: DeserializeOwned + Send>(
        &self,
        args: &[&str],
    ) -> crate::commands::core::CommandResult<T> {
        let cli = crate::commands::core::CliCommand::new();
        cli.execute_json(args).await
    }
    async fn execute_executable<T: DeserializeOwned + Send>(
        &self,
        app: tauri::AppHandle,
        verb: &str,
        executable_id: &str,
        args: &[&str],
        params: Option<std::collections::HashMap<String, String>>,
    ) -> crate::commands::core::CommandResult<T> {
        let cli = crate::commands::core::CliCommand::new();
        cli.execute_executable(app, verb, executable_id, args, params)
            .await
    }
}

pub fn cli_executor() -> CommandRunner<FlowCLI> {
    let executor = Arc::new(FlowCLI);
    CommandRunner::new(executor)
}
