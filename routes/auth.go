package routes

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/issue-notifier/issue-notifier-api/models"
	"github.com/issue-notifier/issue-notifier-api/session"
	"github.com/issue-notifier/issue-notifier-api/utils"
)

type UserInfo struct {
	Name        string `json:"name"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	AvatarImg   string `json:"avatarImg"`
	AccessToken string `json:"accessToken"`
}

func GetAuthenticatedUser(w http.ResponseWriter, r *http.Request) {
	ses, err := session.Store.Get(r, session.CookieName)

	if err != nil {
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

func GitHubLogin(w http.ResponseWriter, r *http.Request) {
	ses, err := session.Store.Get(r, session.CookieName)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	code := r.URL.Query()["code"][0]
	accessToken, err := getAccessToken(code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	userInfo, err := getUserInfo(accessToken)
	if err != nil {
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If no other error returned then save session and respond back
	userSession.UserID = userID

	// Set user as authenticated
	ses.Values["UserSession"] = userSession
	ses.Save(r, w)

	utils.RespondWithJSON(w, http.StatusOK, userInfo)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	ses, err := session.Store.Get(r, session.CookieName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ses.Values["UserSession"] = session.UserSession{}
	ses.Options.MaxAge = -1

	err = ses.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, "Success")
}

func getAccessToken(code string) (string, error) {
	reqBody, _ := json.Marshal(map[string]string{
		"client_id":     GITHUB_CLIENT_ID,
		"client_secret": GITHUB_CLIENT_SECRET,
		"code":          code,
	})

	httpClient := &http.Client{}
	req, _ := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer(reqBody))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	res, err := httpClient.Do(req)

	if err != nil {
		log.Fatalln(err)
	} else {
		defer res.Body.Close()
		data, _ := ioutil.ReadAll(res.Body)

		var loginResponse map[string]interface{}
		json.Unmarshal(data, &loginResponse)

		return loginResponse["access_token"].(string), nil
	}

	return "", err
}

func getUserInfo(accessToken string) (*UserInfo, error) {
	httpClient := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.github.com/user?access_token="+accessToken, nil)
	// req.Header.Set("Authorization", accessToken) // FIXME: See why this is not working

	res, err := httpClient.Do(req)

	if err != nil {
		log.Fatalln(err)
	} else {
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

	return nil, err
}
