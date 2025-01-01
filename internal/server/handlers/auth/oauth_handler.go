package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/ilhamgepe/tengahin/internal/model"
	httpresponse "github.com/ilhamgepe/tengahin/pkg/httpResponse"
	"github.com/labstack/echo/v4"
)

const (
	googleState string = "tengahin_oauth_google"
	githubState string = "tengahin_oauth_github"
)

func (h *AuthHandler) GoogleLogin(c echo.Context) error {
	url := h.oauthProvider.Google.AuthCodeURL(googleState)
	// log.Info().Str("url", url).Msg("url")
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthHandler) GoogleCallback(c echo.Context) error {
	state := c.QueryParam("state")
	if state != googleState {
		return c.JSON(http.StatusBadRequest, httpresponse.RestError{
			ErrCauses: "invalid state",
		})
	}

	code := c.QueryParam("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, httpresponse.RestError{
			ErrCauses: "code not found",
		})
	}

	token, err := h.oauthProvider.Google.Exchange(c.Request().Context(), code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.RestError{
			ErrError:  err.Error(),
			ErrCauses: "failed to exchange token",
		})
	}

	res, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.RestError{
			ErrError:  err.Error(),
			ErrCauses: "failed to get user info",
		})
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return c.JSON(http.StatusInternalServerError, httpresponse.RestError{
			ErrError:  err.Error(),
			ErrCauses: "fetch user status code not ok",
		})
	}

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.RestError{
			ErrError:  err.Error(),
			ErrCauses: "failed to read response body",
		})
	}

	var userInfo model.GoogleUserInfo
	if err := json.Unmarshal(respBody, &userInfo); err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.RestError{
			ErrError:  err.Error(),
			ErrCauses: "failed to unmarshal response body",
		})
	}

	return c.JSON(http.StatusOK, httpresponse.RestSuccess{
		Status: http.StatusOK,
		Data:   userInfo,
	})
}

func (h *AuthHandler) GithubLogin(c echo.Context) error {
	url := h.oauthProvider.Github.AuthCodeURL(githubState)

	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthHandler) GithubCallback(c echo.Context) error {
	state := c.QueryParam("state")
	if state != githubState {
		return c.JSON(http.StatusBadRequest, httpresponse.RestError{
			ErrCauses: "invalid state",
		})
	}

	code := c.QueryParam("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, httpresponse.RestError{
			ErrCauses: "code not found",
		})
	}

	token, err := h.oauthProvider.Github.Exchange(c.Request().Context(), code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.RestError{
			ErrCauses: "failed to exchange token",
		})
	}

	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.RestError{
			ErrCauses: "failed to create request",
		})
	}

	req.Header.Set("Authorization", "bearer "+token.AccessToken)

	resp, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.RestError{
			ErrCauses: "failed to get user info",
		})
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.JSON(http.StatusInternalServerError, httpresponse.RestError{
			ErrCauses: "fetch user status code not ok",
		})
	}

	var userData model.GitHubUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userData); err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.RestError{
			ErrCauses: "failed to unmarshal response body",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"userData": userData,
	})
}
