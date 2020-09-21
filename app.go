package main

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)

type App struct {
	NoAuthRouter *mux.Router
	AuthRouter   *mux.Router
	DB           *sql.DB
}

type User struct {
	Name      string `json:"name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	AvatarImg string `json:"avatarImg"`
}

type UserSession struct {
	Username        string `json:"username"`
	AccessToken     string `json:"accessToken"`
	IsAuthenticated bool   `json:"isAuthenticated"`
}

type Response struct {
	Success bool        `json:"success"`
	Error   string      `json:"error"`
	Payload interface{} `json:"payload"`
}

// store will hold all session data
var store *sessions.CookieStore // TODO: Move to token based authentication in future
var cookieName string = "cookie-name"

func (a *App) Initialize(user, password, dbname string) {
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)
	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	a.NoAuthRouter = mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	a.AuthRouter = a.NoAuthRouter.PathPrefix("/user").Subrouter()
	a.AuthRouter.Use(a.isAuthenticated)

	a.initializeRoutes()

	authKey := securecookie.GenerateRandomKey(64)
	encryptionKey := securecookie.GenerateRandomKey(32)

	store = sessions.NewCookieStore(
		authKey,
		encryptionKey,
	)

	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 60 * 15,
		HttpOnly: true,
	}

	gob.Register(UserSession{})
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.NoAuthRouter))
}

func (a *App) initializeRoutes() {
	// a.Router.HandleFunc("/products", a.getProducts).Methods("GET")
	// a.Router.HandleFunc("/product", a.createProduct).Methods("POST")
	// a.Router.HandleFunc("/product/{id:[0-9]+}", a.getProduct).Methods("GET")
	// a.Router.HandleFunc("/product/{id:[0-9]+}", a.updateProduct).Methods("PUT")
	// a.Router.HandleFunc("/product/{id:[0-9]+}", a.deleteProduct).Methods("DELETE")

	// matches /api/v1/...
	a.NoAuthRouter.HandleFunc("/login/github/oauth2", a.gitHubLogin).Methods("GET")

	// matches /api/v1/user/...
	a.AuthRouter.HandleFunc("/authenticated", a.getAuthenticatedUser).Methods("GET")
	a.AuthRouter.HandleFunc("/logout", a.logout).Methods("GET")
}

func (a *App) isAuthenticated(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, cookieName)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		} else {
			// Check if user is Authenticated
			userSession, ok := session.Values["UserSession"].(UserSession)

			if !userSession.IsAuthenticated || !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			} else {
				next.ServeHTTP(w, r)
			}
		}
	})
}

func getUserSession(session *sessions.Session) (userSession UserSession) {
	userSession, _ = session.Values["UserSession"].(UserSession)

	return
}

func (a *App) getAuthenticatedUser(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, cookieName)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	} else {
		accessToken := getUserSession(session).AccessToken
		user, err := getUser(accessToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		} else {
			response := createResponse(true, "", user)

			respondWithJSON(w, http.StatusOK, response)
		}
	}
}

func (a *App) gitHubLogin(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, cookieName)

	code := r.URL.Query()["code"][0]

	accessToken, err := getAccessToken(code)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	} else {
		user, err := getUser(accessToken)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		} else {
			userSession := &UserSession{
				Username:        user.Username,
				AccessToken:     accessToken,
				IsAuthenticated: true,
			}

			response := createResponse(true, "", user)

			// Set user as authenticated
			session.Values["UserSession"] = userSession
			session.Save(r, w)

			respondWithJSON(w, http.StatusOK, response)
		}
	}
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

func getAccessToken(code string) (string, error) {
	reqBody, _ := json.Marshal(map[string]string{
		"client_id":     "69c3fc731ccb2d116412",                     // TODO: Move to env vars
		"client_secret": "df34a3cadf00452757713d71404eef86731af668", // TODO: Move to env vars
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

func (a *App) logout(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, cookieName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["UserSession"] = UserSession{}
	session.Options.MaxAge = -1

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := createResponse(true, "", nil)

	respondWithJSON(w, http.StatusOK, response)
}

// func (a *App) getProducts(w http.ResponseWriter, r *http.Request) {
// 	// count, _ := strconv.Atoi(r.FormValue("count"))
// 	// start, _ := strconv.Atoi(r.FormValue("start"))

// 	// if count > 10 || count < 1 {
// 	count := 10
// 	// }
// 	// if start < 0 {
// 	start := 0
// 	// }

// 	products, err := getProducts(a.DB, start, count)
// 	if err != nil {
// 		// respondWithError(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	// respondWithJSON(w, http.StatusOK, products)
// }

// func (a *App) getProduct(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	id, err := strconv.Atoi(vars["id"])
// 	if err != nil {
// 		// respondWithError(w, http.StatusBadRequest, "Invalid product ID")
// 		return
// 	}

// 	p := product{ID: id}
// 	if err := p.getProduct(a.DB); err != nil {
// 		switch err {
// 		case sql.ErrNoRows:
// 			// respondWithError(w, http.StatusNotFound, "Product not found")
// 		default:
// 			// respondWithError(w, http.StatusInternalServerError, err.Error())
// 		}
// 		return
// 	}

// 	// respondWithJSON(w, http.StatusOK, p)
// }

// func (a *App) createProduct(w http.ResponseWriter, r *http.Request) {
// 	var p product
// 	decoder := json.NewDecoder(r.Body)
// 	if err := decoder.Decode(&p); err != nil {
// 		// respondWithError(w, http.StatusBadRequest, "Invalid request payload")
// 		return
// 	}
// 	defer r.Body.Close()

// 	if err := p.createProduct(a.DB); err != nil {
// 		// respondWithError(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	// respondWithJSON(w, http.StatusCreated, p)
// }

// func (a *App) updateProduct(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	id, err := strconv.Atoi(vars["id"])
// 	if err != nil {
// 		// respondWithError(w, http.StatusBadRequest, "Invalid product ID")
// 		return
// 	}

// 	var p product
// 	decoder := json.NewDecoder(r.Body)
// 	if err := decoder.Decode(&p); err != nil {
// 		// respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
// 		return
// 	}
// 	defer r.Body.Close()
// 	p.ID = id

// 	if err := p.updateProduct(a.DB); err != nil {
// 		// respondWithError(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	// respondWithJSON(w, http.StatusOK, p)
// }

// func (a *App) deleteProduct(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	id, err := strconv.Atoi(vars["id"])
// 	if err != nil {
// 		// respondWithError(w, http.StatusBadRequest, "Invalid Product ID")
// 		return
// 	}

// 	p := product{ID: id}
// 	if err := p.deleteProduct(a.DB); err != nil {
// 		// respondWithError(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	// respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
// }

func createResponse(success bool, message string, payload interface{}) (response *Response) {
	response = &Response{
		Success: success,
		Error:   message,
		Payload: payload,
	}

	return
}

func respondWithJSON(w http.ResponseWriter, code int, response *Response) {
	res, _ := json.Marshal(response)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(res)
}
