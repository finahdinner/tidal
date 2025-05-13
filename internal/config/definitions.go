package config

type PreferencesFormat struct {
	TwitchConfig           TwitchConfigT    `json:"twitch_config"`
	TwitchVariables        TwitchVariablesT `json:"twitch_variables"`
	AiGeneratedVariables   []LlmVariableT   `json:"ai_generated_variables"`
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

type TwitchVariablesT struct {
	StreamCategory TwitchVariableT `json:"stream_category"`
	StreamUptime   TwitchVariableT `json:"stream_uptime"`
	NumViewers     TwitchVariableT `json:"num_viewers"`
	NumSubscribers TwitchVariableT `json:"num_subscribers"`
	NumFollowers   TwitchVariableT `json:"num_followers"`
}

type TwitchVariableT struct {
	Value       string `json:"value"`
	Description string `json:"description"`
}

type LlmVariableT struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Prompt string `json:"prompt"`
}

// Ensure fields are populated enough to make requests to update twitch variables
func (pf *PreferencesFormat) IsValidForUpdatingTwitchVariables() bool {
	return pf.TwitchConfig.UserName != "" &&
		pf.TwitchConfig.UserId != "" &&
		pf.TwitchConfig.ClientId != "" &&
		pf.TwitchConfig.ClientSecret != "" &&
		pf.TwitchConfig.ClientRedirectUri != "" &&
		pf.TwitchConfig.Credentials.UserAccessToken != "" &&
		pf.TwitchConfig.Credentials.UserAccessRefreshToken != "" &&
		len(pf.TwitchConfig.Credentials.UserAccessScope) > 0
}
