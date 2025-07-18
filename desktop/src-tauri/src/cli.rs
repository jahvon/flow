use std::sync::Arc;
use crate::commands::command_executor::{CommandExecutor, ExecutableExecutor};
use crate::commands::{CommandRunner, ExecutableRunner};

pub struct CmdExecutor;

#[async_trait::async_trait]
impl CommandExecutor for CmdExecutor {
    async fn execute<T: serde::de::DeserializeOwned + Send>(&self, args: &[&str]) -> crate::commands::core::CommandResult<T> {
        let cli = crate::commands::core::CliCommand::new();
        cli.execute(args).await
    }
}

pub struct ExecExecutor;

#[async_trait::async_trait]
impl ExecutableExecutor for ExecExecutor {
    async fn execute<T: serde::de::DeserializeOwned + Send>(
        &self,
        app: tauri::AppHandle,
        verb: &str,
        executable_id: &str,
        args: &[&str],
        params: Option<std::collections::HashMap<String, String>>,
    ) -> crate::commands::core::CommandResult<T> {
        let runner = crate::commands::executable::ExecutableCommands::new();
        runner.execute(app, verb, executable_id, args, params).await
    }
}

pub struct Runners {
    pub cmd: CommandRunner<CmdExecutor>,
    pub exec: ExecutableRunner<ExecExecutor>,
}

pub fn cmd_executor() -> CommandRunner<CmdExecutor> {
    let executor = Arc::new(CmdExecutor);
    CommandRunner::new(executor)
}

pub fn exec_executor() -> ExecutableRunner<ExecExecutor> {
    let executor = Arc::new(ExecExecutor);
    ExecutableRunner::new(executor)
}
