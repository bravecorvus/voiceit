package app

import (
	"log"

	"github.com/garyburd/redigo/redis"
	voiceit2go "github.com/gilgameshskytrooper/VoiceIt2-Go"
	"github.com/gorilla/securecookie"
)

type App struct {
	DB                redis.Conn
	VoiceIt           voiceit2go.VoiceIt2
	CookieHandler     *securecookie.SecureCookie
	ForceSucceedLogin bool
}

func (app *App) Initialize() {
	db, err := redis.Dial("tcp", ":6379")
	// db, err := redis.DialURL(os.Getenv("REDISLOCATION"))
	app.DB = db
	if err != nil {
		log.Println(err.Error())
	}
	_, _ = app.DB.Do("SELECT", "0")

	app.VoiceIt.Initialize()

	app.CookieHandler = securecookie.New(
		securecookie.GenerateRandomKey(64),
		securecookie.GenerateRandomKey(32))

	app.ForceSucceedLogin = false
}
