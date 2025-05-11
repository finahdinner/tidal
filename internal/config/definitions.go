package config

type PreferencesFormat struct {
	TwitchConfig           TwitchConfigT    `json:"twitch_config"`
	StreamVariables        StreamVariablesT `json:"stream_variables"`
	LlmVariables           []LlmVariableT   `json:"llm_variables"`
	VariableUpdateInterval int              `json:"variable_update_interval"`
	ActivityConsoleOutput  string           `json:"activity_console_output"`
}

type TwitchConfigT struct {
	UserName          string       `json:"user_name"`
	UserId            string       `json:"user_id"`
	ClientId          string       `json:"client_id"`
	ClientSecret      string       `json:"client_secret"`
	ClientRedirectUri string       `json:"client_redirect_uri"`
	Credentials       CredentialsT `json:"credentials"`
}

type CredentialsT struct {
	UserAccessToken        string   `json:"user_access_token"`
	UserAccessRefreshToken string   `json:"user_refresh_token"`
	UserAccessScope        []string `json:"user_access_scope"`
	ExpiryUnixTimestamp    int64    `json:"expiry_unix_timestamp"`
}

type StreamVariablesT struct {
	StreamCategory StreamVariableT `json:"stream_category"`
	StreamUptime   StreamVariableT `json:"stream_uptime"`
	NumViewers     StreamVariableT `json:"num_viewers"`
	NumSubscribers StreamVariableT `json:"num_subscribers"`
	NumFollowers   StreamVariableT `json:"num_followers"`
}

type StreamVariableT struct {
	Value       string `json:"value"`
	Description string `json:"description"`
}

type LlmVariableT struct {
	Value  string `json:"value"`
	Prompt string `json:"prompt"`
}

// Ensure fields are populated enough to make requests to update stream variables
func (pf *PreferencesFormat) IsValidForUpdatingStreamVariables() bool {
	return pf.TwitchConfig.UserName != "" &&
		pf.TwitchConfig.UserId != "" &&
		pf.TwitchConfig.ClientId != "" &&
		pf.TwitchConfig.ClientSecret != "" &&
		pf.TwitchConfig.ClientRedirectUri != "" &&
		pf.TwitchConfig.Credentials.UserAccessToken != "" &&
		pf.TwitchConfig.Credentials.UserAccessRefreshToken != "" &&
		len(pf.TwitchConfig.Credentials.UserAccessScope) > 0
}
