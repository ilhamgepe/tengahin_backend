package oauth

import (
	"github.com/ilhamgepe/tengahin/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

type OauthProviders struct {
	Google *oauth2.Config
	Github *oauth2.Config
}

func NewOauthProviders(cfg *config.Config) *OauthProviders {
	return &OauthProviders{
		Google: &oauth2.Config{
			RedirectURL:  cfg.Oauth.GoogleCallbackURL,
			ClientID:     cfg.Oauth.GoogleClientID,
			ClientSecret: cfg.Oauth.GoogleClientSecret,
			Endpoint:     google.Endpoint,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.profile",
			},
		},
		Github: &oauth2.Config{
			RedirectURL:  cfg.Oauth.GithubCallbackURL,
			ClientID:     cfg.Oauth.GithubClientID,
			ClientSecret: cfg.Oauth.GithubClientSecret,
			Endpoint:     github.Endpoint,
			Scopes: []string{
				"read:user",
			},
		},
	}
}
