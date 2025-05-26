package config

var defaultPreferences PreferencesFormat = PreferencesFormat{
	TwitchConfig: TwitchConfigT{
		UserName:          "",
		UserId:            "",
		ClientId:          "",
		ClientSecret:      "",
		ClientRedirectUri: "",
		Credentials: CredentialsT{
			UserAccessToken:        "",
			UserAccessRefreshToken: "",
			UserAccessScope:        []string{},
			ExpiryUnixTimestamp:    0,
		},
	},
	TwitchVariables: TwitchVariablesT{
		StreamCategory: TwitchVariableT{
			Value:       "",
			Description: "Game or category currently being streamed",
		},
		StreamUptime: TwitchVariableT{
			Value:       "",
			Description: "Current stream duration, in seconds",
		},
		NumViewers: TwitchVariableT{
			Value:       "",
			Description: "Current number of viewers of the stream",
		},
		NumSubscribers: TwitchVariableT{
			Value:       "",
			Description: "Current number of subscribers to the channel",
		},
		NumFollowers: TwitchVariableT{
			Value:       "",
			Description: "Current number of followers of the channel",
		},
	},
	TwitchVariableUpdateIntervalSeconds: 10,
	LlmConfig:                           LlmConfigT{},
	AiGeneratedVariables:                []LlmVariableT{},
	Title: TitleT{
		TitleTemplate:              "",
		TitleUpdateIntervalMinutes: 1,
	},
	ActivityConsoleOutput: "",
}
