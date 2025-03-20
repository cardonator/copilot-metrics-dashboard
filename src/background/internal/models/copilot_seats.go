package models

import (
	"fmt"
	"time"
)

// CopilotAssignedSeats represents seats assigned to GitHub Copilot within an organization or enterprise
type CopilotAssignedSeats struct {
	ID           string    `json:"id,omitempty"`
	Date         string    `json:"date"`
	TotalSeats   int       `json:"total_seats"`
	Seats        []Seat    `json:"seats"`
	Enterprise   string    `json:"enterprise,omitempty"`
	Organization string    `json:"organization,omitempty"`
	LastUpdate   time.Time `json:"last_update"`
}

// GetID generates an ID for the seats data
func (c *CopilotAssignedSeats) GetID() string {
	if c.Organization != "" {
		return fmt.Sprintf("%s-ORG-%s", c.Date, c.Organization)
	} else if c.Enterprise != "" {
		return fmt.Sprintf("%s-ENT-%s", c.Date, c.Enterprise)
	}
	return fmt.Sprintf("%s-XXX", c.Date)
}

// Seat represents a seat assigned to a user within GitHub Copilot
type Seat struct {
	CreatedAt               time.Time     `json:"created_at"`
	UpdatedAt               time.Time     `json:"updated_at"`
	PendingCancellationDate string        `json:"pending_cancellation_date"`
	LastActivityAt          *time.Time    `json:"last_activity_at,omitempty"`
	LastActivityEditor      string        `json:"last_activity_editor,omitempty"`
	PlanType                string        `json:"plan_type"`
	Assignee                User          `json:"assignee"`
	AssigningTeam           *Team         `json:"assigning_team,omitempty"`
	Organization            *Organization `json:"organization,omitempty"`
}
