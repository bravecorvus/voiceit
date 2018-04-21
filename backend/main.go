package main

import (
	"net/http"

	"github.com/gilgameshskytrooper/voiceit/backend/app"
	"github.com/gilgameshskytrooper/voiceit/backend/utils"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

var (
	globals app.App
)

func init() {
	globals.Initialize()
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/login", globals.Login)
	r.HandleFunc("/register", globals.Register)
	r.HandleFunc("/secret/{username}", globals.Secret)
	// r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "../frontend/dist/index.html") })
	// r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("../frontend/dist/static/"))))
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, utils.Pwd()+"dist/index.html") })
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("dist/static/"))))
	n := negroni.Classic()
	n.UseHandler(r)

	n.Run(":8080")
}
