use std::process::Command;
use std::fs;

fn main() {
    println!("cargo:rerun-if-changed=../../docs/schemas/");

    // Create output directory
    let output_dir = "src/generated";
    fs::create_dir_all(output_dir).unwrap();

    // Schema files to process
    let schemas = [
        ("flowfile", "../../docs/schemas/flowfile_schema.json"),
        ("workspace", "../../docs/schemas/workspace_schema.json"),
        ("template", "../../docs/schemas/template_schema.json"),
        ("config", "../../docs/schemas/config_schema.json"),
    ];

    for (name, schema_path) in schemas {
        if std::path::Path::new(schema_path).exists() {
            let output_file = format!("{}/{}.rs", output_dir, name);

            let result = Command::new("cargo")
                .args(&["typify", schema_path, "--output", &output_file])
                .output();

            match result {
                Ok(output) if output.status.success() => {
                    println!("Generated types for {} at {}", name, output_file);
                }
                Ok(output) => {
                    println!("cargo:warning=Failed to generate {}: {}", name,
                             String::from_utf8_lossy(&output.stderr));
                }
                Err(e) => {
                    println!("cargo:warning=Failed to run cargo typify for {}: {}", name, e);
                }
            }
        } else {
            println!("cargo:warning=Schema file not found: {}", schema_path);
        }
    }

    // Generate mod.rs
    let mod_content = format!(
        "// Generated module exports\n{}\n",
        schemas.iter()
            .filter(|(_, path)| std::path::Path::new(path).exists())
            .map(|(name, _)| format!("pub mod {};", name))
            .collect::<Vec<_>>()
            .join("\n")
    );

    fs::write(format!("{}/mod.rs", output_dir), mod_content).unwrap();

    // Keep the existing Tauri build process
    tauri_build::build();
}