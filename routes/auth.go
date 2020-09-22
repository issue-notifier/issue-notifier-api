package routes

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/hemakshis/issue-notifier-api/session"
	"github.com/hemakshis/issue-notifier-api/utils"
)

type User struct {
	Name      string `json:"name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	AvatarImg string `json:"avatarImg"`
}

func GetAuthenticatedUser(w http.ResponseWriter, r *http.Request) {
	ses, err := session.Store.Get(r, session.CookieName)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		userSession, _ := ses.Values["UserSession"].(session.UserSession)

		accessToken := userSession.AccessToken
		user, err := getUser(accessToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		} else {
			response := utils.CreateResponse(true, "", user)

			utils.RespondWithJSON(w, http.StatusOK, response)
		}

	}
}

func GitHubLogin(w http.ResponseWriter, r *http.Request) {
	ses, err := session.Store.Get(r, session.CookieName)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		code := r.URL.Query()["code"][0]

		accessToken, err := getAccessToken(code)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		} else {
			user, err := getUser(accessToken)

			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
			} else {
				userSession := &session.UserSession{
					Username:        user.Username,
					AccessToken:     accessToken,
					IsAuthenticated: true,
				}

				response := utils.CreateResponse(true, "", user)

				// Set user as authenticated
				ses.Values["UserSession"] = userSession
				ses.Save(r, w)

				utils.RespondWithJSON(w, http.StatusOK, response)
			}
		}
	}

}

func Logout(w http.ResponseWriter, r *http.Request) {
	ses, err := session.Store.Get(r, session.CookieName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		ses.Values["UserSession"] = session.UserSession{}
		ses.Options.MaxAge = -1

		err = ses.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := utils.CreateResponse(true, "", nil)

		utils.RespondWithJSON(w, http.StatusOK, response)
	}

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

func getUser(accessToken string) (*User, error) {

	httpClient := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.github.com/user?access_token="+accessToken, nil)
	// req.Header.Set("Authorization", accessToken) // FIXME: See why this is not working

	res, err := httpClient.Do(req)

	if err != nil {
		log.Fatalln(err)
	} else {
		defer res.Body.Close()

		data, _ := ioutil.ReadAll(res.Body)

		var userInfo map[string]interface{}
		json.Unmarshal(data, &userInfo)

		user := &User{
			Name:      userInfo["name"].(string),
			Username:  userInfo["login"].(string),
			Email:     userInfo["email"].(string),
			AvatarImg: userInfo["avatar_url"].(string),
		}

		return user, nil
	}

	return nil, err
}
