package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/bappeda-dev-team/penetapan-service/internal/client"
)

type Client struct {
	*client.BaseClient
}

func NewClient(host, apiPath string, httpClient *http.Client) *Client {
	return &Client{
		BaseClient: client.NewBaseClient(host, apiPath, httpClient),
	}
}

func (c *Client) SyncPkPenetapan(
	ctx context.Context,
	request SyncRequest,
) ([]PkPenetapanResponse, error) {
	// url
	baseUrl := fmt.Sprintf("%s/pk/penetapan", c.BaseURL)
	params := url.Values{}
	params.Add("id_pegawai", request.PegawaiId)
	params.Add("kode_opd", request.KodeOpd)
	params.Add("tahun", strconv.Itoa(request.Tahun))

	// request
	req, err := http.NewRequest(http.MethodGet, baseUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("Gagal membuat request: %w", err)
	}
	req.URL.RawQuery = params.Encode()
	req.Header.Set("Content-Type", "application/json")

	// set session id
	sessionID := client.GetSessionID(ctx)
	if sessionID != "" {
		req.Header.Set("X-Session-Id", sessionID)
	} else {
		log.Printf("Tidak ada X-Session-Id ditemukan, mungkin akan 401")
	}

	// send request
	res, err := c.HttpClient.Do(req)
	if err != nil {
		log.Printf("ERROR SYNC TO PERENCANAAN: %v", err)
		return nil, fmt.Errorf("Request gagal: %w", err)
	}
	defer res.Body.Close()

	// response status
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Request ke perencanaan individu gagal. status: %d", res.StatusCode)
	}

	type wrapper struct {
		Code   int                   `json:"code"`
		Status string                `json:"status"`
		Data   []PkPenetapanResponse `json:"data"`
	}

	var result wrapper
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("gagal decode response: %w", err)
	}

	return result.Data, nil
}

func (c *Client) SyncRenjaIndividu(
	ctx context.Context,
	request SyncRequest,
) ([]RenjaIndividuResponse, error) {
	// url
	baseUrl := fmt.Sprintf("%s/pk/penetapan/renja", c.BaseURL)
	params := url.Values{}
	params.Add("id_pegawai", request.PegawaiId)
	params.Add("kode_opd", request.KodeOpd)
	params.Add("tahun", strconv.Itoa(request.Tahun))

	// request
	req, err := http.NewRequest(http.MethodGet, baseUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("Gagal membuat request: %w", err)
	}
	req.URL.RawQuery = params.Encode()
	req.Header.Set("Content-Type", "application/json")

	// set session id
	sessionID := client.GetSessionID(ctx)
	if sessionID != "" {
		req.Header.Set("X-Session-Id", sessionID)
	} else {
		log.Printf("Tidak ada X-Session-Id ditemukan, mungkin akan 401")
	}

	// send request
	res, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Request gagal: %w", err)
	}
	defer res.Body.Close()

	// response status
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Request ke pk penetapan gagal. status: %d", res.StatusCode)
	}

	type wrapper struct {
		Code   int                     `json:"code"`
		Status string                  `json:"status"`
		Data   []RenjaIndividuResponse `json:"data"`
	}

	var result wrapper
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("gagal decode response: %w", err)
	}

	return result.Data, nil
}
