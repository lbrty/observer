package project

import (
	"context"
	"fmt"

	"github.com/lbrty/observer/internal/domain/note"
	"github.com/lbrty/observer/internal/repository"
	"github.com/lbrty/observer/internal/ulid"
	ucaudit "github.com/lbrty/observer/internal/usecase/audit"
)

// NoteUseCase handles person note operations (append-only minus delete).
type NoteUseCase struct {
	repo    repository.PersonNoteRepository
	auditUC *ucaudit.AuditUseCase
}

// NewNoteUseCase creates a NoteUseCase.
func NewNoteUseCase(repo repository.PersonNoteRepository, auditUC *ucaudit.AuditUseCase) *NoteUseCase {
	return &NoteUseCase{repo: repo, auditUC: auditUC}
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
func (uc *NoteUseCase) Create(ctx context.Context, projectID, personID, authorID string, input CreateNoteInput) (*NoteDTO, error) {
	n := &note.Note{
		ID:       ulid.NewString(),
		PersonID: personID,
		AuthorID: &authorID,
		Body:     input.Body,
	}
	if err := uc.repo.Create(ctx, n); err != nil {
		return nil, fmt.Errorf("create note: %w", err)
	}
	uc.auditUC.Record(ctx, &projectID, "note.create", "note", &n.ID, fmt.Sprintf("Created note %s", n.ID))
	dto := noteToDTO(n)
	return &dto, nil
}

// Update updates a note body.
func (uc *NoteUseCase) Update(ctx context.Context, id string, input UpdateNoteInput) (*NoteDTO, error) {
	n, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get note for update: %w", err)
	}
	n.Body = input.Body
	if err := uc.repo.Update(ctx, n); err != nil {
		return nil, fmt.Errorf("update note: %w", err)
	}
	dto := noteToDTO(n)
	return &dto, nil
}

// Delete removes a note.
func (uc *NoteUseCase) Delete(ctx context.Context, projectID, id string) error {
	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete note: %w", err)
	}
	uc.auditUC.Record(ctx, &projectID, "note.delete", "note", &id, fmt.Sprintf("Deleted note %s", id))
	return nil
}
