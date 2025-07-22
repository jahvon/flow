use crate::commands::shell::Shell;
use log::{debug, info};
use serde::Deserialize;
use std::process::{Command, Stdio};
use std::{env, fmt};
use tauri::Emitter;
use tokio::io::{AsyncBufReadExt, BufReader};
use tokio::process::Command as TokioCommand;

#[derive(Debug, Clone)]
pub struct CliConfig {
    pub flow_binary_path: String,
    pub dev_mode: bool,
}

impl CliConfig {
    pub fn new() -> Self {
        let dev_mode = Self::is_dev_mode();
        let flow_binary_path = Self::resolve_binary_path(dev_mode);

        Self {
            flow_binary_path,
            dev_mode,
        }
    }

    fn is_dev_mode() -> bool {
        cfg!(debug_assertions)
            || env::var("TAURI_DEV").is_ok()
            || env::var("DEV_MODE").map(|v| v == "true").unwrap_or(false)
    }

    fn resolve_binary_path(dev_mode: bool) -> String {
        if dev_mode {
            if let Ok(custom_path) = env::var("FLOW_BINARY_PATH") {
                info!("Using custom flow binary path: {}", custom_path);
                return custom_path;
            }
        }

        // Default to system flow binary
        "flow".to_string()
    }

    pub fn flow_binary(&self) -> &str {
        &self.flow_binary_path
    }
}

impl Default for CliConfig {
    fn default() -> Self {
        Self::new()
    }
}

#[derive(Debug)]
pub enum CommandError {
    ExecutionError(String),
    ParseError {
        message: String,
        command: String,
        output: String,
    },
    NonZeroExit(String, i32, String),
}

impl fmt::Display for CommandError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            CommandError::ExecutionError(e) => write!(f, "Failed to execute command: {}", e),
            CommandError::ParseError {
                message,
                command,
                output,
            } => {
                write!(
                    f,
                    "Failed to parse command output for '{}': {}\nOutput: {}",
                    command, message, output
                )
            }
            CommandError::NonZeroExit(command, code, output) => {
                write!(
                    f,
                    "Command '{}' returned non-zero exit code: {}\nOutput: {}",
                    command, code, output
                )
            }
        }
    }
}

impl std::error::Error for CommandError {}

pub type CommandResult<T> = std::result::Result<T, CommandError>;

pub struct CliCommand {
    config: CliConfig,
    shell: Shell,
}

impl CliCommand {
    pub fn new() -> Self {
        Self {
            config: CliConfig::new(),
            shell: Shell::detect(),
        }
    }

    pub fn with_config(config: CliConfig) -> Self {
        Self {
            config,
            shell: Shell::detect(),
        }
    }

    pub async fn check_binary(&self) -> CommandResult<()> {
        use tokio::process::Command;

        let flow_binary = self.config.flow_binary();
        let mut cmd = Command::new(flow_binary);
        cmd.arg("--version");

        match cmd.output().await {
            Ok(output) => {
                if output.status.success() {
                    if !self.config.dev_mode {
                        let version = String::from_utf8_lossy(&output.stdout);
                        info!("Using flow binary version: {}", version);
                    }
                    info!(
                        "flow binary '{}' found and version check passed",
                        flow_binary
                    );
                    Ok(())
                } else {
                    info!(
                        "flow binary '{}' found but failed version check",
                        flow_binary
                    );
                    Err(CommandError::ExecutionError(format!(
                        "flow binary '{}' found but failed version check",
                        flow_binary
                    )))
                }
            }
            Err(e) => {
                info!(
                    "flow binary '{}' not found or not executable: {}",
                    flow_binary, e
                );
                Err(CommandError::ExecutionError(format!(
                    "flow binary '{}' not found or not executable: {}",
                    flow_binary, e
                )))
            }
        }
    }

    pub fn command_string(&self, flow_args: &[&str]) -> String {
        let flow_cmd = format!("{} {}", self.config.flow_binary(), flow_args.join(" "));

        let full_command = if let Some(profile_cmd) = self.shell.source_profile_command() {
            format!("{} && {}", profile_cmd, flow_cmd)
        } else {
            flow_cmd
        };

        full_command
    }

    pub fn build(&self, flow_args: &[&str]) -> Command {
        let flow_cmd = self.command_string(&flow_args);
        let (shell_executable, shell_args) = self.shell.command_args(&flow_cmd);

        let mut cmd = Command::new(shell_executable);
        cmd.args(shell_args);
        cmd.stdout(Stdio::piped()).stderr(Stdio::piped());

        if self.config.dev_mode {
            info!("Executing via shell {:?}: {}", self.shell, flow_cmd);
        }

        cmd
    }

    pub async fn execute(&self, args: &[&str]) -> CommandResult<String> {
        let mut cmd = self.build(args);

        let output = cmd
            .output()
            .map_err(|e| CommandError::ExecutionError(e.to_string()))?;

        if !output.status.success() {
            let stderr = String::from_utf8_lossy(&output.stderr);
            let stdout = String::from_utf8_lossy(&output.stdout);
            return Err(CommandError::NonZeroExit(
                format!("{:?}", args),
                output.status.code().unwrap_or(-1),
                format!("stdout: {}\nstderr: {}", stdout, stderr),
            ));
        }

        let stdout = String::from_utf8(output.stdout).map_err(|e| CommandError::ParseError {
            message: e.to_string(),
            command: format!("{:?}", args),
            output: String::new(),
        })?;

        Ok(stdout)
    }

    pub async fn execute_json<T: for<'de> Deserialize<'de>>(
        &self,
        args: &[&str],
    ) -> CommandResult<T> {
        let mut args = args.to_vec();
        if !args.contains(&"--output") {
            args.push("--output");
            args.push("json");
        }

        let mut cmd = self.build(args.as_slice());

        let output = cmd
            .output()
            .map_err(|e| CommandError::ExecutionError(e.to_string()))?;

        if !output.status.success() {
            let stderr = String::from_utf8_lossy(&output.stderr);
            let stdout = String::from_utf8_lossy(&output.stdout);
            return Err(CommandError::NonZeroExit(
                format!("{:?}", args),
                output.status.code().unwrap_or(-1),
                format!("stdout: {}\nstderr: {}", stdout, stderr),
            ));
        }

        let stdout = String::from_utf8(output.stdout).map_err(|e| CommandError::ParseError {
            message: e.to_string(),
            command: format!("{:?}", args),
            output: String::new(),
        })?;

        serde_json::from_str(&stdout).map_err(|e| CommandError::ParseError {
            message: e.to_string(),
            command: format!("{:?}", args),
            output: stdout.clone(),
        })
    }

    pub async fn execute_executable<T: for<'de> Deserialize<'de>>(
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
        let flow_cmd = self.command_string(&cmd_args_str);

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

        Ok(serde_json::from_str::<T>("null")
            .map_err(|e| CommandError::ExecutionError(e.to_string()))?)
    }
}
