use crate::commands::shell::Shell;
use std::process::{Command, Stdio};
use serde::Deserialize;
use std::{env, fmt};
use log::info;

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
        cfg!(debug_assertions) ||
            env::var("TAURI_DEV").is_ok() ||
            env::var("DEV_MODE").map(|v| v == "true").unwrap_or(false)
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
                    info!("flow binary '{}' found and version check passed", flow_binary);
                    Ok(())
                } else {
                    info!("flow binary '{}' found but failed version check", flow_binary);
                    Err(CommandError::ExecutionError(format!(
                        "flow binary '{}' found but failed version check",
                        flow_binary
                    )))
                }
            }
            Err(e) => {
                info!("flow binary '{}' not found or not executable: {}", flow_binary, e);
                Err(CommandError::ExecutionError(format!(
                    "flow binary '{}' not found or not executable: {}",
                    flow_binary, e
                )))
            }
        }
    }

    pub fn command_string(&self,  flow_args: &[&str]) -> String {
        let flow_cmd = format!(
            "{} {}",
            self.config.flow_binary(),
            flow_args.join(" ")
        );

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

    pub async fn execute<T: for<'de> Deserialize<'de>>(
        &self,
        args: &[&str],
    ) -> CommandResult<T> {
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

        serde_json::from_str(&stdout).map_err(|e| CommandError::ParseError {
            message: e.to_string(),
            command: format!("{:?}", args),
            output: stdout.clone(),
        })
    }
}