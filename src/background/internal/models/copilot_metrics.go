package models

import (
	"fmt"
	"time"
)

// Metrics represents GitHub Copilot usage metrics
type Metrics struct {
	ID                        string              `json:"id,omitempty"`
	Date                      string              `json:"date"`
	TotalActiveUsers          int                 `json:"total_active_users"`
	TotalEngagedUsers         int                 `json:"total_engaged_users"`
	CopilotIdeCodeCompletions *IdeCodeCompletions `json:"copilot_ide_code_completions,omitempty"`
	IdeChat                   *IdeChat            `json:"copilot_ide_chat,omitempty"`
	DotComChat                *DotComChat         `json:"copilot_dotcom_chat,omitempty"`
	DotComPullRequests        *DotComPullRequest  `json:"copilot_dotcom_pull_requests,omitempty"`
	Enterprise                string              `json:"enterprise,omitempty"`
	Organization              string              `json:"organization,omitempty"`
	Team                      string              `json:"team,omitempty"`
	LastUpdate                time.Time           `json:"last_update"`
}

// GetID generates an ID for the metrics data
func (m *Metrics) GetID() string {
	if m.Organization != "" {
		teamSuffix := ""
		if m.Team != "" {
			teamSuffix = "-" + m.Team
		}
		return fmt.Sprintf("%s-ORG-%s%s", m.Date, m.Organization, teamSuffix)
	} else if m.Enterprise != "" {
		teamSuffix := ""
		if m.Team != "" {
			teamSuffix = "-" + m.Team
		}
		return fmt.Sprintf("%s-ENT-%s%s", m.Date, m.Enterprise, teamSuffix)
	}
	return fmt.Sprintf("%s-XXX", m.Date)
}

// IdeCodeCompletions represents IDE code completion metrics
type IdeCodeCompletions struct {
	TotalEngagedUsers int                         `json:"total_engaged_users"`
	Languages         []IdeCodeCompletionLanguage `json:"languages"`
	Editors           []IdeCodeCompletionEditor   `json:"editors"`
}

// IdeCodeCompletionLanguage represents language-specific IDE code completion metrics
type IdeCodeCompletionLanguage struct {
	Name              string `json:"name"`
	TotalEngagedUsers int    `json:"total_engaged_users"`
}

// IdeCodeCompletionEditor represents editor-specific IDE code completion metrics
type IdeCodeCompletionEditor struct {
	Name              string                   `json:"name"`
	TotalEngagedUsers int                      `json:"total_engaged_users"`
	Models            []IdeCodeCompletionModel `json:"models"`
}

// IdeCodeCompletionModel represents model-specific IDE code completion metrics
type IdeCodeCompletionModel struct {
	Name                    string                           `json:"name"`
	IsCustomModel           bool                             `json:"is_custom_model"`
	CustomModelTrainingDate *string                          `json:"custom_model_training_date,omitempty"`
	TotalEngagedUsers       int                              `json:"total_engaged_users"`
	Languages               []IdeCodeCompletionModelLanguage `json:"languages"`
}

// IdeCodeCompletionModelLanguage represents language-specific model metrics for IDE code completions
type IdeCodeCompletionModelLanguage struct {
	Name                    string `json:"name"`
	TotalEngagedUsers       int    `json:"total_engaged_users"`
	TotalCodeSuggestions    int    `json:"total_code_suggestions"`
	TotalCodeAcceptances    int    `json:"total_code_acceptances"`
	TotalCodeLinesSuggested int    `json:"total_code_lines_suggested"`
	TotalCodeLinesAccepted  int    `json:"total_code_lines_accepted"`
}

// IdeChat represents IDE chat metrics
type IdeChat struct {
	TotalEngagedUsers int             `json:"total_engaged_users"`
	Editors           []IdeChatEditor `json:"editors"`
}

// IdeChatEditor represents editor-specific IDE chat metrics
type IdeChatEditor struct {
	Name              string         `json:"name"`
	TotalEngagedUsers int            `json:"total_engaged_users"`
	Models            []IdeChatModel `json:"models"`
}

