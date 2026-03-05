export interface Document {
  id: string;
  person_id: string;
  project_id: string;
  uploaded_by?: string;
  name: string;
  path: string;
  mime_type: string;
  size: number;
  created_at: string;
}

export interface CreateDocumentInput {
  person_id: string;
  name: string;
  path: string;
  mime_type: string;
  size: number;
}

export interface ListDocumentsOutput {
  documents: Document[];
  total: number;
}
