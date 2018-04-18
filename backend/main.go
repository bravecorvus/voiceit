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
	rootmux := mux.NewRouter()

	// Works, but no cookie writing
	// r.PathPrefix("/").Handler(http.FileServer(http.Dir(utils.Pwd() + "dist/")))

	// r.PathPrefix("/").Handler(negroni.New(
	// globals.Middleware,
	// negroni.Wrap(r),
	// ))

	rootmux.HandleFunc("/", globals.Root)

	n := negroni.Classic()
	// n.Use(negroni.NewStatic(http.Dir(utils.Pwd() + "/dist")))
	n.Use(negroni.NewStatic(http.Dir("../frontend/dist")))
	n.UseHandler(rootmux)

	mux := mux.NewRouter()
	mux.HandleFunc("/login", globals.Login)
	n.UseHandler(mux)

	n.Run(":8080")
}
