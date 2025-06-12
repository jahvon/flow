use crate::types::generated::flowfile::Executable;

#[derive(Debug, Clone, serde::Serialize, serde::Deserialize)]
pub struct EnrichedExecutable {
    #[serde(flatten)]
    pub base: Executable,
    pub id: String,
    pub ref_: String,
    pub namespace: Option<String>,
    pub workspace: String,
    pub flowfile: String,
}

impl EnrichedExecutable {
    pub fn new(
        base: Executable,
        id: String,
        ref_: String,
        namespace: Option<String>,
        workspace: String,
        flowfile: String,
    ) -> Self {
        Self {
            base,
            id,
            ref_,
            namespace,
            workspace,
            flowfile,
        }
    }
}

impl std::ops::Deref for EnrichedExecutable {
    type Target = Executable;

    fn deref(&self) -> &Self::Target {
        &self.base
    }
}

impl std::ops::DerefMut for EnrichedExecutable {
    fn deref_mut(&mut self) -> &mut Self::Target {
        &mut self.base
    }
}
