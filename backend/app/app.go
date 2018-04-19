package app

import (
	"log"
	"time"

	"github.com/garyburd/redigo/redis"
	voiceit2go "github.com/gilgameshskytrooper/VoiceIt2-Go"
)

type App struct {
	DB          redis.Conn
	URL         string
	VoiceIt     voiceit2go.VoiceIt2
	UserGroupID string
}

type PasswordResetStruct struct {
	Link           string
	AssociatedUser string
	Expiration     time.Time
}

type AddNewUserStruct struct {
	Link           string
	AssociatedUser string
	Expiration     time.Time
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
}
