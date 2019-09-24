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
	BaseURL() string
}

// HTTPRepository HTTP based artifacts repository
type HTTPRepository struct {
	client  *http.Client
	baseURL string
	header  http.Header
}

var _ Repository = (*HTTPRepository)(nil)

// NewRepository creates a new artifacts repository
func NewRepository(config RepositoryConfig) *HTTPRepository {
	client := &http.Client{
		Timeout: config.Timeout,
	}
	return &HTTPRepository{
		client:  client,
		baseURL: config.URL,
		header:  config.Header,
	}
}

// DownloadFile retrieves a file form the remote artifacts repository over HTTP
func (r *HTTPRepository) DownloadFile(filePath string) (io.ReadCloser, error) {
	u, err := url.Parse(r.baseURL)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, filePath)

	request, err := http.NewRequest(http.MethodGet, u.String(), nil)
	for key, values := range r.header {
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}
	resp, err := r.client.Do(request)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		err := resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("status: %s, closing body err: %s", resp.Status, err)
		}
		return nil, fmt.Errorf("status: %s", resp.Status)
	}
	return resp.Body, nil
}

// BaseURL returns the base URL of the repository
func (r *HTTPRepository) BaseURL() string {
	return r.baseURL
}
