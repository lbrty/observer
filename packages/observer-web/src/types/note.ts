export interface Note {
  id: string;
  person_id: string;
  author_id?: string;
  body: string;
  created_at: string;
  updated_at?: string;
}

export interface CreateNoteInput {
  body: string;
}

export interface UpdateNoteInput {
  body: string;
}

export interface ListNotesOutput {
  notes: Note[];
}
