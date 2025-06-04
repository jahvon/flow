/**
 * This file was automatically generated from workspace_schema.json
 * DO NOT MODIFY IT BY HAND
 */

/**
 * A list of tags.
 * Tags can be used with list commands to filter returned data.
 *
 */
export type CommonTags = string[];

/**
 * Configuration for a workspace in the Flow CLI.
 * This configuration is used to define the settings for a workspace.
 * Every workspace has a workspace config file named `flow.yaml` in the root of the workspace directory.
 *
 */
export interface Workspace {
  /**
   * A description of the workspace. This description is rendered as markdown in the interactive UI.
   */
  description?: string;
  /**
   * A path to a markdown file that contains the description of the workspace.
   */
  descriptionFile?: string;
  /**
   * The display name of the workspace. This is used in the interactive UI.
   */
  displayName?: string;
  executables?: ExecutableFilter;
  tags?: CommonTags;
  [k: string]: unknown;
}
export interface ExecutableFilter {
  /**
   * A list of directories to exclude from the executable search.
   */
  excluded?: string[];
  /**
   * A list of directories to include in the executable search.
   */
  included?: string[];
  [k: string]: unknown;
}
