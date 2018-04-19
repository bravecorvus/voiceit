package app

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/garyburd/redigo/redis"
	voiceit2go "github.com/gilgameshskytrooper/VoiceIt2-Go"
	"github.com/gilgameshskytrooper/voiceit/backend/structs"
	"github.com/gorilla/securecookie"
)

type App struct {
	DB                redis.Conn
	URL               string
	CookieHandler     *securecookie.SecureCookie
	PasswordResets    []PasswordResetStruct
	PendingUsers      []AddNewUserStruct
	VoiceIt           voiceit2go.VoiceIt2
	VoiceItUserGroup  string
	VoiceItAdminGroup string
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
	// db, err := redis.Dial("tcp", ":6379")
	db, err := redis.DialURL(os.Getenv("REDISLOCATION"))
	app.DB = db
	if err != nil {
		log.Println(err.Error())
	}
	_, _ = app.DB.Do("SELECT", "0")

	app.VoiceIt.Initialize()

	usergroupid, _ := redis.String(app.DB.Do("HGET", "system", "usergroupid"))
	if usergroupid == "" { // usergroupid has not been defined in the Redis database yet, so make call to create the group, and we will set the groupid of the response as the user group id
		response := structs.CreateGroupResponse{}
		json.Unmarshal(app.VoiceIt.CreateGroup("users").Bytes(), &response)
		groupid := response.GroupID
		app.DB.Do("HSET", "system", "usergroupid", groupid) // Save this groupid into the system usergroupid as an Redis Hash Set
	}
}
