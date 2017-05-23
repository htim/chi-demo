package main

import (
	"net/http"
	"github.com/pressly/chi"
	"strconv"
	"fmt"
	"golang.org/x/oauth2"
	"context"
	"github.com/google/go-github/github"
	"github.com/dgrijalva/jwt-go"

)

func (a *App) GetLink(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		fmt.Println(err)
	}
	link, err := GetLinkById(a.DB, id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	respondWithJson(w, http.StatusOK, link)
}

func (a *App) GitHubLogin(w http.ResponseWriter, r *http.Request) {
	url := a.GithubOauth2Config.AuthCodeURL(a.Oauth2StateString, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (a *App) GitHubLoginCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != a.Oauth2StateString {
		fmt.Println("Invalid state: " + state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	code := r.FormValue("code")
	ctx := context.Background()
	token, err := a.GithubOauth2Config.Exchange(ctx, code)
	if err != nil {
		fmt.Println("error while getting token: ", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	oauthClient := a.GithubOauth2Config.Client(ctx, token)
	client := github.NewClient(oauthClient)
	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		fmt.Println("error while getting user info", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	fmt.Println(user)

	claims := jwt.MapClaims{
		"login":user.Login,
		"id": user.ID,
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := jwtToken.SignedString(a.SigningKey)
	if err != nil {
		fmt.Println(err)
	}
	w.Write([]byte(tokenString))
}
