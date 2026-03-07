package project_test

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/lbrty/observer/internal/domain/document"
	mock_repo "github.com/lbrty/observer/internal/repository/mock"
	mock_storage "github.com/lbrty/observer/internal/storage/mock"
	ucproject "github.com/lbrty/observer/internal/usecase/project"
)

func TestDocumentUseCase_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockDocumentRepository(ctrl)
	mockFS := mock_storage.NewMockFileStorage(ctrl)
	uc := ucproject.NewDocumentUseCase(mockRepo, mockFS)

	existing := &document.Document{
		ID:        "d1",
		PersonID:  "p1",
		ProjectID: "proj1",
		Name:      "old-name.pdf",
		Path:      "proj1/p1/d1_old-name.pdf",
		MimeType:  "application/pdf",
		Size:      1024,
		CreatedAt: time.Now().UTC(),
	}

	mockRepo.EXPECT().GetByID(gomock.Any(), "d1").Return(existing, nil)
	mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, d *document.Document) error {
		assert.Equal(t, "new-name.pdf", d.Name)
		assert.Equal(t, "proj1/p1/d1_old-name.pdf", d.Path, "path should not change")
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
	mockFS := mock_storage.NewMockFileStorage(ctrl)
	uc := ucproject.NewDocumentUseCase(mockRepo, mockFS)

	mockRepo.EXPECT().GetByID(gomock.Any(), "d1").Return(nil, errors.New("not found"))

	_, err := uc.Update(context.Background(), "d1", ucproject.UpdateDocumentInput{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "get document for update")
}

func TestDocumentUseCase_Update_NilName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockDocumentRepository(ctrl)
	mockFS := mock_storage.NewMockFileStorage(ctrl)
	uc := ucproject.NewDocumentUseCase(mockRepo, mockFS)

	existing := &document.Document{
		ID:        "d1",
		PersonID:  "p1",
		ProjectID: "proj1",
		Name:      "unchanged.pdf",
		Path:      "proj1/p1/d1_unchanged.pdf",
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

func TestDocumentUseCase_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockDocumentRepository(ctrl)
	mockFS := mock_storage.NewMockFileStorage(ctrl)
	uc := ucproject.NewDocumentUseCase(mockRepo, mockFS)

	existing := &document.Document{
		ID:   "d1",
		Path: "proj1/p1/d1_test.pdf",
	}

	mockRepo.EXPECT().GetByID(gomock.Any(), "d1").Return(existing, nil)
	mockRepo.EXPECT().Delete(gomock.Any(), "d1").Return(nil)
	mockFS.EXPECT().Delete(gomock.Any(), "proj1/p1/d1_test.pdf").Return(nil)

	err := uc.Delete(context.Background(), "d1")
	require.NoError(t, err)
}

func TestDocumentUseCase_Thumbnail_NotImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockDocumentRepository(ctrl)
	mockFS := mock_storage.NewMockFileStorage(ctrl)
	uc := ucproject.NewDocumentUseCase(mockRepo, mockFS)

	mockRepo.EXPECT().GetByID(gomock.Any(), "d1").Return(&document.Document{
		ID:       "d1",
		Path:     "proj1/p1/d1_test.pdf",
		MimeType: "application/pdf",
	}, nil)

	_, _, err := uc.Thumbnail(context.Background(), "d1")
	require.ErrorIs(t, err, document.ErrNotImage)
}

func TestDocumentUseCase_Thumbnail_CacheHit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockDocumentRepository(ctrl)
	mockFS := mock_storage.NewMockFileStorage(ctrl)
	uc := ucproject.NewDocumentUseCase(mockRepo, mockFS)

	mockRepo.EXPECT().GetByID(gomock.Any(), "d1").Return(&document.Document{
		ID:        "d1",
		PersonID:  "p1",
		ProjectID: "proj1",
		Path:      "proj1/p1/d1_photo.jpg",
		MimeType:  "image/jpeg",
		Name:      "photo.jpg",
		Size:      5000,
	}, nil)

	thumbReader := io.NopCloser(strings.NewReader("thumb-data"))
	mockFS.EXPECT().Open(gomock.Any(), "proj1/p1/d1_photo.jpg_thumb.jpg").Return(thumbReader, nil)

	dto, rc, err := uc.Thumbnail(context.Background(), "d1")
	require.NoError(t, err)
	defer rc.Close()
	assert.Equal(t, "d1", dto.ID)
}

func TestDocumentUseCase_Thumbnail_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockDocumentRepository(ctrl)
	mockFS := mock_storage.NewMockFileStorage(ctrl)
	uc := ucproject.NewDocumentUseCase(mockRepo, mockFS)

	mockRepo.EXPECT().GetByID(gomock.Any(), "d1").Return(nil, document.ErrDocumentNotFound)

	_, _, err := uc.Thumbnail(context.Background(), "d1")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "get document")
}