// IdeChatModel represents model-specific IDE chat metrics
type IdeChatModel struct {
	Name                     string  `json:"name"`
	IsCustomModel            bool    `json:"is_custom_model"`
	CustomModelTrainingDate  *string `json:"custom_model_training_date,omitempty"`
	TotalEngagedUsers        int     `json:"total_engaged_users"`
	TotalChats               int     `json:"total_chats"`
	TotalChatInsertionEvents int     `json:"total_chat_insertion_events"`
	TotalChatCopyEvents      int     `json:"total_chat_copy_events"`
}

// DotComChat represents GitHub.com chat metrics
type DotComChat struct {
	TotalEngagedUsers int               `json:"total_engaged_users"`
	Models            []DotComChatModel `json:"models"`
}

// DotComChatModel represents model-specific GitHub.com chat metrics
type DotComChatModel struct {
	Name                    string  `json:"name"`
	IsCustomModel           bool    `json:"is_custom_model"`
	CustomModelTrainingDate *string `json:"custom_model_training_date,omitempty"`
	TotalEngagedUsers       int     `json:"total_engaged_users"`
	TotalChats              int     `json:"total_chats"`
}

// DotComPullRequest represents GitHub.com pull request metrics
type DotComPullRequest struct {
	TotalEngagedUsers int                           `json:"total_engaged_users"`
	Repositories      []DotComPullRequestRepository `json:"repositories"`
}

// DotComPullRequestRepository represents repository-specific GitHub.com pull request metrics
type DotComPullRequestRepository struct {
	Name              string                             `json:"name"`
	TotalEngagedUsers int                                `json:"total_engaged_users"`
	Models            []DotComPullRequestRepositoryModel `json:"models"`
}

// DotComPullRequestRepositoryModel represents model-specific GitHub.com pull request metrics
type DotComPullRequestRepositoryModel struct {
	Name                    string  `json:"name"`
	IsCustomModel           bool    `json:"is_custom_model"`
	CustomModelTrainingDate *string `json:"custom_model_training_date,omitempty"`
	TotalEngagedUsers       int     `json:"total_engaged_users"`
	TotalPrSummariesCreated int     `json:"total_pr_summaries_created"`
}

// CopilotUsage represents GitHub Copilot usage statistics
type CopilotUsage struct {
	ID                    string           `json:"id,omitempty"`
	Day                   string           `json:"day"`
	Enterprise            string           `json:"enterprise,omitempty"`
	Organization          string           `json:"organization,omitempty"`
	Team                  string           `json:"team,omitempty"`
	LastUpdate            time.Time        `json:"last_update"`
	TotalSuggestionsCount int              `json:"total_suggestions_count"`
	TotalAcceptancesCount int              `json:"total_acceptances_count"`
	TotalLinesSuggested   int              `json:"total_lines_suggested"`
	TotalLinesAccepted    int              `json:"total_lines_accepted"`
	TotalActiveUsers      int              `json:"total_active_users"`
	TotalChatAcceptances  int              `json:"total_chat_acceptances"`
	TotalChatTurns        int              `json:"total_chat_turns"`
	TotalActiveChatUsers  int              `json:"total_active_chat_users"`
	Breakdown             []UsageBreakdown `json:"breakdown"`
}

// GetID generates an ID for the usage data
func (c *CopilotUsage) GetID() string {
	if c.Organization != "" {
		return fmt.Sprintf("%s-ORG-%s", c.Day, c.Organization)
	} else if c.Enterprise != "" {
		return fmt.Sprintf("%s-ENT-%s", c.Day, c.Enterprise)
	}
	return fmt.Sprintf("%s-XXX", c.Day)
}

// UsageBreakdown represents usage statistics broken down by language and editor
type UsageBreakdown struct {
	Day              string `json:"day"`
	Language         string `json:"language"`
	Editor           string `json:"editor"`
	SuggestionsCount int    `json:"suggestions_count"`
	AcceptancesCount int    `json:"acceptances_count"`
	LinesSuggested   int    `json:"lines_suggested"`
	LinesAccepted    int    `json:"lines_accepted"`
	ActiveUsers      int    `json:"active_users"`
	Enterprise       string `json:"enterprise,omitempty"`
	Organization     string `json:"organization,omitempty"`
	Team             string `json:"team,omitempty"`
}
