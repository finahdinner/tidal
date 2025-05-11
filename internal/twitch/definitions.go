package twitch

const (
	twitchApiAuthoriseUrl = "https://id.twitch.tv/oauth2/authorize"
	twitchApiTokenUrl     = "https://id.twitch.tv/oauth2/token"
)

type userAccessTokenInfoT struct {
	AccessToken  string   `json:"access_token"`
	ExpiresIn    int      `json:"expires_in"`
	RefreshToken string   `json:"refresh_token"`
	Scope        []string `json:"scope"`
	TokenType    string   `json:"token_type"`
}

const (
	twitchApiUsersUrl         = "https://api.twitch.tv/helix/users"
	twitchApiStreamsUrl       = "https://api.twitch.tv/helix/streams"
	twitchApiSubscriptionsUrl = "https://api.twitch.tv/helix/subscriptions"
	twitchApiFollowersUrl     = "https://api.twitch.tv/helix/channels/followers"
)

type getUsersApiResponseT struct {
	Data []struct {
		Id              string `json:"id"`
		Login           string `json:"login"`
		DisplayName     string `json:"display_name"`
		UserType        string `json:"type"`
		BroadcasterType string `json:"broadcaster_type"`
		Description     string `json:"description"`
		ProfileImageUrl string `json:"profile_image_url"`
		OfflineImageUrl string `json:"offline_image_url"`
		ViewCount       int    `json:"view_count"`
		Email           string `json:"email"`
		CreatedAt       string `json:"created_at"`
	} `json:"data"`
}

type paginationApiResponse struct {
	Cursor string `json:"cursor"`
}

type streamInfoT struct {
	Id           string   `json:"id"`
	UserId       string   `json:"user_id"`
	UserLogin    string   `json:"user_login"`
	UserName     string   `json:"user_name"`
	GameId       string   `json:"game_id"`
	GameName     string   `json:"game_name"`
	StreamType   string   `json:"type"`
	Title        string   `json:"title"`
	Tags         []string `json:"tags"`
	ViewerCount  int      `json:"viewer_count"`
	StartedAt    string   `json:"started_at"` // RFC3339 format
	Language     string   `json:"language"`
	ThumbnailUrl string   `json:"thumbnail_url"`
	TagIds       []string `json:"tag_ids"`
	IsMature     bool     `json:"is_mature"`
}

type getStreamInfoApiResponseT struct {
	Data       []streamInfoT         `json:"data"`
	Pagination paginationApiResponse `json:"pagination"`
}

type getChannelSubscribersResponseT struct {
	Data []struct {
		BroadcasterId    string `json:"broadcaster_id"`
		BroadcasterLogin string `json:"broadcaster_login"`
		BroadcasterName  string `json:"broadcaster_name"`
		GifterId         string `json:"gifter_id"`
		GifterLogin      string `json:"gifter_login"`
		GifterName       string `json:"gifter_name"`
		IsGift           bool   `json:"is_gift"`
		PlanName         string `json:"plan_name"`
		Tier             string `json:"tier"`
		UserId           string `json:"user_id"`
		UserName         string `json:"user_name"`
		UserLogin        string `json:"user_login"`
	} `json:"data"`
	Pagination paginationApiResponse `json:"pagination"`
	Points     int                   `json:"points"`
	Total      int                   `json:"total"`
}

type getChannelFollowersResponseT struct {
	Data []struct {
		FollowedAt string `json:"followed_at"`
		UserID     string `json:"user_id"`
		UserLogin  string `json:"user_login"`
		UserName   string `json:"user_name"`
	} `json:"data"`
	Pagination paginationApiResponse `json:"pagination"`
	Total      int                   `json:"total"`
}

type RawApiResponses struct {
	StreamInfo      *streamInfoT
	SubscribersInfo *getChannelSubscribersResponseT
	FollowersInfo   *getChannelFollowersResponseT
}
