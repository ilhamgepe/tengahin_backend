package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

type GitHubUserInfo struct {
	AvatarURL         string    `json:"avatar_url"`
	Bio               *string   `json:"bio"`
	Blog              string    `json:"blog"`
	Company           *string   `json:"company"`
	CreatedAt         time.Time `json:"created_at"`
	Email             string    `json:"email"`
	EventsURL         string    `json:"events_url"`
	Followers         int       `json:"followers"`
	FollowersURL      string    `json:"followers_url"`
	Following         int       `json:"following"`
	FollowingURL      string    `json:"following_url"`
	GistsURL          string    `json:"gists_url"`
	GravatarID        string    `json:"gravatar_id"`
	Hireable          *bool     `json:"hireable"`
	HTMLURL           string    `json:"html_url"`
	ID                int       `json:"id"`
	Location          *string   `json:"location"`
	Login             string    `json:"login"`
	Name              *string   `json:"name"`
	NodeID            string    `json:"node_id"`
	NotificationEmail string    `json:"notification_email"`
	OrganizationsURL  string    `json:"organizations_url"`
	PublicGists       int       `json:"public_gists"`
	PublicRepos       int       `json:"public_repos"`
	ReceivedEventsURL string    `json:"received_events_url"`
	ReposURL          string    `json:"repos_url"`
	SiteAdmin         bool      `json:"site_admin"`
	StarredURL        string    `json:"starred_url"`
	SubscriptionsURL  string    `json:"subscriptions_url"`
	TwitterUsername   *string   `json:"twitter_username"`
	Type              string    `json:"type"`
	UpdatedAt         time.Time `json:"updated_at"`
	URL               string    `json:"url"`
	UserViewType      string    `json:"user_view_type"`
}

type RegisterDTO struct {
	Email    string `json:"email" form:"email" validate:"required,email" db:"email"`
	Username string `json:"username" form:"username" validate:"required" db:"username"`
	Fullname string `json:"fullname" form:"fullname" validate:"required,min=3" db:"fullname"`
	Password string `json:"password" form:"password" validate:"required,min=6" db:"password"`
}

type LoginDTO struct {
	Email    string `json:"email" form:"email" validate:"required,email" db:"email"`
	Password string `json:"password" form:"password" validate:"required,min=6" db:"password"`
}

type RefreshTokenDTO struct {
	RefreshToken string `json:"refresh_token" form:"refresh_token" validate:"required"`
}

func (u *RegisterDTO) HashPassword() error {
	byte, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(byte)
	return nil
}
