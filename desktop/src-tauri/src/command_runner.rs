use crate::types::enriched::{Executable, Workspace};
use crate::types::generated::config::Config;
use serde::Deserialize;
use std::fmt;
use std::process::{Command, Stdio};

#[derive(Debug)]
pub enum CommandError {
    ExecutionError(String),
    ParseError {
        message: String,
        command: String,
        output: String,
    },
    NonZeroExit(i32),
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
            CommandError::NonZeroExit(code) => {
                write!(f, "Command returned non-zero exit code: {}", code)
            }
        }
    }
}

impl std::error::Error for CommandError {}

pub type CommandResult<T> = std::result::Result<T, CommandError>;

#[derive(Debug, Clone)]
pub struct CommandRunner;

#[derive(Deserialize, Debug)]
struct ExecutableResponse {
    executables: Vec<Executable>,
}

#[derive(Deserialize, Debug)]
struct WorkspaceResponse {
    workspaces: Vec<Workspace>,
}

impl CommandRunner {
    pub fn new() -> Self {
        Self
    }

    fn build_base_command(&self) -> Command {
        // TODO: Make this configurable / use the main flow binary
        let mut cmd = Command::new("/Users/jahvon/workspaces/github.com/jahvon/flow/.bin/flow");
        cmd.stdout(Stdio::piped())
            .stderr(Stdio::piped());

        cmd
    }

    pub async fn execute_command<T: for<'de> Deserialize<'de>>(
        &self,
        args: &[&str],
    ) -> CommandResult<T> {
        let mut cmd = self.build_base_command();
        cmd.args(args);

        println!("cmd: {:?}", cmd);

        let output = cmd
            .output()
            .map_err(|e| CommandError::ExecutionError(e.to_string()))?;

        if !output.status.success() {
            return Err(CommandError::NonZeroExit(
                output.status.code().unwrap_or(-1),
            ));
        }

        let stdout = String::from_utf8(output.stdout).map_err(|e| CommandError::ParseError {
            message: e.to_string(),
            command: format!("{:?}", cmd),
            output: String::new(),
        })?;

        serde_json::from_str(&stdout).map_err(|e| CommandError::ParseError {
            message: e.to_string(),
            command: format!("{:?}", cmd),
            output: stdout.clone(),
        })
    }

    pub async fn get_config(&self) -> CommandResult<Config> {
        self.execute_command(&["config", "get", "--output", "json"])
            .await
    }

    pub async fn sync(&self) -> CommandResult<()> {
        self.execute_command::<()>(&["sync"]).await
    }

    pub async fn list_workspaces(&self) -> CommandResult<Vec<Workspace>> {
        let response: WorkspaceResponse = self
            .execute_command(&["workspace", "list", "--output", "json"])
            .await?;
        Ok(response.workspaces)
    }

    pub async fn get_workspace(&self, workspace: &str) -> CommandResult<Workspace> {
        let response: Workspace = self
            .execute_command(&["workspace", "get", workspace, "--output", "json"])
            .await?;
        Ok(response)
    }

    pub async fn list_executables(
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

        let response: ExecutableResponse = self.execute_command(&args).await?;
        Ok(response.executables)
    }

    pub async fn get_executable(&self, exec_ref: &str) -> CommandResult<Executable> {
        let split_ref: Vec<&str> = exec_ref.split(" ").collect();
        match split_ref.len() {
            1 => {
                // Just a verb
                self.execute_command(&["browse", split_ref[0], "--output", "json"])
                    .await
            }
            2 => {
                // Verb and ID
                self.execute_command(&[
                    "browse",
                    split_ref[0],
                    split_ref[1],
                    "--output",
                    "json",
                ])
                .await
            }
            _ => Err(CommandError::ParseError {
                message: format!("Invalid executable reference format: {}", exec_ref),
                command: format!("{:?}", exec_ref),
                output: String::new(),
            }),
        }
    }

    pub async fn execute(
        &self,
        verb: &str,
        executable_id: &str,
        args: &[&str],
    ) -> CommandResult<()> {
        let mut cmd_args = vec![verb, executable_id];
        cmd_args.extend(args);
        self.execute_command::<()>(&cmd_args).await
    }
}
