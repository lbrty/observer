# ADR-010: S3-Compatible Object Storage

| Field      | Value             |
| ---------- | ----------------- |
| Status     | Proposed          |
| Date       | 2026-03-13        |
| Supersedes | —                 |
| Components | observer, storage |


## Decision

Add an `S3Storage` implementation of `FileStorage` backed by `aws-sdk-go-v2`, selectable via `STORAGE_BACKEND`. `local` remains the default so existing deployments require no changes.


## Why aws-sdk-go-v2

- De-facto standard S3 client in the Go ecosystem, actively maintained by AWS.
- Works against any S3-compatible endpoint (AWS S3, MinIO, Wasabi, Backblaze B2) via a configurable endpoint override; leaving it empty uses the AWS default resolver.
- Maps cleanly to the `FileStorage` interface (`Save`, `Open`, `Delete`).
- Verbosity is contained inside `internal/storage/s3.go`; the interface and all consumers are unchanged.
- Supports the SDK default credential provider chain, so AWS-native deployments do not need to duplicate credentials.


## Implementation

### New file: `internal/storage/s3.go`

| Interface method | S3 operation   |
| ---------------- | -------------- |
| `Save`           | `PutObject`    |
| `Open`           | `GetObject`    |
| `Delete`         | `DeleteObject` |

`Delete` on a missing key returns `nil` — AWS S3 `DeleteObject` is unconditionally idempotent, naturally matching `LocalStorage`'s `os.IsNotExist` guard.

Object keys are the `path` argument as-is (e.g. `projects/<id>/documents/file.pdf`).

`NewS3Storage` validates that `Bucket` is non-empty and returns an error immediately, rather than letting a confusing SDK error surface at the first operation.

`AccessKey` and `SecretKey` are optional. When both are empty, `NewS3Storage` uses the SDK default credential provider chain (env vars → shared credentials file → IAM instance role → ECS task role).

### Config additions (`internal/config/config.go`)

`StorageConfig` is replaced with:

```go
type StorageConfig struct {
	Path        string // STORAGE_PATH,    default "data/uploads"
	Backend     string // STORAGE_BACKEND, default "local"
	S3Endpoint  string // S3_ENDPOINT,     default "" (AWS default)
	S3Bucket    string // S3_BUCKET
	S3Region    string // S3_REGION,       default "us-east-1"
	S3AccessKey string // S3_ACCESS_KEY    (optional — falls back to SDK chain)
	S3SecretKey string // S3_SECRET_KEY    (optional — falls back to SDK chain)
}
```

### DI wiring (`internal/app/container.go`)

```go
var (
	fileStorage storage.FileStorage
	err         error
)
switch cfg.Storage.Backend {
case "s3":
	fileStorage, err = storage.NewS3Storage(cfg.Storage)
default:
	fileStorage, err = storage.NewLocalStorage(cfg.Storage.Path)
}
if err != nil {
	return nil, fmt.Errorf("init file storage: %w", err)
}
```


## Alternatives Considered

### A. minio-go

S3-compatible client maintained by MinIO. Functional against AWS S3 and MinIO, but optimised for MinIO-specific features this project does not need. `aws-sdk-go-v2` is the more widely recognised standard for S3 in Go and has broader ecosystem support.

### B. gocloud.dev/blob

Generic blob abstraction over S3, GCS, and Azure. Uses `aws-sdk-go-v2` under the hood, so it adds an abstraction layer with no benefit for an S3-only target. Rejected as unnecessary complexity.
