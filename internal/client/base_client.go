package client

import (
	"context"
	"net/http"
	"strings"
)

type BaseClient struct {
	BaseURL    string
	HttpClient *http.Client
}

// constructor
func NewBaseClient(host, apiPath string, httpClient *http.Client) *BaseClient {
	baseURL := strings.TrimRight(host, "/")
	if apiPath != "" {
		baseURL += "/" + strings.Trim(apiPath, "/")
	}
	return &BaseClient{
		BaseURL:    baseURL,
		HttpClient: httpClient,
	}
}

// key untuk context session id
type ctxKey string

const SessionIDKey ctxKey = "X-Session-Id"

// Inject session ID ke context (opsional)
func WithSessionID(ctx context.Context, sessionID string) context.Context {
	return context.WithValue(ctx, SessionIDKey, sessionID)
}

// Ambil session ID dari context
func GetSessionID(ctx context.Context) string {
	if v := ctx.Value(SessionIDKey); v != nil {
		if s, ok := v.(string); ok && s != "" {
			return s
		}
	}
	return ""
}
