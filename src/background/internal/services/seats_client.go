package services

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/cardonator/copilot-metrics-dashboard/internal/models"
	"go.uber.org/zap"
)

// CopilotSeatsClient handles fetching Copilot seats from GitHub API
type CopilotSeatsClient struct {
	githubClient *GitHubClient
	logger       *zap.Logger
}

// NewCopilotSeatsClient creates a new Copilot seats client
func NewCopilotSeatsClient(githubClient *GitHubClient, logger *zap.Logger) *CopilotSeatsClient {
	return &CopilotSeatsClient{
		githubClient: githubClient,
		logger:       logger,
	}
}

// GetEnterpriseAssignedSeats fetches Copilot seats for an enterprise
func (c *CopilotSeatsClient) GetEnterpriseAssignedSeats(enterprise string) (*models.CopilotAssignedSeats, error) {
	path := fmt.Sprintf("/enterprises/%s/copilot/billing/seats", enterprise)
	allSeats := []models.Seat{}

	for path != "" {
		req, err := c.githubClient.createRequest("GET", path, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		resp, err := c.githubClient.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch seats: %w", err)
		}

		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, body)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		var data models.CopilotAssignedSeats
		if err := json.Unmarshal(body, &data); err != nil {
			return nil, fmt.Errorf("failed to unmarshal seats: %w", err)
		}

		allSeats = append(allSeats, data.Seats...)

		// Get next page path, handling both relative and absolute URLs
		nextLink := GetNextPageURL(resp.Header.Get("Link"))
		if nextLink == "" {
			path = ""
		} else if strings.HasPrefix(nextLink, "http") {
			// If it's an absolute URL, extract just the path
			// First, find where the base URL ends
			baseURL := c.githubClient.baseURL
			if !strings.HasPrefix(nextLink, baseURL) {
				c.logger.Warn("Next page URL doesn't match base URL",
					zap.String("nextLink", nextLink),
					zap.String("baseURL", baseURL))
				// Try to extract the path anyway by finding the path after hostname
				parts := strings.SplitN(nextLink, "/", 4)
				if len(parts) >= 4 {
					path = "/" + parts[3]
				} else {
					path = ""
				}
			} else {
				// Extract path from absolute URL
				path = nextLink[len(baseURL):]
			}
		} else {
			// It's already a relative path
			path = nextLink
		}

		c.logger.Debug("Pagination", zap.String("nextPath", path))
	}

	now := time.Now().UTC()
	currentDate := now.Format("2006-01-02")

	result := &models.CopilotAssignedSeats{
		TotalSeats: len(allSeats),
		Enterprise: enterprise,
		LastUpdate: now,
		Date:       currentDate,
		Seats:      allSeats,
	}

	return result, nil
}

// GetOrganizationAssignedSeats fetches Copilot seats for an organization
func (c *CopilotSeatsClient) GetOrganizationAssignedSeats(organization string) (*models.CopilotAssignedSeats, error) {
	path := fmt.Sprintf("/orgs/%s/copilot/billing/seats", organization)
	allSeats := []models.Seat{}

	for path != "" {
		req, err := c.githubClient.createRequest("GET", path, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		resp, err := c.githubClient.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch seats: %w", err)
		}

		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, body)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		var data models.CopilotAssignedSeats
		if err := json.Unmarshal(body, &data); err != nil {
			return nil, fmt.Errorf("failed to unmarshal seats: %w", err)
		}

		allSeats = append(allSeats, data.Seats...)

		// Get next page path, handling both relative and absolute URLs
		nextLink := GetNextPageURL(resp.Header.Get("Link"))
		if nextLink == "" {
			path = ""
		} else if strings.HasPrefix(nextLink, "http") {
			// If it's an absolute URL, extract just the path
			// First, find where the base URL ends
			baseURL := c.githubClient.baseURL
			if !strings.HasPrefix(nextLink, baseURL) {
				c.logger.Warn("Next page URL doesn't match base URL",
					zap.String("nextLink", nextLink),
					zap.String("baseURL", baseURL))
				// Try to extract the path anyway by finding the path after hostname
				parts := strings.SplitN(nextLink, "/", 4)
				if len(parts) >= 4 {
					path = "/" + parts[3]
				} else {
					path = ""
				}
			} else {
				// Extract path from absolute URL
				path = nextLink[len(baseURL):]
			}
		} else {
			// It's already a relative path
			path = nextLink
		}

		c.logger.Debug("Pagination", zap.String("nextPath", path))
	}

	now := time.Now().UTC()
	currentDate := now.Format("2006-01-02")

	result := &models.CopilotAssignedSeats{
		TotalSeats:   len(allSeats),
		Organization: organization,
		LastUpdate:   now,
		Date:         currentDate,
		Seats:        allSeats,
	}

	return result, nil
}

// LoadTestSeatsData loads test seats data from a file
func (c *CopilotSeatsClient) LoadTestSeatsData(isEnterprise bool) (*models.CopilotAssignedSeats, error) {
	data, err := loadTestData("seats.json")
	if err != nil {
		return nil, err
	}

	var seats models.CopilotAssignedSeats
	if err := json.Unmarshal(data, &seats); err != nil {
		return nil, fmt.Errorf("failed to unmarshal test seats data: %w", err)
	}

	// Add appropriate metadata
	now := time.Now().UTC()
	currentDate := now.Format("2006-01-02")
	seats.LastUpdate = now
	seats.Date = currentDate

	if isEnterprise {
		seats.Enterprise = "test-enterprise"
		seats.Organization = ""
	} else {
		seats.Organization = "test-organization"
		seats.Enterprise = ""
	}

	return &seats, nil
}
