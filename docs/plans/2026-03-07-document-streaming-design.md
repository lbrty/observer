# Document Streaming + Thumbnails — Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add inline document streaming and on-demand image thumbnail generation to the document API.

**Architecture:** Three new endpoints under the existing project-scoped read group: `/stream` (inline), `/download` (attachment, already exists), `/thumbnail` (300px JPEG). Thumbnails are generated on first request and cached colocated with the original file. All endpoints share the same auth/permission middleware.

**Tech Stack:** Go stdlib `image`, `image/jpeg`, `image/png`, `image/gif` for decode; `golang.org/x/image/draw` for resizing; existing `storage.FileStorage` for persistence.

---

### Task 1: Add `golang.org/x/image` dependency

**Step 1: Add the dependency**

Run:
```bash
go get golang.org/x/image
```

**Step 2: Verify it's in go.mod**

Run:
```bash
grep 'golang.org/x/image' go.mod
```
Expected: a line with `golang.org/x/image`

**Step 3: Commit**

```bash
git add go.mod go.sum
git commit -m "add golang.org/x/image dependency for thumbnail resizing"
```

---

### Task 2: Add `ErrNotImage` domain error

**Files:**
- Modify: `internal/domain/document/errors.go`

**Step 1: Add error**

Add `ErrNotImage` to the existing var block in `internal/domain/document/errors.go`:

```go
var (
	ErrDocumentNotFound = errors.New("document not found")
	ErrNotImage         = errors.New("document is not an image")
)
```

**Step 2: Map error to HTTP 400 in handler**

In `internal/handler/errors.go`, add a case in `MapDomainError` right after the existing document error case (line ~125):

```go
case errors.Is(err, document.ErrNotImage):
	return http.StatusBadRequest, "errors.document.notImage"
```

**Step 3: Commit**

```bash
git add internal/domain/document/errors.go internal/handler/errors.go
git commit -m "add ErrNotImage domain error for thumbnail requests on non-image documents"
```

---

### Task 3: Add `Thumbnail` method to `DocumentUseCase` with tests

**Files:**
- Modify: `internal/usecase/project/document_usecase.go`
- Modify: `internal/usecase/project/document_usecase_test.go`

The `Thumbnail` method:
1. Fetches document metadata from repo
2. Checks if MIME type starts with `image/` — returns `document.ErrNotImage` if not
3. Builds the thumbnail path: same as original but with `_thumb.jpg` suffix appended
4. Tries `fs.Open` on the thumbnail path — if found, returns it (cache hit)
5. On cache miss: opens original, decodes image, resizes to 300px wide (aspect preserved), encodes as JPEG, saves via `fs.Save`, then opens and returns the thumbnail

**Step 1: Write the helper function `thumbnailPath`**

In `internal/usecase/project/document_usecase.go`, add after `storagePath`:

```go
func thumbnailPath(originalPath string) string {
	return originalPath + "_thumb.jpg"
}
```

**Step 2: Write `isImage` helper**

```go
func isImage(mimeType string) bool {
	return strings.HasPrefix(mimeType, "image/")
}
```

Add `"strings"` to imports.

**Step 3: Write `generateThumbnail` private method**

```go
const thumbnailWidth = 300

func (uc *DocumentUseCase) generateThumbnail(ctx context.Context, originalPath, thumbPath string) error {
	rc, err := uc.fs.Open(ctx, originalPath)
	if err != nil {
		return fmt.Errorf("open original: %w", err)
	}
	defer rc.Close()

	src, _, err := image.Decode(rc)
	if err != nil {
		return fmt.Errorf("decode image: %w", err)
	}

	bounds := src.Bounds()
	origW := bounds.Dx()
	origH := bounds.Dy()
	newW := thumbnailWidth
	newH := origH * newW / origW
	if newH < 1 {
		newH = 1
	}

	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
	draw.ApproxBiLinear.Scale(dst, dst.Bounds(), src, bounds, draw.Over, nil)

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, dst, &jpeg.Options{Quality: 80}); err != nil {
		return fmt.Errorf("encode thumbnail: %w", err)
	}

	if err := uc.fs.Save(ctx, thumbPath, &buf); err != nil {
		return fmt.Errorf("save thumbnail: %w", err)
	}
	return nil
}
```

Add these imports: `"bytes"`, `"image"`, `"image/jpeg"`, `_ "image/png"`, `_ "image/gif"`, `"golang.org/x/image/draw"`.

**Step 4: Write `Thumbnail` method**

