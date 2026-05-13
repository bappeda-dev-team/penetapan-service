package main

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/bappeda-dev-team/penetapan-service/internal/api"
	"github.com/joho/godotenv"
)

// loadConfig digunakan untuk mengambil konfigurasi dari env
// APP_ENV untuk menunjukkan dimana aplikasi berjalan
// APP_PORT untuk port yang digunakan aplikasi
// DB_DSN untuk connection string ke database
func loadConfig(logger *slog.Logger) api.Config {

	_ = godotenv.Load()

	var cfg api.Config

	cfg.Env = getEnv(
		"APP_ENV",
		"development",
	)

	cfg.Port = getIntEnv(
		"APP_PORT",
		4000,
	)

	cfg.DB.Dsn = getRequiredEnv(
		logger,
		"DB_DSN",
	)

	cfg.Services.Perencanaan.BaseURL = getEnv(
		"SERVICES_PERENCANAAN_BASE_URL",
		"http://localhost:8080",
	)

	cfg.Services.Perencanaan.ApiPath = getEnv(
		"SERVICES_PERENCANAAN_API_PATH",
		"",
	)

	return cfg
}

// getRequiredEnv adalah fungsi untuk mengambil env wajib
// tanpa default value
// crash jika key env tidak ditemukan
func getRequiredEnv(
	logger *slog.Logger,
	key string,
) string {

	value := os.Getenv(key)

	if value == "" {

		logger.Error(
			key + " environment variable is required",
		)

		os.Exit(1)
	}

	return value
}

// getEnv adalah fungsi untuk mengambil env string
// dengan default value
func getEnv(
	key string,
	defaultValue string,
) string {

	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}

// getIntEnv adalah fungsi untuk mengambil env integer (angka)
// dengan default value
// crash jika env value bukan angka
func getIntEnv(
	key string,
	defaultValue int,
) int {

	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	result, err := strconv.Atoi(value)

	if err != nil {
		return defaultValue
	}

	return result
}
