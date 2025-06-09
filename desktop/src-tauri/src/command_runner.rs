use std::process::{Command, Stdio};
use serde::Deserialize;
use std::fmt;
use crate::generated::{flowfile, config};

#[derive(Debug)]
pub enum CommandError {
    ExecutionError(String),
    ParseError(String),
    NonZeroExit(i32),
}

impl fmt::Display for CommandError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            CommandError::ExecutionError(e) => write!(f, "Failed to execute command: {}", e),
            CommandError::ParseError(e) => write!(f, "Failed to parse command output: {}", e),
            CommandError::NonZeroExit(code) => write!(f, "Command returned non-zero exit code: {}", code),
        }
    }
}

impl std::error::Error for CommandError {}

pub type CommandResult<T> = std::result::Result<T, CommandError>;

#[derive(Debug, Clone)]
pub struct CommandRunner;

#[derive(Deserialize)]
struct ExecutableResponse {
    executables: Vec<ExecutableWithId>,
}

#[derive(Deserialize)]
struct ExecutableWithId {
    id: String,
    spec: flowfile::Executable,
}

impl CommandRunner {
    pub fn new() -> Self {
        Self
    }

    fn build_base_command(&self) -> Command {
        // TODO: Make this configurable / use the main flow binary
        let mut cmd = Command::new("/Users/jahvon/workspaces/github.com/jahvon/flow/.bin/flow");
        cmd.stdout(Stdio::piped())
           .stderr(Stdio::piped())
           .arg("-x")  // Always non-interactive
           .arg("--verbosity")
           .arg("-1"); // Always minimum verbosity

        cmd
    }

    pub async fn execute_command<T: for<'de> Deserialize<'de>>(
        &self,
        args: &[&str],
    ) -> CommandResult<T> {
        let mut cmd = self.build_base_command();
        cmd.args(args);

        let output = cmd.output()
            .map_err(|e| CommandError::ExecutionError(e.to_string()))?;

        if !output.status.success() {
            return Err(CommandError::NonZeroExit(output.status.code().unwrap_or(-1)));
        }

        let stdout = String::from_utf8(output.stdout)
            .map_err(|e| CommandError::ParseError(e.to_string()))?;

        println!("Raw command output: {}", stdout);
        println!("Command args: {:?}", args);

        serde_json::from_str(&stdout)
            .map_err(|e| CommandError::ParseError(e.to_string()))
    }

    pub async fn get_config(&self) -> CommandResult<config::Config> {
      self.execute_command(&["config", "view", "--output", "json"]).await
    }

    pub async fn sync(&self) -> CommandResult<()> {
        self.execute_command::<()>(&["sync"]).await
    }

    pub async fn list_executables(
        &self,
        workspace: Option<&str>,
        namespace: Option<&str>,
    ) -> CommandResult<Vec<flowfile::Executable>> {
        let mut args = vec!["library", "glance", "--output", "json"];

        if let Some(ws) = workspace {
            args.extend_from_slice(&["--workspace", ws]);
        }

        if let Some(ns) = namespace {
            args.extend_from_slice(&["--namespace", ns]);
        }

        let response: ExecutableResponse = self.execute_command(&args).await?;
        Ok(response.executables.into_iter().map(|e| e.spec).collect())
    }

    pub async fn get_executable(
        &self,
        exec_ref: &str,
    ) -> CommandResult<flowfile::Executable> {
        self.execute_command(&["library", "view", exec_ref, "--output", "json"]).await
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