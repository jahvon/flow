// Learn more about Tauri commands at https://tauri.app/develop/calling-rust/
use std::process::Stdio;
use std::process::Command;
use serde_json::Result;
use serde::{Deserialize, Serialize};

// Add this to expose the generated types
pub mod generated;

// Re-export for convenience
pub use generated::*;

#[derive(Serialize, Deserialize)]
struct WorkspacesData {
    workspaces: Vec<workspace::Workspace>,
}

#[tauri::command]
fn greet(name: &str) -> String {
    format!("Hello, {}! You've been greeted from Rust!", name)
}

#[tauri::command]
fn workspaces() -> WorkspacesData {
//     let mut output = run_cmd!(flow ws list --output json -x --verbosity -1);
// //     println!("{}", output.ok())
//     // print the output from the above command
//     match output {
//         Ok(result) => Ok(result),
//         Err(err) => eprintln!("Error executing command: {}", err),
//     }
    let output = Command::new("flow")
        .args(["ws", "list", "--output", "json", "-x", "--verbosity", "-1"])
        // Tell the OS to record the command's output
        .stdout(Stdio::piped())
        // execute the command, wait for it to complete, then capture the output
        .output()
        // Blow up if the OS was unable to start the program
        .unwrap();

    // extract the raw bytes that we captured and interpret them as a string
    let stdout = String::from_utf8(output.stdout).unwrap();

    let d: Result<WorkspacesData> = serde_json::from_str(&stdout);
//     match d {
//         Ok(result) => Ok(result),
//         Err(err) => eprintln!("{}", err),
//     }

    return d.unwrap();
}

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
//     run_cmd!(flow ws list --output json -x --verbosity -1);

workspaces();
    tauri::Builder::default()
        .plugin(tauri_plugin_opener::init())
        .invoke_handler(tauri::generate_handler![greet])
        .invoke_handler(tauri::generate_handler![workspaces])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
