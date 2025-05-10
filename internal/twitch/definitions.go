package twitch

const (
	twitchApiAuthoriseUrl = "https://id.twitch.tv/oauth2/authorize"
	twitchApiTokenUrl     = "https://id.twitch.tv/oauth2/token"
	twitchApiChannelUrl   = "https://api.twitch.tv/helix/channels"
	twitchApiUsersUrl     = "https://api.twitch.tv/helix/users"
)

type userAccessTokenInfoT struct {
	AccessToken  string   `json:"access_token"`
	ExpiresIn    int      `json:"expires_in"`
	RefreshToken string   `json:"refresh_token"`
	Scope        []string `json:"scope"`
	TokenType    string   `json:"token_type"`
}

type userDataT struct {
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
}

type getUsersApiResponseT struct {
	Data []userDataT
}
