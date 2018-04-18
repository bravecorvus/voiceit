package app

import (
	"net/http"
)

func (app *App) Middleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	next(w, r)
}
