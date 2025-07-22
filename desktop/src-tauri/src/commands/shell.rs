use std::env;

#[derive(Debug, Clone)]
pub enum Shell {
    Bash,
    Zsh,
    Fish,
    Unknown(String),
}

impl Shell {
    pub fn detect() -> Self {
        if let Ok(shell_path) = env::var("SHELL") {
            return Self::from_path(&shell_path);
        }

        Shell::Bash
    }

    fn from_path(path: &str) -> Self {
        let shell_name = std::path::Path::new(path)
            .file_name()
            .and_then(|name| name.to_str())
            .unwrap_or(path);

        match shell_name {
            "bash" => Shell::Bash,
            "zsh" => Shell::Zsh,
            "fish" => Shell::Fish,
            _ => Shell::Unknown(path.to_string()),
        }
    }

    pub fn command_args(&self, command: &str) -> (String, Vec<String>) {
        match self {
            Shell::Bash | Shell::Zsh | Shell::Fish => (
                self.executable(),
                vec!["-c".to_string(), command.to_string()],
            ),
            Shell::Unknown(path) => {
                // Assume POSIX-like behavior
                (path.clone(), vec!["-c".to_string(), command.to_string()])
            }
        }
    }

    pub fn executable(&self) -> String {
        match self {
            Shell::Bash => "bash".to_string(),
            Shell::Zsh => "zsh".to_string(),
            Shell::Fish => "fish".to_string(),
            Shell::Unknown(path) => path.clone(),
        }
    }

    pub fn source_profile_command(&self) -> Option<String> {
        match self {
            Shell::Bash => Some(
                "source ~/.bashrc 2>/dev/null || source ~/.bash_profile 2>/dev/null || true"
                    .to_string(),
            ),
            Shell::Zsh => Some("source ~/.zshrc 2>/dev/null || true".to_string()),
            Shell::Fish => {
                Some("source ~/.config/fish/config.fish 2>/dev/null || true".to_string())
            }
            Shell::Unknown(_) => None,
        }
    }
}
