package perencanaan

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/bappeda-dev-team/penetapan-service/internal/client"
)

type PerencanaanClient struct {
	*client.BaseClient
}

func NewPerencanaanClient(host, apiPath string, httpClient *http.Client) *PerencanaanClient {
	return &PerencanaanClient{
		BaseClient: client.NewBaseClient(host, apiPath, httpClient),
	}
}

func (c *PerencanaanClient) GetPenetapanTujuanOpd(ctx context.Context, request PerencanaanRequest) ([]PerencanaanTujuanOpdResponse, error) {
	// url penetapan tujuan opd
	url := fmt.Sprintf("%s/tujuan_opd/penetapan/%s/%d", c.BaseURL, request.KodeOpd, request.Tahun)

	// request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("Gagal membuat request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	sessionID := client.GetSessionID(ctx)
	if sessionID != "" {
		req.Header.Set("X-Session-Id", sessionID)
	} else {
		log.Printf("Tidak ada X-Session-Id ditemukan, mungkin akan 401")
	}

	// send request
	res, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Request ke perencanaan gagal: %w", err)
	}
	defer res.Body.Close()

	// response status
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Tujuan Opd Penetapan: gagal sync. status: %d", res.StatusCode)
	}

	// safe, response pasti ada
	type wrapper struct {
		Code   int                            `json:"code"`
		Status string                         `json:"status"`
		Data   []PerencanaanTujuanOpdResponse `json:"data"`
	}

	var result wrapper
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("gagal decode response: %w", err)
	}

	return result.Data, nil
}

func (c *PerencanaanClient) GetPenetapanSasaranOpd(ctx context.Context, request PerencanaanRequest) ([]PerencanaanSasaranOpdResponse, error) {
	// url penetapan tujuan opd
	url := fmt.Sprintf("%s/sasaran_opd/penetapan/%s/%d", c.BaseURL, request.KodeOpd, request.Tahun)

	// request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("Gagal membuat request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	sessionID := client.GetSessionID(ctx)
	if sessionID != "" {
		req.Header.Set("X-Session-Id", sessionID)
	} else {
		log.Printf("Tidak ada X-Session-Id ditemukan, mungkin akan 401")
	}

	// send request
	res, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Request ke perencanaan gagal: %w", err)
	}
	defer res.Body.Close()

	// response status
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Sasaran Opd Penetapan: gagal sync. status: %d", res.StatusCode)
	}

	// safe, response pasti ada
	type wrapper struct {
		Code   int                             `json:"code"`
		Status string                          `json:"status"`
		Data   []PerencanaanSasaranOpdResponse `json:"data"`
	}

	var result wrapper
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("gagal decode response: %w", err)
	}

	return result.Data, nil
}
