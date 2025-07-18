pub mod cache;
pub mod command_executor;
pub mod core;
pub mod config;
pub mod executable;
pub mod workspace;
pub mod shell;
pub mod vault;

pub struct CommandRunner<E: command_executor::CommandExecutor + 'static> {
    pub config: config::ConfigCommands<E>,
    pub workspace: workspace::WorkspaceCommands<E>,
    pub vault: vault::VaultCommands<E>,
    pub cache: cache::CacheCommands<E>,
}

impl<E: command_executor::CommandExecutor + 'static> CommandRunner<E> {
    pub fn new(executor: std::sync::Arc<E>) -> Self {
        Self {
            config: config::ConfigCommands::new(executor.clone()),
            workspace: workspace::WorkspaceCommands::new(executor.clone()),
            vault: vault::VaultCommands::new(executor.clone()),
            cache: cache::CacheCommands::new(executor),
        }
    }
}

pub struct ExecutableRunner<E: command_executor::ExecutableExecutor + 'static> {
    pub executable: executable::ExecutableCommands,
    pub executor: std::sync::Arc<E>,
}

impl<E: command_executor::ExecutableExecutor + 'static> ExecutableRunner<E> {
    pub fn new(executor: std::sync::Arc<E>) -> Self {
        Self {
            executable: executable::ExecutableCommands::new(),
            executor,
        }
    }
}