export interface Tag {
  id: string;
  project_id: string;
  name: string;
  created_at: string;
}

export interface CreateTagInput {
  name: string;
}

export interface ListTagsOutput {
  tags: Tag[];
}
