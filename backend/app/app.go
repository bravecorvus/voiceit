package app

import (
	"log"
	"os"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/gilgameshskytrooper/voiceit/backend/email"
	"github.com/gilgameshskytrooper/voiceit/backend/utils"
	"github.com/gorilla/securecookie"
	"golang.org/x/crypto/bcrypt"
)

type App struct {
	DB             redis.Conn
	URL            string
	CookieHandler  *securecookie.SecureCookie
	PasswordResets []PasswordResetStruct
	PendingUsers   []AddNewUserStruct
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

	adminemail := os.Getenv("VOICEITADMINEMAIL")
	adminpass := os.Getenv("VOICEITADMINPASSWORD")

	if adminpass == "" {
		newpassword := utils.GenerateRandomHash(40)
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newpassword), bcrypt.DefaultCost)
		if err != nil {
			log.Println(err.Error())
		}
		_, _ = app.DB.Do("HSET", "logins", "admin:password", hashedPassword)
		sentornotsent, status := email.SendEmail(adminemail, "VoiceIt admin password change", "VoiceIt superadmin password was not defined as an environment variable and hence, it has been randomly generated.\n\nThe password is ["+newpassword+"].")
		if !sentornotsent {
			log.Println("Sending email failed at", status)
		}
	} else {
		storedhashedpass, _ := redis.String(app.DB.Do("HGET", "logins", "admin:password"))
		samepassword := bcrypt.CompareHashAndPassword([]byte(storedhashedpass), []byte(adminpass))

		if samepassword != nil {
			newpasswordhash, err := bcrypt.GenerateFromPassword([]byte(adminpass), bcrypt.DefaultCost)
			if err != nil {
				log.Println("Couldn't create password hash from given password")
			}
			_, _ = app.DB.Do("HSET", "logins", "admin:password", newpasswordhash)
			sentornotsent, status := email.SendEmail(adminemail, "VoiceIt superadmin password change", "Stored VoiceIt admin password differed from the one provided by the ADMINPASSWORD environment variable and hence, it has been changed to be the one parsed from environment variables.\n\nThe password is ["+adminpass+"].")
			if !sentornotsent {
				log.Println("Sending email failed at", status)
			}
		}
	}

	app.CookieHandler = securecookie.New(
		securecookie.GenerateRandomKey(64),
		securecookie.GenerateRandomKey(32))
	app.URL = os.Getenv("VOICEITURL")
	if app.URL == "" {
		app.URL = "http://localhost:8080"
	}
}
