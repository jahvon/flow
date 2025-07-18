use std::process::Stdio;
use serde::Deserialize;
use tauri::Emitter;
use tokio::io::{AsyncBufReadExt, BufReader};
use crate::commands::core::{CliCommand, CliConfig, CommandError, CommandResult};
use crate::types::enriched::Executable;
use tokio::process::Command as TokioCommand;
use log::{info,debug};

#[derive(Deserialize, Debug)]
struct ExecutableResponse {
    executables: Vec<Executable>,
}

pub struct ExecutableCommands{
    config: CliConfig,
    flow: CliCommand,
}

impl ExecutableCommands {
    pub fn new() -> Self {
        Self {
            config: CliConfig::new(),
            flow: CliCommand::new(),
        }
    }


    pub async fn sync(&self) -> CommandResult<()> {
        self.flow.execute::<()>(&["sync"]).await
    }

    pub async fn list(
        &self,
        workspace: Option<&str>,
        namespace: Option<&str>,
    ) -> CommandResult<Vec<Executable>> {
        let mut args = vec!["browse", "--list", "--output", "json"];

        if let Some(ws) = workspace {
            args.extend_from_slice(&["--workspace", ws]);
        }

        if let Some(ns) = namespace {
            args.extend_from_slice(&["--namespace", ns]);
        } else {
            args.push("--all");
        }

        let response: ExecutableResponse = self.flow.execute(&args).await?;
        Ok(response.executables)
    }

    pub async fn get(&self, exec_ref: &str) -> CommandResult<Executable> {
        let split_ref: Vec<&str> = exec_ref.split(" ").collect();
        match split_ref.len() {
            1 => {
                // Just a verb
                self.flow.execute(&["browse", split_ref[0], "--output", "json"])
                    .await
            }
            2 => {
                // Verb and ID
                self.flow.execute(&["browse", split_ref[0], split_ref[1], "--output", "json"])
                    .await
            }
            _ => Err(CommandError::ParseError {
                message: format!("Invalid executable reference format: {}", exec_ref),
                command: format!("{:?}", exec_ref),
                output: String::new(),
            }),
        }
    }

    pub async fn execute<T: for<'de> Deserialize<'de>>(
        &self,
        app: tauri::AppHandle,
        verb: &str,
        executable_id: &str,
        args: &[&str],
        params: Option<std::collections::HashMap<String, String>>,
    ) -> CommandResult<T> {
        let mut cmd_args: Vec<String> = vec![verb.to_string(), executable_id.to_string()];
        cmd_args.extend(args.iter().map(|&s| s.to_string()));

        // Add parameters
        if let Some(params) = params {
            for (key, value) in params {
                cmd_args.push("--param".to_string());
                cmd_args.push(format!("{}={}", key, value));
            }
        }

        let cmd_args_str: Vec<&str> = cmd_args.iter().map(|s| s.as_str()).collect();
        let flow_cmd = self.flow.command_string(&cmd_args_str);

        let mut cmd = if flow_cmd.contains(" && ") {
            // Handle shell sourcing case: "source ~/.bashrc && flow exec ..."
            let mut shell_cmd = TokioCommand::new("sh");
            shell_cmd.arg("-c").arg(&flow_cmd);
            shell_cmd
        } else {
            // Handle simple case: "flow exec ..."
            let parts: Vec<&str> = flow_cmd.split_whitespace().collect();
            if parts.is_empty() {
                return Err(CommandError::ExecutionError("Empty command".to_string()));
            }
            let mut cmd = TokioCommand::new(parts[0]);
            if parts.len() > 1 {
                cmd.args(&parts[1..]);
            }
            cmd
        };
        cmd.stdout(Stdio::piped()).stderr(Stdio::piped());

        let dev_mode = self.config.dev_mode;
        let flow_cmd_clone = flow_cmd.clone();
        if dev_mode {
            info!("Executing with {:?}", &flow_cmd_clone);
        }

        let mut child = cmd
            .spawn()
            .map_err(|e| CommandError::ExecutionError(format!("Failed to spawn command: {}", e)))?;

        let stdout = child
            .stdout
            .take()
            .ok_or_else(|| CommandError::ExecutionError("Failed to capture stdout".to_string()))?;
        let stderr = child
            .stderr
            .take()
            .ok_or_else(|| CommandError::ExecutionError("Failed to capture stderr".to_string()))?;

        let app_stdout = app.clone();
        let app_stderr = app.clone();

        // Handle stdout
        let stdout_handle = tokio::spawn(async move {
            let reader = BufReader::new(stdout);
            let mut lines = reader.lines();

            while let Some(line) = lines.next_line().await.unwrap_or(None) {
                if dev_mode {
                    debug!("stdout: {}", line);
                }

                let _ = app_stdout.emit(
                    "command-output",
                    serde_json::json!({
                    "type": "stdout",
                    "line": line,
                    "timestamp": std::time::SystemTime::now()
                        .duration_since(std::time::UNIX_EPOCH)
                        .unwrap()
                        .as_millis()
                }),
                );
            }
        });

        // Handle stderr
        let stderr_handle = tokio::spawn(async move {
            let reader = BufReader::new(stderr);
            let mut lines = reader.lines();

            while let Some(line) = lines.next_line().await.unwrap_or(None) {
                if dev_mode {
                    debug!("stderr: {}", line);
                }

                let _ = app_stderr.emit(
                    "command-output",
                    serde_json::json!({
                    "type": "stderr",
                    "line": line,
                    "timestamp": std::time::SystemTime::now()
                        .duration_since(std::time::UNIX_EPOCH)
                        .unwrap()
                        .as_millis()
                }),
                );
            }
        });

        let (stdout_result, stderr_result) = tokio::join!(stdout_handle, stderr_handle);
        stdout_result.map_err(|e| CommandError::ExecutionError(format!("Stdout error: {}", e)))?;
        stderr_result.map_err(|e| CommandError::ExecutionError(format!("Stderr error: {}", e)))?;

        // Wait for the process to complete
        let status = child
            .wait()
            .await
            .map_err(|e| CommandError::ExecutionError(e.to_string()))?;

        let success = status.success();
        let exit_code = status.code();

        if dev_mode {
            debug!("Command completed with exit code: {:?}", exit_code);
        }

        // Emit completion event
        let _ = app.emit(
            "command-complete",
            serde_json::json!({
            "success": success,
            "exit_code": exit_code,
            "command": flow_cmd_clone
        }),
        );

        if !success {
            return Err(CommandError::NonZeroExit(
                flow_cmd_clone,
                exit_code.unwrap_or(-1),
                "Command failed - check output for details".to_string(),
            ));
        }

        Ok(serde_json::from_str::<T>("null").map_err(|e| CommandError::ExecutionError(e.to_string()))?)
    }
}