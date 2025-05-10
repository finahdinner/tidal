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
			Description: "Current stream duration, in minutes",
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
		MostRecentSubscriber: StreamVariableT{
			Value:       "",
			Description: "Username of the most recent subscriber to the channel",
		},
		MostRecentFollower: StreamVariableT{
			Value:       "",
			Description: "Username of the most recent follower of the channel",
		},
	},
	LlmVariables:           []LlmVariableT{},
	VariableUpdateInterval: 0,
	ActivityConsoleOutput:  "",
}
