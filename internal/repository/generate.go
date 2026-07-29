package repository

//go:generate go tool mockgen -source=interfaces.go -exclude_interfaces=LoginAttemptStore -destination=mock/repository.go -package=mock
