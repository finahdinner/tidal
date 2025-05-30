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
	LlmConfig: LlmConfigT{
		DefaultPromptSuffix: "Ensure that your response only contains text relevant to the above information. " +
			"Do not exceed 100 characters. " +
			"Ensure your response does not contain profanities and cannot be construed as political or divisive.",
	},
	AiGeneratedVariables: []LlmVariableT{},
	Title: TitleT{
		Value:                           "",
		TitleTemplate:                   "",
		TitleUpdateIntervalMinutes:      1,
		UpdateImmediatelyOnStart:        true,
		ThrowErrorIfEmptyVariable:       true,
		ThrowErrorIfNonExistentVariable: true,
	},
	ActivityConsoleOutput: "",
}
