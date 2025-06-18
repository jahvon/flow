use crate::types::generated::flowfile::Executable as GeneratedExecutable;
use crate::types::generated::workspace::Workspace as GeneratedWorkspace;

#[derive(Debug, Clone, serde::Serialize, serde::Deserialize)]
pub struct Workspace {
    #[serde(flatten)]
    pub base: GeneratedWorkspace,
    pub id: String,
    pub path: String,
    #[serde(rename = "fullDescription")]
    pub full_description: String,
}

impl Workspace {
    pub fn new(
        base: GeneratedWorkspace,
        id: String,
        path: String,
        full_description: String,
    ) -> Self {
        Self {
            base,
            id,
            path,
            full_description,
        }
    }
}

impl std::ops::Deref for Workspace {
    type Target = GeneratedWorkspace;

    fn deref(&self) -> &Self::Target {
        &self.base
    }
}

impl std::ops::DerefMut for Workspace {
    fn deref_mut(&mut self) -> &mut Self::Target {
        &mut self.base
    }
}

#[derive(Debug, Clone, serde::Serialize, serde::Deserialize)]
pub struct Executable {
    #[serde(flatten)]
    pub base: GeneratedExecutable,
    pub id: String,
    #[serde(rename = "ref")]
    pub ref_: String,
    pub namespace: Option<String>,
    pub workspace: String,
    pub flowfile: String,
    #[serde(rename = "fullDescription")]
    pub full_description: String,
}

impl Executable {
    pub fn new(
        base: GeneratedExecutable,
        id: String,
        ref_: String,
        namespace: Option<String>,
        workspace: String,
        flowfile: String,
        full_description: String,
    ) -> Self {
        Self {
            base,
            id,
            ref_,
            namespace,
            workspace,
            flowfile,
            full_description,
        }
    }
}

impl std::ops::Deref for Executable {
    type Target = GeneratedExecutable;

    fn deref(&self) -> &Self::Target {
        &self.base
    }
}

impl std::ops::DerefMut for Executable {
    fn deref_mut(&mut self) -> &mut Self::Target {
        &mut self.base
    }
}
