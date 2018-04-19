package app

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/garyburd/redigo/redis"
	"github.com/gilgameshskytrooper/voiceit/backend/structs"
	"github.com/gilgameshskytrooper/voiceit/backend/utils"
	"github.com/yosssi/ace"
)

func (app *App) Login(w http.ResponseWriter, r *http.Request) {
	// Retreive the file and save to disk using the FormFile method
	file, header, err := r.FormFile("file")
	if err != nil {
		log.Println(err.Error())
		fmt.Fprint(w, "Failed. Please Try Again")
		http.Redirect(w, r, "/", 302)
		return
	}
	username := header.Filename
	defer file.Close()

	out, err1 := os.Create(utils.Pwd() + "files/" + username + ".mp4")

	if err1 != nil {
		fmt.Fprint(w, "Failed. Please Try Again")
		log.Println("Failed to os.Create")
		http.Redirect(w, r, "/", 302)
		return
	}
	_, err2 := io.Copy(out, file)
	if err2 != nil {
		fmt.Fprint(w, "Failed. Please Try Again")
		log.Println("Failed to io.Copy")
		http.Redirect(w, r, "/", 302)
		return
	}

	response := structs.VideoVerificationResponse{}
	userid, _ := redis.String(app.DB.Do("HGET", "users", username+":userid"))
	json.Unmarshal(app.VoiceIt.VideoIdentification(userid, "en-US", utils.Pwd()+"files/"+username+".mp4").Bytes(), &response)
	if response.ResponseCode != "SUCC" { // Verification failed. Return user to root
		fmt.Fprint(w, "Failed to log in")
		log.Println("Failed to log in")
		http.Redirect(w, r, "/", 302)
		return
	}

	out.Close()
	os.Remove(utils.Pwd() + "files/" + username + ".mp4")

	template, err := ace.Load(utils.Pwd()+"templates/secret", "", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = template.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (app *App) Register(w http.ResponseWriter, r *http.Request) {
	// Grab file, save it to disk
	file, header, err := r.FormFile("file")
	if err != nil {
		log.Println(err.Error())
		fmt.Fprint(w, "Failed. Please Try Again")
		http.Redirect(w, r, "/", 302)
		return
	}
	username := header.Filename

	defer file.Close()

	out, err1 := os.Create(utils.Pwd() + "files/" + username + ".mp4")

	if err1 != nil {
		fmt.Fprint(w, "Failed. Please Try Again")
		log.Println("Failed to os.Create")
		http.Redirect(w, r, "/", 302)
		return
	}

	_, err2 := io.Copy(out, file)
	if err2 != nil {
		out.Close()
		os.Remove(utils.Pwd() + "files/" + username + ".mp4")
		fmt.Fprint(w, "Failed. Please Try Again")
		log.Println("Failed to io.Copy")
		http.Redirect(w, r, "/", 302)
		return
	}

	// Check if user already exists in database, return user to root if user already exists in the database
	is_member, _ := redis.Bool(app.DB.Do("SISMEMBER", "users", username))
	if is_member {
		out.Close()
		os.Remove(utils.Pwd() + "files/" + username + ".mp4")
		log.Println("User tried to register existing username")
		fmt.Fprint(w, "User already exists. Please choose another username")
		http.Redirect(w, r, "/", 302)
		return
	}

	// Register user in VoiceIt API
	create_user_response := structs.CreateNewUserResponse{}
	json.Unmarshal(app.VoiceIt.CreateUser().Bytes(), &create_user_response)

	if create_user_response.ResponseCode != "SUCC" {
		out.Close()
		os.Remove(utils.Pwd() + "files/" + username + ".mp4")
		log.Println("Create user caused failure\n" + create_user_response.Message)
		fmt.Fprint(w, "Error communicating with VoiceIt API to create new user.")
		http.Redirect(w, r, "/", 302)
		return
	}

	app.DB.Do("HSET", "users", username+":userid", create_user_response.UserID)
	log.Println("New User ID", create_user_response.UserID)

	// Create new video enrollment for user for given group
	create_user_video_enrollment_response := structs.CreateUserVideoEnrollmentResponse{}

	json.Unmarshal(
		app.VoiceIt.CreateVideoEnrollment(
			create_user_response.UserID,
			"en-US",
			utils.Pwd()+"files/"+username+".mp4").Bytes(),
		&create_user_video_enrollment_response)

	if create_user_video_enrollment_response.ResponseCode != "SUCC" {
		out.Close()
		os.Remove(utils.Pwd() + "files/" + username + ".mp4")
		log.Println("Creating user video enrollment failed.")
		fmt.Fprint(w, "Creating user video enrollment failed.")
		http.Redirect(w, r, "/", 302)
		return
	}

	out.Close()
	os.Remove(utils.Pwd() + "files/" + username + ".mp4")
	fmt.Fprint(w, "Successfully created account. Please try logging in.")
	http.Redirect(w, r, "/", 302)
}
