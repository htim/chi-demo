package main

import (
	"github.com/pressly/chi"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"net/http"
	"encoding/json"
	"fmt"
	"log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"os"
)

type App struct {
	Router             *chi.Mux
	DB                 *sqlx.DB
	GithubOauth2Config *oauth2.Config
	Oauth2StateString  string
	SigningKey         []byte
}

func (a *App) Initialize() {
	var err error
	connectionString := os.Getenv("POSTGRES_URL")
	a.DB, err = sqlx.Open("postgres", connectionString)
	if err != nil {
		fmt.Println(err)
	}
	a.Router = chi.NewRouter()
	a.initRoutes()
	a.GithubOauth2Config, a.Oauth2StateString = configureGitHubOauth2()
	a.SigningKey = []byte(os.Getenv("SIGNING_KEY"))
}

func (a *App) Run(port string) {
	log.Fatal(http.ListenAndServe(port, a.Router))
}

func (a *App) initRoutes() {
	a.Router.Route("/links", func(r chi.Router) {
		r.Get("/:id", a.GetLink)
	})
	a.Router.Get("/github-login", a.GitHubLogin)
	a.Router.Get("/github-callback", a.GitHubLoginCallback)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJson(w, code, map[string]string{"error": message})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println(err)
	}
	w.Write(response)
}

func configureGitHubOauth2() (*oauth2.Config, string) {
	conf := &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		Scopes:       []string{"user"},
		Endpoint:     github.Endpoint,
	}
	return conf, "randomString"
}
