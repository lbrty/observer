package project_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/lbrty/observer/internal/domain/document"
	mock_repo "github.com/lbrty/observer/internal/repository/mock"
	ucproject "github.com/lbrty/observer/internal/usecase/project"
)

func TestDocumentUseCase_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockDocumentRepository(ctrl)
	uc := ucproject.NewDocumentUseCase(mockRepo)

	existing := &document.Document{
		ID:        "d1",
		PersonID:  "p1",
		ProjectID: "proj1",
		Name:      "old-name.pdf",
		Path:      "/docs/old-name.pdf",
		MimeType:  "application/pdf",
		Size:      1024,
		CreatedAt: time.Now().UTC(),
	}

	mockRepo.EXPECT().GetByID(gomock.Any(), "d1").Return(existing, nil)
	mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, d *document.Document) error {
		assert.Equal(t, "new-name.pdf", d.Name)
		assert.Equal(t, "/docs/old-name.pdf", d.Path, "path should not change")
		return nil
	})

	newName := "new-name.pdf"
	out, err := uc.Update(context.Background(), "d1", ucproject.UpdateDocumentInput{Name: &newName})
	require.NoError(t, err)
	assert.Equal(t, "d1", out.ID)
	assert.Equal(t, "new-name.pdf", out.Name)
}

func TestDocumentUseCase_Update_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockDocumentRepository(ctrl)
	uc := ucproject.NewDocumentUseCase(mockRepo)

	mockRepo.EXPECT().GetByID(gomock.Any(), "d1").Return(nil, errors.New("not found"))

	_, err := uc.Update(context.Background(), "d1", ucproject.UpdateDocumentInput{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "get document for update")
}

func TestDocumentUseCase_Update_NilName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockDocumentRepository(ctrl)
	uc := ucproject.NewDocumentUseCase(mockRepo)

	existing := &document.Document{
		ID:        "d1",
		PersonID:  "p1",
		ProjectID: "proj1",
		Name:      "unchanged.pdf",
		Path:      "/docs/unchanged.pdf",
		MimeType:  "application/pdf",
		Size:      512,
		CreatedAt: time.Now().UTC(),
	}

	mockRepo.EXPECT().GetByID(gomock.Any(), "d1").Return(existing, nil)
	mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, d *document.Document) error {
		assert.Equal(t, "unchanged.pdf", d.Name, "name should remain unchanged")
		return nil
	})

	out, err := uc.Update(context.Background(), "d1", ucproject.UpdateDocumentInput{})
	require.NoError(t, err)
	assert.Equal(t, "unchanged.pdf", out.Name)
}
