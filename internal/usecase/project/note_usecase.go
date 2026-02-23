package project

import (
	"context"
	"fmt"

	"github.com/lbrty/observer/internal/domain/note"
	"github.com/lbrty/observer/internal/repository"
	"github.com/lbrty/observer/internal/ulid"
)

// NoteUseCase handles person note operations (append-only minus delete).
type NoteUseCase struct {
	repo repository.PersonNoteRepository
}

// NewNoteUseCase creates a NoteUseCase.
func NewNoteUseCase(repo repository.PersonNoteRepository) *NoteUseCase {
	return &NoteUseCase{repo: repo}
}

// List returns all notes for a person.
func (uc *NoteUseCase) List(ctx context.Context, personID string) ([]NoteDTO, error) {
	notes, err := uc.repo.List(ctx, personID)
	if err != nil {
		return nil, fmt.Errorf("list notes: %w", err)
	}
	dtos := make([]NoteDTO, len(notes))
	for i, n := range notes {
		dtos[i] = noteToDTO(n)
	}
	return dtos, nil
}

// Create creates a new note. authorID is auto-set from auth context.
func (uc *NoteUseCase) Create(ctx context.Context, personID, authorID string, input CreateNoteInput) (*NoteDTO, error) {
	n := &note.Note{
		ID:       ulid.NewString(),
		PersonID: personID,
		AuthorID: &authorID,
		Body:     input.Body,
	}
	if err := uc.repo.Create(ctx, n); err != nil {
		return nil, fmt.Errorf("create note: %w", err)
	}
	dto := noteToDTO(n)
	return &dto, nil
}

// Delete removes a note.
func (uc *NoteUseCase) Delete(ctx context.Context, id string) error {
	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete note: %w", err)
	}
	return nil
}
