package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	REFRESH_TOKEN_URL = "https://chat.openai.com/api/auth/session"
)

type SessionResult struct {
	Err         string `json:"error"`
	Expires     string `json:"expires"`
	AccessToken string `json:"accessToken"`
}

type JwtToken struct {
	accessToken  string
	refreshToken string
	cfClearance  string
	mutex        sync.Mutex
	expiry       time.Time
}

// Creates a JwtToken instance that is expired and with a empty token value.
func NewEmptyJwtToken(refreshToken string, cfClearance string) JwtToken {
	return JwtToken{
		accessToken:  "",
		refreshToken: refreshToken,
		cfClearance:  cfClearance,
		expiry:       time.Now().Add(-10 * time.Second),
	}
}

func (j *JwtToken) Get() (string, error) {
	if (j.accessToken == "") || (time.Now().Before(j.expiry)) {
		_, err := j.refreshJwt()
		if err != nil {
			return "", err
		}
	}
	return j.accessToken, nil
}

func (j *JwtToken) set(accessToken string, expiry time.Time) string {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	j.accessToken = accessToken
	j.expiry = expiry
	return j.accessToken
}

func (j *JwtToken) refreshJwt() (string, error) {
	req, err := http.NewRequest("GET", REFRESH_TOKEN_URL, nil)
	if err != nil {
		return "", fmt.Errorf("[RefreshJwt] failed to create request: %v", err)
	}
	req.Header.Set("User-Agent", USER_AGENT)
	req.Header.Set("Cookie", fmt.Sprintf("cf_clearance=%s; __Secure-next-auth.session-token=%s", j.cfClearance, j.refreshToken))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("[RefreshJwt] failed to perform request: %v", err)
	}
	defer res.Body.Close()

	var result SessionResult
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return "", fmt.Errorf("[RefreshJwt] failed to decode response: %v", err)
	}

	accessToken := result.AccessToken
	if accessToken == "" {
		return "", errors.New("[RefreshJwt] unauthorized")
	}
	if result.Err != "" {
		if result.Err == "RefreshAccessTokenError" {
			return "", errors.New("[RefreshJwt] session token has expired")
		}

		return "", errors.New(result.Err)
	}

	expiryTime, err := time.Parse(time.RFC3339, result.Expires)
	if err != nil {
		return "", fmt.Errorf("[RefreshJwt] failed to parse expiry time: %v", err)
	}

	j.set(accessToken, expiryTime)
	return accessToken, nil
}
