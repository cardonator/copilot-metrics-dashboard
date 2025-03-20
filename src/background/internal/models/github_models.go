package models

// Organization represents a GitHub organization
type Organization struct {
	Login            string `json:"login"`
	ID               int    `json:"id"`
	NodeID           string `json:"node_id"`
	URL              string `json:"url"`
	ReposURL         string `json:"repos_url"`
	EventsURL        string `json:"events_url"`
	HooksURL         string `json:"hooks_url"`
	IssuesURL        string `json:"issues_url"`
	MembersURL       string `json:"members_url"`
	PublicMembersURL string `json:"public_members_url"`
	AvatarURL        string `json:"avatar_url"`
	Description      string `json:"description"`
}

// Team represents a GitHub team
type Team struct {
	ID                  int         `json:"id"`
	NodeID              string      `json:"node_id"`
	URL                 string      `json:"url"`
	HTMLURL             string      `json:"html_url"`
	Name                string      `json:"name"`
	Slug                string      `json:"slug"`
	Description         string      `json:"description"`
	Privacy             string      `json:"privacy"`
	NotificationSetting string      `json:"notification_setting"`
	Permission          string      `json:"permission"`
	MembersURL          string      `json:"members_url"`
	RepositoriesURL     string      `json:"repositories_url"`
	Parent              interface{} `json:"parent"`
}

// User represents a GitHub user
type User struct {
	ID                int    `json:"id"`
	Login             string `json:"login"`
	Name              string `json:"name"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}
