package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/garyburd/redigo/redis"
	"github.com/gilgameshskytrooper/voiceit/backend/structs"
	"github.com/gilgameshskytrooper/voiceit/backend/utils"
)

func (app *App) Login(w http.ResponseWriter, r *http.Request) {
	// Retreive the file and save to disk using the FormFile method
	file, _, err := r.FormFile("file")
	if err != nil {
		fmt.Println(err.Error())
		fmt.Fprint(w, "Failed. Please Try Again")
		http.Redirect(w, r, "/", 302)
		return
	}
	defer file.Close()

	userhash := utils.GenerateRandomHash(40)
	out, err1 := os.Create(utils.Pwd() + "files/" + userhash + ".mp4")

	if err1 != nil {
		fmt.Fprint(w, "Failed. Please Try Again")
		fmt.Println("Failed to os.Create")
		http.Redirect(w, r, "/", 302)
		return
	}
	defer out.Close()
	_, err2 := io.Copy(out, file)
	if err2 != nil {
		fmt.Fprint(w, "Failed. Please Try Again")
		fmt.Println("Failed to io.Copy")
		http.Redirect(w, r, "/", 302)
		return
	}

	usergroupid, _ := redis.String(app.DB.Do("HGET", "system", "usergroupid"))
	response := structs.CreateNewUserResponse{}
	json.Unmarshal(app.VoiceIt.VideoIdentification(usergroupid, "en-US", utils.Pwd()+"files/"+userhash+".mp4", true).Bytes(), &response)

}

func (app *App) Register(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	fmt.Println("filename", header.Filename)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Fprint(w, "Failed. Please Try Again")
		http.Redirect(w, r, "/", 302)
		return
	}
	defer file.Close()

	out, err1 := os.Create(utils.Pwd() + "files/" + header.Filename + ".mp4")

	if err1 != nil {
		fmt.Fprint(w, "Failed. Please Try Again")
		fmt.Println("Failed to os.Create")
		http.Redirect(w, r, "/", 302)
		return
	}
	defer out.Close()
	_, err2 := io.Copy(out, file)
	if err2 != nil {
		fmt.Fprint(w, "Failed. Please Try Again")
		fmt.Println("Failed to io.Copy")
		http.Redirect(w, r, "/", 302)
		return
	}

}
