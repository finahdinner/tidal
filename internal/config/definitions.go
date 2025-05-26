package config

import "github.com/finahdinner/tidal/internal/helpers"

type PreferencesFormat struct {
	TwitchConfig                        TwitchConfigT    `json:"twitch_config"`
	TwitchVariables                     TwitchVariablesT `json:"twitch_variables"`
	TwitchVariableUpdateIntervalSeconds int              `json:"twitch_variable_update_interval_seconds"`
	LlmConfig                           LlmConfigT       `json:"llm_config"`
	AiGeneratedVariables                []LlmVariableT   `json:"ai_generated_variables"`
	TitleConfig                         TitleConfigT     `json:"title_config"`
	ActivityConsoleOutput               string           `json:"activity_console_output"`
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

type LlmConfigT struct {
	Provider string `json:"provider"`
	ApiKey   string `json:"api_key"`
}

type LlmVariableT struct {
	Name         string `json:"name"`
	Value        string `json:"value"`
	PromptMain   string `json:"prompt_main"`
	PromptSuffix string `json:"prompt_suffix"`
}

type TitleConfigT struct {
	TitleTemplate              string `json:"title_template"`
	TitleUpdateIntervalMinutes int    `json:"title_update_interval_minutes"`
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
		len(pf.TwitchConfig.Credentials.UserAccessScope) > 0 &&
		pf.TitleConfig.TitleTemplate != "" &&
		pf.TitleConfig.TitleUpdateIntervalMinutes >= helpers.MinTitleUpdateIntervalMinutes &&
		pf.TitleConfig.TitleUpdateIntervalMinutes <= helpers.MaxTitleUpdateIntervalMinutes
}
