import { Workspace } from "./generated/workspace";

export interface EnrichedWorkspace extends Workspace {
  name: string;
  path: string;
  fullDescription: string;
}
