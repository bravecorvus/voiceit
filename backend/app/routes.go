package app

import (
	"fmt"
	"net/http"
)

func (app *App) Root(w http.ResponseWriter, r *http.Request) {
	app.setInitialSession(w, r)
}

func (app *App) Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println(getTokenFromRequest(r))
}
