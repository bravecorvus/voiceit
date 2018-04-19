package app

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/garyburd/redigo/redis"
	"golang.org/x/crypto/bcrypt"
)

func getIPFromRequest(req *http.Request) (string, error) {
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return "", fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		return "", fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
	}
	return userIP.String(), nil
}

func (app *App) authenticateLogin(username, password string, r *http.Request) bool {
	ip, err := getIPFromRequest(r)
	if err != nil {
		log.Println(err.Error())
	}
	isbanned, _ := redis.Bool(app.DB.Do("SISMEMBER", "bannedips", ip))
	if isbanned {
		return false
	}
	savedhashpassword, _ := redis.String(app.DB.Do("HGET", "logins", username+":password"))
	valid := bcrypt.CompareHashAndPassword([]byte(savedhashpassword), []byte(password))
	if valid != nil {
		log.Println("Is Not Valid Password")
		return false
	}
	return true
}
