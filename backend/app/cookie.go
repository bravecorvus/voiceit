package app

import (
	"log"
	"net/http"

	"github.com/garyburd/redigo/redis"
	"github.com/gilgameshskytrooper/bigdisk/crypto"
	"golang.org/x/crypto/bcrypt"
)

func (app *App) setSession(username string, w http.ResponseWriter) {
	tokenval := crypto.GenerateRandomHash(20)
	encryptedtoken, err := bcrypt.GenerateFromPassword([]byte(tokenval), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err.Error())
	}
	app.DB.Do("HSET", "logins", username+":clienttoken", encryptedtoken)
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

	// Use bcrypt.CompareHashAndPassword() method in order to compare the encrypted token hash received from the database with the token extracted from the users cookie
	tokenerror := bcrypt.CompareHashAndPassword([]byte(tokenfromdb), []byte(tokenfromcookie))
	if tokenerror == nil {
		return true
	} else {
		return false
	}
}
