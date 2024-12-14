package lib

const (
	OAuthStateKey          = "oauth.state"
	OAuthSessionName       = "oauth.session"
	OAuthStateToken        = "oauth.token"
	OAuthCloudId           = "oauth.cloud.id"
	ProfileSessionName     = "profile"
	ProfileUserDisplayName = "profile.user.name"
	ProfileUserImage       = "profile.user.image"

	AuthResults = "auth.results"
)

type AuthSession struct {
	CloudId      string
	Token        string
	DisplayName  string
	ProfileImage string
}
