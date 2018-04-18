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
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(utils.Pwd() + "dist/"))) // Serve Vue.js built assets at root
	n := negroni.Classic()
	n.UseHandler(r)
	err := http.ListenAndServe(":8080", n)
	if err != nil {
		panic(err.Error())
	}
}
