package main

import (
	"net/http"

	"github.com/gilgameshskytrooper/voiceit/backend/app"
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
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("../frontend/dist/")))
	// r.PathPrefix("/").Handler(http.FileServer(http.Dir(utils.Pwd() + "dist/")))
	n := negroni.Classic()
	n.UseHandler(r)

	n.Run(":8080")
}
