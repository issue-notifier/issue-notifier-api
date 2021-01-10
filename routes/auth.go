package routes

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/issue-notifier/issue-notifier-api/models"
	"github.com/issue-notifier/issue-notifier-api/session"
	"github.com/issue-notifier/issue-notifier-api/utils"
)

// UserInfo struct to store user related information as received from Github. This is passed to UI for display
type UserInfo struct {
	Name        string `json:"name"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	AvatarImg   string `json:"avatarImg"`
	AccessToken string `json:"accessToken"`
}

// GetAuthenticatedUser godoc
// @Summary Get user information from Github for authenticated user
// @Description Get user information from Github for authenticated user via session token
// @Tags user
// @Produce json
// @Security Github OAuth
// @Success 200 {object} UserInfo
// @Failure 401 {string} Unauthorized
// @Router /api/v1/user/authenticated [get]
func GetAuthenticatedUser(w http.ResponseWriter, r *http.Request) {
	ses, err := session.Store.Get(r, session.CookieName)

	if err != nil {
		utils.LogError.Println("Failed to get user session. Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userSession, _ := ses.Values["UserSession"].(session.UserSession)
	accessToken := userSession.AccessToken
	userInfo, err := getUserInfo(accessToken)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, userInfo)
}

// GitHubLogin godoc
// @Summary Login via Github OAuth
// @Description Performs a Github OAuth by using the `code` provided in query param and then uses the received accessToken to fetch user information from Github. If it is a new user then user is also saved in the database. It also creates a user session for future authenticated calls.
// @Tags user
// @Produce json
// @Success 200 {object} UserInfo
// @Failure 401 {string} Unauthorized
// @Failure 500 {string} Internal Server Error
// @Router /api/v1/login/github/oauth2 [get]
func GitHubLogin(w http.ResponseWriter, r *http.Request) {
	ses, err := session.Store.Get(r, session.CookieName)
	if err != nil {
		utils.LogError.Println("Failed to get session information. Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	code := r.URL.Query()["code"][0]
	accessToken, err := getAccessToken(code)
	if err != nil {
		utils.LogError.Println("Failed to get access token from Github. Error:", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	userInfo, err := getUserInfo(accessToken)
	if err != nil {
		utils.LogError.Println("Failed to get user information from Github. Error:", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	userSession := &session.UserSession{
		AccessToken:     accessToken,
		IsAuthenticated: true,
	}

	var userID string
	userID, err = models.GetUserIDByUsername(userInfo.Username)
	// If no users found with the given username, create the user in DB
	if err == sql.ErrNoRows {
		userID, err = models.CreateUser(userInfo.Username, userInfo.Email)
		if err != nil {
			utils.LogError.Println("Failed to create new user. Error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		utils.LogInfo.Println("Successfully created new user with userID:", userID)

	} else if err != nil {
		utils.LogError.Println("Failed to get username for userID:", userID)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If no other error returned then save session and respond back
	userSession.UserID = userID

	// Set user as authenticated
	ses.Values["UserSession"] = userSession
	ses.Save(r, w)

	utils.LogInfo.Println("Successfully logged in user")
	utils.RespondWithJSON(w, http.StatusOK, userInfo)
}

// Logout godoc
// @Summary Logout from user session
// @Description Performs the logout function by deleting user session
// @Tags user
// @Security Github OAuth
// @Success 200 {string} Success
// @Failure 500 {string} Internal Server Error
// @Router /api/v1/user/logout [get]
func Logout(w http.ResponseWriter, r *http.Request) {
	ses, err := session.Store.Get(r, session.CookieName)
	if err != nil {
		utils.LogError.Println("Failed to get session information. Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ses.Values["UserSession"] = session.UserSession{}
	ses.Options.MaxAge = -1

	err = ses.Save(r, w)
	if err != nil {
		utils.LogError.Println("Failed save session. Unabled to logout. Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.LogInfo.Println("Successfully logged out user")
	utils.RespondWithJSON(w, http.StatusOK, "Success")
}

func getAccessToken(code string) (string, error) {
	reqBody, _ := json.Marshal(map[string]string{
		"client_id":     GithubClientID,
		"client_secret": GithubClientSecret,
		"code":          code,
	})

	httpClient := &http.Client{}
	req, _ := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer(reqBody))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	res, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("[getAccessToken] %v", err)
	}
	defer res.Body.Close()

	data, _ := ioutil.ReadAll(res.Body)

	var loginResponse map[string]interface{}
	json.Unmarshal(data, &loginResponse)

	return loginResponse["access_token"].(string), nil
}

func getUserInfo(accessToken string) (*UserInfo, error) {
	httpClient := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	req.Header.Set("Authorization", "token "+accessToken)

	res, err := httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("[getUserInfo] %v", err)
	}
	defer res.Body.Close()

	dataBytes, _ := ioutil.ReadAll(res.Body)

	var data map[string]interface{}
	json.Unmarshal(dataBytes, &data)

	userInfo := &UserInfo{
		Name:        data["name"].(string),
		Username:    data["login"].(string),
		Email:       data["email"].(string),
		AvatarImg:   data["avatar_url"].(string),
		AccessToken: accessToken,
	}

	return userInfo, nil
}
