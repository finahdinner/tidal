package preferences

var DefaultPreferences PreferencesFormat = PreferencesFormat{
	TwitchConfig: twitchConfigT{
		UserName:     "",
		UserId:       "",
		ClientId:     "",
		ClientSecret: "",
		RedirectUri:  "",
		Credentials: credentialsT{
			UserAccessToken:        "",
			UserAccessRefreshToken: "",
			UserAccessScope:        []string{},
		},
	},
	StreamVariables: streamVariablesT{
		StreamCategory: streamVariableT{
			Value:       "",
			Description: "Game or category currently being streamed",
		},
		StreamUptime: streamVariableT{
			Value:       "",
			Description: "Current stream duration, in minutes",
		},
		NumViewers: streamVariableT{
			Value:       "",
			Description: "Current number of viewers of the stream",
		},
		NumSubscribers: streamVariableT{
			Value:       "",
			Description: "Current number of subscribers to the channel",
		},
		NumFollowers: streamVariableT{
			Value:       "",
			Description: "Current number of followers of the channel",
		},
		MostRecentSubscriber: streamVariableT{
			Value:       "",
			Description: "Username of the most recent subscriber to the channel",
		},
		MostRecentFollower: streamVariableT{
			Value:       "",
			Description: "Username of the most recent follower of the channel",
		},
	},
	LlmVariables:          []llmVariableT{},
	ActivityConsoleOutput: "",
	UpdateFrequency:       0,
}
