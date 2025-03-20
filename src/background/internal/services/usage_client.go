package services

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/cardonator/copilot-metrics-dashboard/internal/models"
	"go.uber.org/zap"
)

// CopilotUsageClient handles fetching Copilot usage data from GitHub API
type CopilotUsageClient struct {
	githubClient *GitHubClient
	logger       *zap.Logger
}

// NewCopilotUsageClient creates a new Copilot usage client
func NewCopilotUsageClient(githubClient *GitHubClient, logger *zap.Logger) *CopilotUsageClient {
	return &CopilotUsageClient{
		githubClient: githubClient,
		logger:       logger,
	}
}

// GetCopilotUsageForOrganization fetches Copilot usage statistics for an organization
func (c *CopilotUsageClient) GetCopilotUsageForOrganization(organization string) ([]models.CopilotUsage, error) {
	requestURI := fmt.Sprintf("/orgs/%s/copilot/usage", organization)

	req, err := c.githubClient.createRequest("GET", requestURI, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.githubClient.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch usage data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var usage []models.CopilotUsage
	if err := json.Unmarshal(body, &usage); err != nil {
		return nil, fmt.Errorf("failed to unmarshal usage data: %w", err)
	}

	// Set IDs for each usage record if not already set
	for i := range usage {
		if usage[i].ID == "" {
			usage[i].ID = usage[i].GetID()
		}
	}

	return usage, nil
}

// GetCopilotUsageForEnterprise fetches Copilot usage statistics for an enterprise
func (c *CopilotUsageClient) GetCopilotUsageForEnterprise(enterprise string) ([]models.CopilotUsage, error) {
	requestURI := fmt.Sprintf("/enterprises/%s/copilot/usage", enterprise)

	req, err := c.githubClient.createRequest("GET", requestURI, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.githubClient.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch usage data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var usage []models.CopilotUsage
	if err := json.Unmarshal(body, &usage); err != nil {
		return nil, fmt.Errorf("failed to unmarshal usage data: %w", err)
	}

	// Set IDs for each usage record if not already set
	for i := range usage {
		if usage[i].ID == "" {
			usage[i].ID = usage[i].GetID()
		}
	}

	return usage, nil
}

// LoadTestUsageData loads test usage data from a file
func (c *CopilotUsageClient) LoadTestUsageData() ([]models.CopilotUsage, error) {
	data, err := loadTestData("usage.json")
	if err != nil {
		return nil, err
	}

	var usage []models.CopilotUsage
	if err := json.Unmarshal(data, &usage); err != nil {
		return nil, fmt.Errorf("failed to unmarshal test usage data: %w", err)
	}

	// Set IDs for each usage record if not already set
	for i := range usage {
		if usage[i].ID == "" {
			usage[i].ID = usage[i].GetID()
		}
	}

	return usage, nil
}
