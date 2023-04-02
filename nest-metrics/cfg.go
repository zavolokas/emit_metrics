package nestmetrics

type Config struct {
	ClientID     string
	ClientSecret string
	ProjectID    string
	AuthURL      string
	TokenURL     string
	Scopes       []string
	RedirectURL  string
	RefreshToken string
	AccessToken  string
	EmitFreqSec  int
}
