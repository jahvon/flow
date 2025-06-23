use std::result::Result;
use std::sync::Arc;
use tauri::Manager;

pub mod cache;
pub mod command_runner;
pub mod types;

pub use cache::Cache;
pub use command_runner::{CommandError, CommandResult, CommandRunner};
pub use types::*;

#[tauri::command]
async fn get_config() -> Result<config::Config, String> {
    let runner = CommandRunner::new();
    runner.get_config().await.map_err(|e| e.to_string())
}

#[tauri::command]
async fn get_workspace(name: String) -> Result<enriched::Workspace, String> {
    let runner = CommandRunner::new();
    runner.get_workspace(&name).await.map_err(|e| e.to_string())
}

#[tauri::command]
async fn list_workspaces() -> Result<Vec<enriched::Workspace>, String> {
    let runner = CommandRunner::new();
    runner.list_workspaces().await.map_err(|e| e.to_string())
}

#[tauri::command]
async fn get_executable(executable_ref: String) -> Result<enriched::Executable, String> {
    let runner = CommandRunner::new();
    runner
        .get_executable(&executable_ref)
        .await
        .map_err(|e| e.to_string())
}

#[tauri::command]
async fn list_executables(
    workspace: Option<String>,
    namespace: Option<String>,
) -> Result<Vec<enriched::Executable>, String> {
    let runner = CommandRunner::new();
    runner
        .list_executables(workspace.as_deref(), namespace.as_deref())
        .await
        .map_err(|e| e.to_string())
}

#[tauri::command]
async fn sync() -> Result<(), String> {
    let runner = CommandRunner::new();
    runner.sync().await.map_err(|e| e.to_string())
}

#[tauri::command]
async fn execute(
    app: tauri::AppHandle,
    verb: String,
    executable_id: String,
    args: Vec<String>,
) -> Result<(), String> {
    let runner = CommandRunner::new();
    let args: Vec<&str> = args.iter().map(|s| s.as_str()).collect();
    runner
        .execute(app, &verb, &executable_id, &args)
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
        .setup(|app| {
            let cache = Arc::new(Cache::new(app.handle().clone()));
            cache.init()?;
            app.manage(cache);
            Ok(())
        })
        .invoke_handler(tauri::generate_handler![
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
