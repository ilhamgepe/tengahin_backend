package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
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

	// check user in db if exist
	user, err := h.userService.FindByEmail(c.Request().Context(), userInfo.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return httpresponse.KnownSQLError(c, err)
	}

	// if not exist, create new user
	if errors.Is(err, sql.ErrNoRows) {
		user, err = h.userService.CreateUser(c.Request().Context(), model.RegisterDTO{
			Email:    userInfo.Email,
			Username: userInfo.GivenName,
			Fullname: userInfo.Name,
		})
		if err != nil {
			return httpresponse.KnownSQLError(c, err)
		}
		user.Sanitize()
	}

	// generate token
	accessToken, payload, err := h.tokenMaker.CreateToken(user.ID, h.cfg.Server.TokenDuration)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.RestError{
			ErrError:  echo.ErrInternalServerError.Error(),
			ErrCauses: err.Error(),
		})
	}

	refreshToken, payloadRefresh, err := h.tokenMaker.CreateRefreshToken(user.ID, h.cfg.Server.RefreshTokenDuration)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.RestError{
			ErrError:  echo.ErrInternalServerError.Error(),
			ErrCauses: err.Error(),
		})
	}

	userJSON, err := json.Marshal(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.RestError{
			ErrError:  echo.ErrInternalServerError.Error(),
			ErrCauses: err.Error(),
		})
	}

	expiresIn := time.Until(payload.ExpiresAt.Time)
	rediRes := h.rdb.Set(c.Request().Context(), payload.ID, userJSON, expiresIn)
	if rediRes.Err() != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.RestError{
			ErrError:  echo.ErrInternalServerError.Error(),
			ErrCauses: rediRes.Err().Error(),
		})
	}

	expiresIn = time.Until(payloadRefresh.ExpiresAt.Time)
	redisRefreshRes := h.rdb.Set(c.Request().Context(), payloadRefresh.ID, refreshToken, expiresIn)
	if redisRefreshRes.Err() != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.RestError{
			ErrError:  echo.ErrInternalServerError.Error(),
			ErrCauses: redisRefreshRes.Err().Error(),
		})
	}

	return c.JSON(200, httpresponse.RestSuccess{
		Status: http.StatusOK,
		Data: map[string]interface{}{
			"access_token":             accessToken,
			"expires_at":               payload.ExpiresAt,
			"refresh_token":            refreshToken,
			"refresh_token_expires_at": payloadRefresh.ExpiresAt,
			"user":                     user,
		},
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

	var userInfo model.GitHubUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.RestError{
			ErrCauses: "failed to unmarshal response body",
		})
	}

	// check user in db if exist
	user, err := h.userService.FindByEmail(c.Request().Context(), userInfo.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return httpresponse.KnownSQLError(c, err)
	}

	// if not exist, create new user
	if errors.Is(err, sql.ErrNoRows) {
		user, err = h.userService.CreateUser(c.Request().Context(), model.RegisterDTO{
			Email:    userInfo.Email,
			Username: userInfo.Login,
			Fullname: *userInfo.Name,
		})
		if err != nil {
			return httpresponse.KnownSQLError(c, err)
		}
		user.Sanitize()
	}

	// generate token
	accessToken, payload, err := h.tokenMaker.CreateToken(user.ID, h.cfg.Server.TokenDuration)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.RestError{
			ErrError:  echo.ErrInternalServerError.Error(),
			ErrCauses: err.Error(),
		})
	}

	refreshToken, payloadRefresh, err := h.tokenMaker.CreateRefreshToken(user.ID, h.cfg.Server.RefreshTokenDuration)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.RestError{
			ErrError:  echo.ErrInternalServerError.Error(),
			ErrCauses: err.Error(),
		})
	}

	userJSON, err := json.Marshal(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.RestError{
			ErrError:  echo.ErrInternalServerError.Error(),
			ErrCauses: err.Error(),
		})
	}

	expiresIn := time.Until(payload.ExpiresAt.Time)
	rediRes := h.rdb.Set(c.Request().Context(), payload.ID, userJSON, expiresIn)
	if rediRes.Err() != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.RestError{
			ErrError:  echo.ErrInternalServerError.Error(),
			ErrCauses: rediRes.Err().Error(),
		})
	}

	expiresIn = time.Until(payloadRefresh.ExpiresAt.Time)
	redisRefreshRes := h.rdb.Set(c.Request().Context(), payloadRefresh.ID, refreshToken, expiresIn)
	if redisRefreshRes.Err() != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.RestError{
			ErrError:  echo.ErrInternalServerError.Error(),
			ErrCauses: redisRefreshRes.Err().Error(),
		})
	}

	return c.JSON(200, httpresponse.RestSuccess{
		Status: http.StatusOK,
		Data: map[string]interface{}{
			"access_token":             accessToken,
			"expires_at":               payload.ExpiresAt,
			"refresh_token":            refreshToken,
			"refresh_token_expires_at": payloadRefresh.ExpiresAt,
			"user":                     user,
		},
	})
}
