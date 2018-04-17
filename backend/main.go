package main

import (
	"net/http"

	"github.com/gilgameshskytrooper/voiceit/backend/app"
	"github.com/gilgameshskytrooper/voiceit/backend/utils"
	"github.com/gorilla/mux"
)

var (
	globals app.App
)

func init() {
	globals.Initialize()
}
func main() {
	r := mux.NewRouter()
	http.Handle("/",
		http.FileServer(http.Dir(utils.Pwd()+"dist")))
	http.ListenAndServe(":8080", r)
}
