package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
)

type Repository interface {
	DownloadFile(path string) (io.ReadCloser, error)
}

type HttpRepository struct {
	client  *http.Client
	baseURL string
}

var _ Repository = (*HttpRepository)(nil)

func NewRepository(config RepositoryConfig) *HttpRepository {
	client := &http.Client{
		Timeout: config.Timeout,
	}
	return &HttpRepository{
		client:  client,
		baseURL: config.URL,
	}
}

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
