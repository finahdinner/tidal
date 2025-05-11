package twitch

const (
	twitchApiAuthoriseUrl     = "https://id.twitch.tv/oauth2/authorize"
	twitchApiTokenUrl         = "https://id.twitch.tv/oauth2/token"
	twitchApiChannelUrl       = "https://api.twitch.tv/helix/channels"
	twitchApiUsersUrl         = "https://api.twitch.tv/helix/users"
	twitchApiSubscriptionsUrl = "https://api.twitch.tv/helix/subscriptions"
)

type paginationApiResponse struct {
	Cursor string `json:"cursor"`
}

type userAccessTokenInfoT struct {
	AccessToken  string   `json:"access_token"`
	ExpiresIn    int      `json:"expires_in"`
	RefreshToken string   `json:"refresh_token"`
	Scope        []string `json:"scope"`
	TokenType    string   `json:"token_type"`
}

// type userDataT struct {
// 	Id              string `json:"id"`
// 	Login           string `json:"login"`
// 	DisplayName     string `json:"display_name"`
// 	UserType        string `json:"type"`
// 	BroadcasterType string `json:"broadcaster_type"`
// 	Description     string `json:"description"`
// 	ProfileImageUrl string `json:"profile_image_url"`
// 	OfflineImageUrl string `json:"offline_image_url"`
// 	ViewCount       int    `json:"view_count"`
// 	Email           string `json:"email"`
// 	CreatedAt       string `json:"created_at"`
// }

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

type streamInfoApiResponseT struct {
	id           string   `json:"id"`
	userId       string   `json:"user_id"`
	userLogin    string   `json:"user_login"`
	userName     string   `json:"user_name"`
	gameId       string   `json:"game_id"`
	gameName     string   `json:"game_name"`
	streamType   string   `json:"type"`
	title        string   `json:"title"`
	tags         []string `json:"tags"`
	viewerCount  int      `json:"viewer_count"`
	startedAt    string   `json:"started_at"`
	language     string   `json:"language"`
	thumbnailUrl string   `json:"thumbnail_url"`
	tagIds       []string `json:"tag_ids"`
	isMature     bool     `json:"is_mature"`
}

type getStreamsApiResponseT struct {
	data       []streamInfoApiResponseT
	pagination paginationApiResponse
}
