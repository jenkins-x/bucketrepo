package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
)

// Repository artifacts repository
type Repository interface {
	DownloadFile(path string) (io.ReadCloser, error)
}

// HttpRepository HTTP based artifacts repository
type HttpRepository struct {
	client  *http.Client
	baseURL string
}

var _ Repository = (*HttpRepository)(nil)

// NewRepository creates a new artifacts repository
func NewRepository(config RepositoryConfig) *HttpRepository {
	client := &http.Client{
		Timeout: config.Timeout,
	}
	return &HttpRepository{
		client:  client,
		baseURL: config.URL,
	}
}

// DownloadFile retrieves a file form the remote artifacts repository over HTTP
func (r *HttpRepository) DownloadFile(filePath string) (io.ReadCloser, error) {
	u, err := url.Parse(r.baseURL)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, filePath)

	resp, err := r.client.Get(u.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("status: %s", resp.Status)
	}
	return resp.Body, nil
}
