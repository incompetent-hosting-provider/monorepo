package authentication

type AccessToken string
type RefreshToken string

type SessionTokens struct {
	AccessToken  AccessToken  `json:"access_token,omitempty"`
	RefreshToken RefreshToken `json:"refresh_token,omitempty"`
}
