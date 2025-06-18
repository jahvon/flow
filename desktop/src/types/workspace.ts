import { Workspace } from "./generated/workspace";

export interface EnrichedWorkspace extends Workspace {
  id: string;
  path: string;
  fullDescription: string;
}
