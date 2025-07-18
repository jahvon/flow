use std::result::Result;
use tauri::Manager;

pub mod types;
pub mod commands;
pub mod cli;

pub use types::{enriched, generated};
use crate::commands::command_executor::ExecutableExecutor;

pub fn cli_runners() -> cli::Runners {
    let cmd = cli::cmd_executor();
    let exec = cli::exec_executor();
    cli::Runners { cmd, exec }
}

#[tauri::command]
async fn check_flow_binary() -> Result<(), String> {
    let cli = commands::core::CliCommand::new();
    cli.check_binary().await.map_err(|e| e.to_string())
}

#[tauri::command]
async fn get_config() -> Result<crate::types::generated::config::Config, String> {
    let runner = cli_runners();
    runner.cmd.config.get().await.map_err(|e| e.to_string())
}

#[tauri::command]
async fn get_workspace(name: String) -> Result<enriched::Workspace, String> {
    let runner = cli_runners();
    runner.cmd.workspace.get(&name).await.map_err(|e| e.to_string())
}

#[tauri::command]
async fn list_workspaces() -> Result<Vec<enriched::Workspace>, String> {
    let runner = cli_runners();
    runner.cmd.workspace.list().await.map_err(|e| e.to_string())
}

#[tauri::command]
async fn get_executable(executable_ref: String) -> Result<enriched::Executable, String> {
    let runner = cli_runners();
    runner.exec.executable.get(&executable_ref).await.map_err(|e| e.to_string())
}

#[tauri::command]
async fn list_executables(
    workspace: Option<String>,
    namespace: Option<String>,
) -> Result<Vec<enriched::Executable>, String> {
    let runner = cli_runners();
    runner.exec.executable.list(workspace.as_deref(), namespace.as_deref())
        .await
        .map_err(|e| e.to_string())
}

#[tauri::command]
async fn sync() -> Result<(), String> {
    let runner = cli_runners();
    runner.exec.executable.sync().await.map_err(|e| e.to_string())
}

#[tauri::command]
async fn execute(
    app: tauri::AppHandle,
    verb: String,
    executable_id: String,
    args: Vec<String>,
    params: Option<std::collections::HashMap<String, String>>,
) -> Result<(), String> {
    let runner = cli_runners();
    let args: Vec<&str> = args.iter().map(|s| s.as_str()).collect();
    runner.exec.executor.execute::<()>(app, &verb, &executable_id, &args, params)
        .await
        .map_err(|e| e.to_string())
}

#[tauri::command]
async fn reload_window(app: tauri::AppHandle) -> Result<(), String> {
    app.get_webview_window("main")
        .ok_or_else(|| "Main window not found".to_string())?
        .reload()
        .map_err(|e| e.to_string())
}

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .plugin(tauri_plugin_shell::init())
        .plugin(tauri_plugin_opener::init())
        .setup(|_app| {
            Ok(())
        })
        .invoke_handler(tauri::generate_handler![
            check_flow_binary,
            sync,
            execute,
            list_executables,
            get_executable,
            get_workspace,
            list_workspaces,
            get_config,
            reload_window,
        ])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
