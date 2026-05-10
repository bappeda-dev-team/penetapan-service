package api

import (
	"log/slog"
)

const Version = "1.0.0"

// port: which port of server should be listening to
// env: which env should used, prod, staging, and dev
type Config struct {
	Port int
	Env  string
}

// hold application dependencies for
// http handler, middleware and helpers
type Application struct {
	Config Config
	Logger *slog.Logger

	// inject service
	PenetapanOpdService *PenetapanOpdService
}
