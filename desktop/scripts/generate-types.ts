import { compile } from 'json-schema-to-typescript';
import fs from 'fs/promises';
import path from 'path';

const SCHEMAS_DIR = path.resolve(process.cwd(), '../docs/schemas');
const OUTPUT_DIR = path.resolve(process.cwd(), 'src/types');

async function generateTypes() {
  try {
    // Ensure output directory exists
    await fs.mkdir(OUTPUT_DIR, { recursive: true });

    // Read all schema files
    const schemaFiles = await fs.readdir(SCHEMAS_DIR);

    for (const file of schemaFiles) {
      if (!file.endsWith('.json')) continue;

      const schemaPath = path.join(SCHEMAS_DIR, file);
      const schema = JSON.parse(await fs.readFile(schemaPath, 'utf-8'));

      // Generate TypeScript types
      const types = await compile(schema, file.replace('.json', ''), {
        bannerComment: `/**
 * This file was automatically generated from ${file}
 * DO NOT MODIFY IT BY HAND
 */`,
        style: {
          singleQuote: true,
          semi: true,
        },
      });

      // Write the generated types to a file
      const outputPath = path.join(OUTPUT_DIR, `${file.replace('.json', '').replace('_schema', '')}.ts`);
      await fs.writeFile(outputPath, types);

      console.log(`Generated types for ${file}`);
    }

    console.log('Type generation complete!');
  } catch (error) {
    console.error('Error generating types:', error);
    process.exit(1);
  }
}

generateTypes();