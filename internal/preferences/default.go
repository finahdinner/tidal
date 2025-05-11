package preferences

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
		},
	},
	StreamVariables: StreamVariablesT{
		StreamCategory: StreamVariableT{
			Value:       "",
			Description: "Game or category currently being streamed",
		},
		StreamUptime: StreamVariableT{
			Value:       "",
			Description: "Current stream duration, in seconds",
		},
		NumViewers: StreamVariableT{
			Value:       "",
			Description: "Current number of viewers of the stream",
		},
		NumSubscribers: StreamVariableT{
			Value:       "",
			Description: "Current number of subscribers to the channel",
		},
		NumFollowers: StreamVariableT{
			Value:       "",
			Description: "Current number of followers of the channel",
		},
	},
	LlmVariables:           []LlmVariableT{},
	VariableUpdateInterval: -1,
	ActivityConsoleOutput:  "",
}
