pub mod cache;
pub mod command_executor;
pub mod config;
pub mod core;
pub mod executable;
pub mod shell;
pub mod vault;
pub mod workspace;

pub struct CommandRunner<E: command_executor::CommandExecutor + 'static> {
    pub config: config::ConfigCommands<E>,
    pub workspace: workspace::WorkspaceCommands<E>,
    pub vault: vault::VaultCommands<E>,
    pub cache: cache::CacheCommands<E>,
    pub executable: executable::ExecutableCommands<E>,
}

impl<E: command_executor::CommandExecutor + 'static> CommandRunner<E> {
    pub fn new(executor: std::sync::Arc<E>) -> Self {
        Self {
            config: config::ConfigCommands::new(executor.clone()),
            workspace: workspace::WorkspaceCommands::new(executor.clone()),
            vault: vault::VaultCommands::new(executor.clone()),
            cache: cache::CacheCommands::new(executor.clone()),
            executable: executable::ExecutableCommands::new(executor.clone()),
        }
    }
}
