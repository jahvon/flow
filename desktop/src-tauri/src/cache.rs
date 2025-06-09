use std::collections::HashMap;
use std::path::PathBuf;
use std::sync::RwLock;
use serde::{Deserialize, Serialize};
use tauri::AppHandle;
use tauri::Emitter;
use notify::{Watcher, RecursiveMode, Event};
use std::sync::Arc;
use std::env;

use crate::generated::workspace::Workspace;

const WS_CACHE_KEY: &str = "workspace";
const CACHE_DIR: &str = "flow";
const ENV_CACHE_DIR: &str = "FLOW_CACHE_DIR";

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct WorkspaceCacheData {
    pub workspaces: HashMap<String, Workspace>,
    #[serde(rename = "workspaceLocations")]
    pub workspace_locations: HashMap<String, String>,
}

// At this time, we only care about the workspace cache since it gives us
// the assigned name and path to the workspaces. Executables should be loaded
// from the command runner to avoid duplicating core logic.
pub struct Cache {
    workspace_data: Arc<RwLock<WorkspaceCacheData>>,
    app_handle: AppHandle,
}

impl Cache {
    pub fn new(app_handle: AppHandle) -> Self {
        Self {
            workspace_data: Arc::new(RwLock::new(WorkspaceCacheData {
                workspaces: HashMap::new(),
                workspace_locations: HashMap::new(),
            })),
            app_handle,
        }
    }

    pub fn init(&self) -> Result<(), String> {
        // Load initial cache data
        self.load_workspace_cache()?;

        // Set up file watchers
        self.setup_watchers()?;

        Ok(())
    }

    fn load_workspace_cache(&self) -> Result<(), String> {
        let cache_path = self.get_cache_path(WS_CACHE_KEY)?;
        if !cache_path.exists() {
            return Ok(());
        }

        let data = std::fs::read_to_string(cache_path)
            .map_err(|e| format!("Failed to read workspace cache: {}", e))?;

        let cache_data: WorkspaceCacheData = serde_yaml::from_str(&data)
            .map_err(|e| format!("Failed to parse workspace cache: {}", e))?;

        let mut workspace_data = self.workspace_data.write()
            .map_err(|_| "Failed to acquire workspace cache lock".to_string())?;
        *workspace_data = cache_data;

        Ok(())
    }

    fn setup_watchers(&self) -> Result<(), String> {
        let workspace_data = self.workspace_data.clone();
        let app_handle = self.app_handle.clone();

        // Watch workspace cache
        let ws_cache_path = self.get_cache_path(WS_CACHE_KEY)?;
        let ws_cache_path_for_watch = ws_cache_path.clone();
        let mut ws_watcher = notify::recommended_watcher(move |res: Result<Event, _>| {
            if let Ok(event) = res {
                if event.kind.is_modify() {
                    if let Ok(mut data) = workspace_data.write() {
                        if let Ok(cache_str) = std::fs::read_to_string(&ws_cache_path) {
                            if let Ok(cache_data) = serde_yaml::from_str::<WorkspaceCacheData>(&cache_str) {
                                *data = cache_data;
                                let _ = app_handle.emit("workspace-cache-updated", ());
                            }
                        }
                    }
                }
            }
        }).map_err(|e| format!("Failed to create workspace cache watcher: {}", e))?;

        ws_watcher.watch(&ws_cache_path_for_watch, RecursiveMode::NonRecursive)
            .map_err(|e| format!("Failed to watch workspace cache: {}", e))?;

        Ok(())
    }

    fn get_cache_path(&self, key: &str) -> Result<PathBuf, String> {
        let cache_dir = if let Ok(custom_dir) = env::var(ENV_CACHE_DIR) {
            PathBuf::from(custom_dir)
        } else {
            // Use the same path as Go's os.UserCacheDir()
            let home = dirs::home_dir()
                .ok_or_else(|| "Failed to get home directory".to_string())?;

            #[cfg(target_os = "macos")]
            let cache_dir = home.join("Library/Caches");

            #[cfg(target_os = "linux")]
            let cache_dir = home.join(".cache");

            #[cfg(target_os = "windows")]
            let cache_dir = home.join("AppData/Local");

            cache_dir.join(CACHE_DIR)
        };

        std::fs::create_dir_all(&cache_dir)
            .map_err(|e| format!("Failed to create cache directory: {}", e))?;

        Ok(cache_dir.join(format!("latestcache/{}", key)))
    }

    pub fn get_workspace_cache(&self) -> Option<WorkspaceCacheData> {
        self.workspace_data.read()
            .ok()
            .map(|data| data.clone())
    }
}