```go
// Thumbnail returns a reader for the document's thumbnail image.
func (uc *DocumentUseCase) Thumbnail(ctx context.Context, id string) (*DocumentDTO, io.ReadCloser, error) {
	d, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, nil, fmt.Errorf("get document: %w", err)
	}
	if !isImage(d.MimeType) {
		return nil, nil, document.ErrNotImage
	}

	thumbPath := thumbnailPath(d.Path)
	rc, err := uc.fs.Open(ctx, thumbPath)
	if err == nil {
		dto := documentToDTO(d)
		return &dto, rc, nil
	}

	if err := uc.generateThumbnail(ctx, d.Path, thumbPath); err != nil {
		return nil, nil, fmt.Errorf("generate thumbnail: %w", err)
	}

	rc, err = uc.fs.Open(ctx, thumbPath)
	if err != nil {
		return nil, nil, fmt.Errorf("open thumbnail: %w", err)
	}
	dto := documentToDTO(d)
	return &dto, rc, nil
}
```

**Step 5: Write tests**

In `internal/usecase/project/document_usecase_test.go`, add:

```go
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
```

Add `"io"`, `"strings"` to test imports.

**Step 6: Run tests**

Run:
```bash
just test
```
Expected: all pass

**Step 7: Commit**

```bash
git add internal/usecase/project/document_usecase.go internal/usecase/project/document_usecase_test.go
git commit -m "add Thumbnail use case method with cache-hit and non-image tests"
```

---

### Task 4: Add `Stream` and `Thumbnail` handler methods

**Files:**
- Modify: `internal/handler/document_handler.go`

**Step 1: Add `Stream` method**

```go
// Stream handles GET /projects/:project_id/documents/:id/stream.
func (h *DocumentHandler) Stream(c *gin.Context) {
	if !middleware.CanViewDocumentsFrom(c) {
		c.JSON(http.StatusForbidden, errJSON("errors.document.insufficientPermissions", "insufficient permissions to view documents"))
		return
	}

	doc, rc, err := h.uc.Download(c.Request.Context(), c.Param("id"))
	if err != nil {
		HandleError(c, err)
		return
	}
	defer rc.Close()

	c.Header("Content-Disposition", "inline; filename=\""+doc.Name+"\"")
	c.DataFromReader(http.StatusOK, doc.Size, doc.MimeType, rc, nil)
}
```

**Step 2: Add `Thumbnail` method**

```go
// Thumbnail handles GET /projects/:project_id/documents/:id/thumbnail.
func (h *DocumentHandler) Thumbnail(c *gin.Context) {
	if !middleware.CanViewDocumentsFrom(c) {
		c.JSON(http.StatusForbidden, errJSON("errors.document.insufficientPermissions", "insufficient permissions to view documents"))
		return
	}

	_, rc, err := h.uc.Thumbnail(c.Request.Context(), c.Param("id"))
	if err != nil {
		HandleError(c, err)
		return
	}
	defer rc.Close()

	c.Header("Cache-Control", "public, max-age=86400")
	c.DataFromReader(http.StatusOK, -1, "image/jpeg", rc, nil)
}
```

**Step 3: Commit**

```bash
git add internal/handler/document_handler.go
git commit -m "add Stream and Thumbnail handler methods"
```

---

### Task 5: Register new routes

**Files:**
- Modify: `internal/server/server.go` (~line 232)

**Step 1: Add routes**

After the existing line `read.GET("/documents/:id/download", documentHandler.Download)` (line 232), add:

```go
read.GET("/documents/:id/stream", documentHandler.Stream)
read.GET("/documents/:id/thumbnail", documentHandler.Thumbnail)
```

**Step 2: Verify build**

Run:
```bash
go build ./...
```
Expected: no errors

**Step 3: Run all tests**

Run:
```bash
just test
```
Expected: all pass

**Step 4: Commit**

```bash
git add internal/server/server.go
git commit -m "register /stream and /thumbnail document routes"
```

---

### Summary of changes

| File | Change |
|------|--------|
| `go.mod` / `go.sum` | Add `golang.org/x/image` |
| `internal/domain/document/errors.go` | Add `ErrNotImage` |
| `internal/handler/errors.go` | Map `ErrNotImage` → 400 |
| `internal/usecase/project/document_usecase.go` | Add `thumbnailPath`, `isImage`, `generateThumbnail`, `Thumbnail` |
| `internal/usecase/project/document_usecase_test.go` | Add 3 thumbnail tests |
| `internal/handler/document_handler.go` | Add `Stream`, `Thumbnail` methods |
| `internal/server/server.go` | Register 2 new routes |
