package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"
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
	if err != nil {
		return nil, err
	}
	for key, values := range r.header {
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}
	resp, err := r.client.Do(request)
	l := log.WithField("requestUrl", u).
		WithField("requestHeaders", r.header)
	if err != nil {
		l.WithError(err).Debugf("Failed to download file")
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		if log.IsLevelEnabled(log.DebugLevel) {
			if resp.Body != nil {
				var bout strings.Builder
				_, err = io.Copy(&bout, resp.Body)
				if err != nil {
					l = l.WithError(err)
				} else {
					l = l.WithField("responseBody", bout.String())
				}
			}
			l.WithField("responseStatus", resp.Status).Debugf("Failed to download file")
		}
		err = resp.Body.Close()
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
