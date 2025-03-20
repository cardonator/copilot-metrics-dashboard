package services

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go.uber.org/zap"
)

// GitHubClient is a client for the GitHub API
type GitHubClient struct {
	client     *http.Client
	baseURL    string
	token      string
	apiVersion string
	logger     *zap.Logger
}

// NewGitHubClient creates a new GitHub API client
func NewGitHubClient(baseURL, token, apiVersion string, logger *zap.Logger) *GitHubClient {
	if baseURL == "" {
		baseURL = "https://api.github.com"
	}

	// Ensure baseURL doesn't end with a slash
	baseURL = strings.TrimSuffix(baseURL, "/")

	return &GitHubClient{
		client:     &http.Client{Timeout: 30 * time.Second},
		baseURL:    baseURL,
		token:      token,
		apiVersion: apiVersion,
		logger:     logger,
	}
}

// GetNextPageURL extracts the next page URL from Link header if present
func GetNextPageURL(linkHeader string) string {
	if linkHeader == "" {
		return ""
	}

	links := strings.Split(linkHeader, ",")
	for _, link := range links {
		parts := strings.Split(link, ";")
		if len(parts) == 2 && strings.Contains(parts[1], `rel="next"`) {
			urlPart := strings.TrimSpace(parts[0])
			// Remove angle brackets
			return strings.Trim(urlPart, "<>")
		}
	}

	return ""
}

// createRequest creates a new HTTP request with the appropriate headers
func (g *GitHubClient) createRequest(method, path string, body io.Reader) (*http.Request, error) {
	// Ensure path starts with a slash
	if !strings.HasPrefix(path, "/") && !strings.HasPrefix(path, "http") {
		path = "/" + path
	}

	// If path is already an absolute URL, use it directly
	var fullURL string
	if strings.HasPrefix(path, "http") {
		fullURL = path
	} else {
		fullURL = fmt.Sprintf("%s%s", g.baseURL, path)
	}

	// Validate the URL
	_, err := url.Parse(fullURL)
	if err != nil {
		g.logger.Error("Invalid URL", zap.String("url", fullURL), zap.Error(err))
		return nil, fmt.Errorf("invalid URL %s: %w", fullURL, err)
	}

	req, err := http.NewRequest(method, fullURL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+g.token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", g.apiVersion)
	req.Header.Set("User-Agent", "GitHubCopilotMetricsDashboard")

	return req, nil
}
