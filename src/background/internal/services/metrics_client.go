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

// GetCopilotUsageFromMetrics converts metrics to usage format
func (c *CopilotMetricsClient) GetCopilotUsageFromMetrics(metrics []models.Metrics) ([]models.CopilotUsage, error) {
	usagesByDay := make(map[string]*models.CopilotUsage)

	// Convert metrics to usage format
	for _, metric := range metrics {
		var breakdown []models.UsageBreakdown
		totalSuggestionsCount := 0
		totalAcceptancesCount := 0
		totalLinesSuggested := 0
		totalLinesAccepted := 0
		totalActiveUsers := 0
		totalChatAcceptances := 0
		totalChatTurns := 0
		totalActiveChatUsers := 0

		// Set total engaged users from the metrics
		if metric.TotalEngagedUsers > 0 {
			totalActiveUsers = metric.TotalEngagedUsers
		}

		// Sum up code completion metrics from all editors and languages
		if metric.CopilotIdeCodeCompletions != nil {
			// Track unique users across all editors and languages
			if metric.CopilotIdeCodeCompletions.TotalEngagedUsers > 0 {
				totalActiveUsers = metric.CopilotIdeCodeCompletions.TotalEngagedUsers
			}

			for _, editor := range metric.CopilotIdeCodeCompletions.Editors {
				for _, model := range editor.Models {
					for _, language := range model.Languages {
						// Add breakdown by language/editor
						usageBreakdown := models.UsageBreakdown{
							Day:              metric.Date,
							Language:         language.Name,
							Editor:           editor.Name,
							SuggestionsCount: language.TotalCodeSuggestions,
							AcceptancesCount: language.TotalCodeAcceptances,
							LinesSuggested:   language.TotalCodeLinesSuggested,
							LinesAccepted:    language.TotalCodeLinesAccepted,
							ActiveUsers:      language.TotalEngagedUsers,
							Enterprise:       metric.Enterprise,
							Organization:     metric.Organization,
							Team:             metric.Team,
						}
						breakdown = append(breakdown, usageBreakdown)

						// Aggregate totals
						totalSuggestionsCount += language.TotalCodeSuggestions
						totalAcceptancesCount += language.TotalCodeAcceptances
						totalLinesSuggested += language.TotalCodeLinesSuggested
						totalLinesAccepted += language.TotalCodeLinesAccepted
					}
				}
			}
		}

		// Get chat data
		if metric.IdeChat != nil {
			totalActiveChatUsers = metric.IdeChat.TotalEngagedUsers

			for _, editor := range metric.IdeChat.Editors {
				for _, model := range editor.Models {
					if model.TotalChats > 0 {
						totalChatTurns += model.TotalChats
					}
					if model.TotalChatCopyEvents > 0 {
						totalChatAcceptances += model.TotalChatCopyEvents
					}
					if model.TotalChatInsertionEvents > 0 {
						totalChatAcceptances += model.TotalChatInsertionEvents
					}
				}
			}
		}

		// Add chat data from GitHub.com
		if metric.DotComChat != nil && metric.DotComChat.TotalEngagedUsers > 0 {
			totalActiveChatUsers += metric.DotComChat.TotalEngagedUsers

			for _, model := range metric.DotComChat.Models {
				if model.TotalChats > 0 {
					totalChatTurns += model.TotalChats
				}
			}
		}

		// Get or create usage record for this date
		usage, exists := usagesByDay[metric.Date]
		if !exists {
			usage = &models.CopilotUsage{
				ID:                    "", // Will be set below
				Day:                   metric.Date,
				Organization:          metric.Organization,
				Enterprise:            metric.Enterprise,
				Team:                  metric.Team,
				LastUpdate:            time.Now().UTC(),
				TotalSuggestionsCount: totalSuggestionsCount,
				TotalAcceptancesCount: totalAcceptancesCount,
				TotalLinesSuggested:   totalLinesSuggested,
				TotalLinesAccepted:    totalLinesAccepted,
				TotalActiveUsers:      totalActiveUsers,
				TotalChatAcceptances:  totalChatAcceptances,
				TotalChatTurns:        totalChatTurns,
				TotalActiveChatUsers:  totalActiveChatUsers,
				Breakdown:             []models.UsageBreakdown{},
			}
			usagesByDay[metric.Date] = usage
		} else {
			// Increment existing totals
			usage.TotalSuggestionsCount += totalSuggestionsCount
			usage.TotalAcceptancesCount += totalAcceptancesCount
			usage.TotalLinesSuggested += totalLinesSuggested
			usage.TotalLinesAccepted += totalLinesAccepted
			usage.TotalChatAcceptances += totalChatAcceptances
			usage.TotalChatTurns += totalChatTurns

			// Take the max value for users to avoid double counting
			if totalActiveUsers > usage.TotalActiveUsers {
				usage.TotalActiveUsers = totalActiveUsers
			}
			if totalActiveChatUsers > usage.TotalActiveChatUsers {
				usage.TotalActiveChatUsers = totalActiveChatUsers
			}
		}

		// Add breakdowns to usage
		usage.Breakdown = append(usage.Breakdown, breakdown...)
	}

	// Convert map to array and set IDs
	var result []models.CopilotUsage
	for _, usage := range usagesByDay {
		usage.ID = usage.GetID()
		result = append(result, *usage)
	}

	return result, nil
}

// LoadTestUsageData loads test usage data by transforming test metrics data
func (c *CopilotMetricsClient) LoadTestUsageData() ([]models.CopilotUsage, error) {
	// Get test metrics data
	metrics, err := c.LoadTestMetrics("")
	if err != nil {
		return nil, err
	}

	// Convert to usage format
	return c.GetCopilotUsageFromMetrics(metrics)
}

// GetCopilotUsageForEnterprise generates Copilot usage data for an enterprise from metrics
func (c *CopilotMetricsClient) GetCopilotUsageForEnterprise(enterprise string) ([]models.CopilotUsage, error) {
	// Get metrics data first
	metrics, err := c.GetCopilotMetricsForEnterprise(enterprise, "")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch metrics: %w", err)
	}

	// Convert metrics to usage format
	usages, err := c.GetCopilotUsageFromMetrics(metrics)
	if err != nil {
		return nil, fmt.Errorf("failed to convert metrics to usage: %w", err)
	}

	return usages, nil
}

// GetCopilotUsageForOrganization generates Copilot usage data for an organization from metrics
func (c *CopilotMetricsClient) GetCopilotUsageForOrganization(organization string) ([]models.CopilotUsage, error) {
	// Get metrics data first
	metrics, err := c.GetCopilotMetricsForOrganization(organization, "")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch metrics: %w", err)
	}

	// Convert metrics to usage format
	usages, err := c.GetCopilotUsageFromMetrics(metrics)
	if err != nil {
		return nil, fmt.Errorf("failed to convert metrics to usage: %w", err)
	}

	return usages, nil
}
