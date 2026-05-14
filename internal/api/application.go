package api

import (
	"log/slog"

	"github.com/bappeda-dev-team/penetapan-service/internal/service/opd"
)

const Version = "1.0.0"

// port: which port of server should be listening to
// env: which env should used, prod, staging, and dev
type Config struct {
	Port int
	Env  string

	DB struct {
		Dsn string
	}

	Services struct {
		Perencanaan ServiceConfig
	}
}

// ServiceConfig used for external services
type ServiceConfig struct {
	BaseURL string
	ApiPath string
}

// hold application dependencies for
// http handler, middleware and helpers
type Application struct {
	Config Config
	Logger *slog.Logger

	// inject service
	PenetapanOpdService *opd.PenetapanOpdService
}
