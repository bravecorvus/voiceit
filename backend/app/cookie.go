package app

import (
	"net/http"

	"github.com/garyburd/redigo/redis"
	"github.com/gilgameshskytrooper/bigdisk/crypto"
)

func (app *App) setSession(username string, w http.ResponseWriter) {
	tokenval := crypto.GenerateRandomHash(20)
	app.DB.Do("HSET", "logins", username+":token", tokenval)
	value := map[string]string{
		"token": tokenval,
	}

	token, err := app.CookieHandler.Encode("token", value)
	if err == nil {
		cookie := &http.Cookie{
			Name:  "token",
			Value: token,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
	}
}

func (app *App) authenticateBrowserToken(username string, r *http.Request) bool {
	tokenfromcookie := ""
	if tokencookie, err := r.Cookie("token"); err == nil {
		tokencookieValue := make(map[string]string)
		if err = app.CookieHandler.Decode("token", tokencookie.Value, &tokencookieValue); err == nil {
			tokenfromcookie = tokencookieValue["token"]
		}
	} else {
		return false
	}
	tokenfromdb, _ := redis.String(app.DB.Do("HGET", "logins", username+":token"))
	if tokenfromcookie != tokenfromdb && tokenfromcookie != "" && tokenfromdb != "" {
		return false
	}
	return true
}
