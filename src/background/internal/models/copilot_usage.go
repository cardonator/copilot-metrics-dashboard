package models

// CopilotUsage represents GitHub Copilot usage statistics
type CopilotUsage struct {
	TotalSuggestionsCount int         `json:"total_suggestions_count"`
	TotalAcceptancesCount int         `json:"total_acceptances_count"`
	TotalLinesSuggested   int         `json:"total_lines_suggested"`
	TotalLinesAccepted    int         `json:"total_lines_accepted"`
	TotalActiveUsers      int         `json:"total_active_users"`
	TotalChatAcceptances  int         `json:"total_chat_acceptances"`
	TotalChatTurns        int         `json:"total_chat_turns"`
	TotalActiveChatUsers  int         `json:"total_active_chat_users"`
	Day                   string      `json:"day"`
	ID                    string      `json:"id,omitempty"`
	Breakdown             []Breakdown `json:"breakdown"`
}

// GetID generates an ID for the usage data
func (c *CopilotUsage) GetID() string {
	return c.Day
}

// Breakdown represents usage statistics broken down by language and editor
type Breakdown struct {
	Language         string `json:"language"`
	Editor           string `json:"editor"`
	SuggestionsCount int    `json:"suggestions_count"`
	AcceptancesCount int    `json:"acceptances_count"`
	LinesSuggested   int    `json:"lines_suggested"`
	LinesAccepted    int    `json:"lines_accepted"`
	ActiveUsers      int    `json:"active_users"`
}
