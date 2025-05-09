package preferences

type preferencesFormat struct {
	TwitchConfig          twitchConfigT    `json:"twitch_config"`
	StreamVariables       streamVariablesT `json:"stream_variables"`
	LlmVariables          []llmVariableT   `json:"llm_variables"`
	ActivityConsoleOutput string           `json:"activity_console_output"`
	UpdateFrequency       uint16           `json:"update_frequency"`
}

type twitchConfigT struct {
	UserName     string       `json:"user_name"`
	UserId       string       `json:"user_id"`
	ClientId     string       `json:"client_id"`
	ClientSecret string       `json:"client_secret"`
	RedirectUri  string       `json:"redirect_uri"`
	Credentials  credentialsT `json:"credentials"`
}

type credentialsT struct {
	UserAccessToken        string   `json:"user_access_token"`
	UserAccessRefreshToken string   `json:"user_refresh_token"`
	UserAccessScope        []string `json:"user_access_scope"`
}

type streamVariablesT struct {
	StreamCategory       streamVariableT `json:"stream_category"`
	StreamUptime         streamVariableT `json:"stream_uptime"`
	NumViewers           streamVariableT `json:"num_viewers"`
	NumSubscribers       streamVariableT `json:"num_subscribers"`
	NumFollowers         streamVariableT `json:"num_followers"`
	MostRecentSubscriber streamVariableT `json:"most_recent_subscriber"`
	MostRecentFollower   streamVariableT `json:"most_recent_follower"`
}

type streamVariableT struct {
	Value       string `json:"value"`
	Description string `json:"description"`
}

type llmVariableT struct {
	Value  string `json:"value"`
	Prompt string `json:"prompt"`
}
