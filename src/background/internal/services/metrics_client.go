package services

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/cardonator/copilot-metrics-dashboard/internal/models"

	"go.uber.org/zap"
)

// CopilotMetricsClient handles fetching Copilot metrics from GitHub API
type CopilotMetricsClient struct {
	githubClient *GitHubClient
	logger       *zap.Logger
}

// NewCopilotMetricsClient creates a new Copilot metrics client
func NewCopilotMetricsClient(githubClient *GitHubClient, logger *zap.Logger) *CopilotMetricsClient {
	return &CopilotMetricsClient{
		githubClient: githubClient,
		logger:       logger,
	}
}

// GetCopilotMetricsForEnterprise fetches Copilot metrics for an enterprise
func (c *CopilotMetricsClient) GetCopilotMetricsForEnterprise(enterprise, team string) ([]models.Metrics, error) {
	var requestURI string
	if team == "" {
		requestURI = fmt.Sprintf("/enterprises/%s/copilot/metrics", enterprise)
	} else {
		requestURI = fmt.Sprintf("/enterprises/%s/team/%s/copilot/metrics", enterprise, team)
	}

	req, err := c.githubClient.createRequest("GET", requestURI, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.githubClient.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch metrics: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		c.logger.Warn("Team not found", zap.String("team", team))
		return []models.Metrics{}, nil
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var metrics []models.Metrics
	if err := json.Unmarshal(body, &metrics); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metrics: %w", err)
	}

	// Add metadata
	for i := range metrics {
		metrics[i].Enterprise = enterprise
		metrics[i].Team = team
		metrics[i].LastUpdate = time.Now().UTC()
	}

	return metrics, nil
}

// GetCopilotMetricsForOrganization fetches Copilot metrics for an organization
func (c *CopilotMetricsClient) GetCopilotMetricsForOrganization(organization, team string) ([]models.Metrics, error) {
	var requestURI string
	if team == "" {
		requestURI = fmt.Sprintf("/orgs/%s/copilot/metrics", organization)
	} else {
		requestURI = fmt.Sprintf("/orgs/%s/team/%s/copilot/metrics", organization, team)
	}

	req, err := c.githubClient.createRequest("GET", requestURI, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.githubClient.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch metrics: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		c.logger.Warn("Team not found", zap.String("team", team))
		return []models.Metrics{}, nil
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var metrics []models.Metrics
	if err := json.Unmarshal(body, &metrics); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metrics: %w", err)
	}

	// Add metadata
	for i := range metrics {
		metrics[i].Organization = organization
		metrics[i].Team = team
		metrics[i].LastUpdate = time.Now().UTC()
	}

	return metrics, nil
}

// LoadTestMetrics loads test metrics from a file
func (c *CopilotMetricsClient) LoadTestMetrics(team string) ([]models.Metrics, error) {
	data, err := loadTestData("metrics.json")
	if err != nil {
		return nil, err
	}

	var metrics []models.Metrics
	if err := json.Unmarshal(data, &metrics); err != nil {
		return nil, fmt.Errorf("failed to unmarshal test metrics data: %w", err)
	}

	// Add metadata
	for i := range metrics {
		metrics[i].Organization = "test"
		metrics[i].Team = team
		metrics[i].LastUpdate = time.Now().UTC()
	}

	return metrics, nil
}
