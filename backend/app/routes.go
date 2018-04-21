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
	"github.com/gorilla/mux"
	"github.com/yosssi/ace"
)

func (app *App) Secret(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	if app.authenticateBrowserToken(username, r) {
		template, err := ace.Load(utils.Pwd()+"templates/secret", "", nil)
		if err != nil {
			http.Error(w, "Failed to load template", http.StatusInternalServerError)
			return
		}

		if err = template.Execute(w, nil); err != nil {
			http.Error(w, "Failed to execute template", http.StatusInternalServerError)
			return
		}
	} else {
		fmt.Fprint(w, "<p>Authentication failed</p>", 403)
	}

}

func (app *App) Login(w http.ResponseWriter, r *http.Request) {
	// Retreive the file and save to disk using the FormFile method
	file, header, err := r.FormFile("file")
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(403)
		return
	}
	username := header.Filename
	defer file.Close()

	// Check is username is saved in the database
	is_member, _ := redis.Bool(app.DB.Do("SISMEMBER", "users", username))

	if !is_member {
		log.Println("Tried to login without a valid username")
		w.WriteHeader(401)
		return
	}

	out, err1 := os.Create(utils.Pwd() + "files/" + username + ".mp4")

	if err1 != nil {
		log.Println("Failed to os.Create")
		w.WriteHeader(403)
		return
	}
	_, err2 := io.Copy(out, file)
	if err2 != nil {
		os.Remove(utils.Pwd() + "files/" + username + ".mp4")
		log.Println("Failed to io.Copy enrollment #1")
		w.WriteHeader(403)
		return
	}

	// out.Close()
	// video.ConvertToH264MP4(utils.Pwd()+"files/", username)
	out, err = os.Open(utils.Pwd() + "files/" + username + ".mp4")
	if err != nil {
		os.Remove(utils.Pwd() + "files/" + username + ".mp4")
		log.Println("Failed to open converted .mp4 file")
		w.WriteHeader(403)
		return
	}

	response := structs.VideoVerificationResponse{}
	userid, _ := redis.String(app.DB.Do("HGET", "logins", username+":userid"))
	json.Unmarshal(app.VoiceIt.VideoVerification(userid, "en-US", utils.Pwd()+"files/"+username+".mp4").Bytes(), &response)
	if response.ResponseCode != "SUCC" { // Verification failed. Return user to root
		log.Println("Failed to log in")
		log.Println("mesage:", response.Message)
		log.Println("ResponseCode:", response.ResponseCode)
		os.Remove(utils.Pwd() + "files/" + username + ".mp4")
		w.WriteHeader(403)
		return
	}

	out.Close()
	os.Remove(utils.Pwd() + "files/" + username + ".mp4")
	app.setSession(username, w)
	w.WriteHeader(200)
}

func (app *App) Register(w http.ResponseWriter, r *http.Request) {
	// Grab file, save it to disk
	file, header, err := r.FormFile("file")
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(403)
		return
	}

	username := header.Filename

	defer file.Close()

	out, err1 := os.Create(utils.Pwd() + "files/" + username + ".mp4")

	if err1 != nil {
		w.WriteHeader(403)
		log.Println("Failed to os.Create")
		return
	}

	_, err2 := io.Copy(out, file)
	if err2 != nil {
		out.Close()
		os.Remove(utils.Pwd() + "files/" + username + ".mp4")
		w.WriteHeader(403)
		log.Println("Failed to io.Copy")
		return
	}
	// out.Close()
	// video.ConvertToH264MP4(utils.Pwd()+"files/", username)
	// out, err = os.Open(utils.Pwd() + "files/" + username + ".mp4")
	// if err != nil {
	// log.Println("Failed to open converted .mp4 file")
	// w.WriteHeader(403)
	// return
	// }

	// Check if user already exists in database, return user to root if user already exists in the database
	is_member, _ := redis.Bool(app.DB.Do("SISMEMBER", "users", username))
	if is_member {
		out.Close()
		os.Remove(utils.Pwd() + "files/" + username + ".mp4")
		log.Println("User tried to register existing username")
		w.WriteHeader(403)
		return
	}

	// Since doesn't already exist in the system, add user to database
	app.DB.Do("SADD", "users", username)

	// Register user in VoiceIt API
	create_user_response := structs.CreateNewUserResponse{}
	json.Unmarshal(app.VoiceIt.CreateUser().Bytes(), &create_user_response)

	if create_user_response.ResponseCode != "SUCC" {
		app.DB.Do("SREM", "users", username)
		out.Close()
		os.Remove(utils.Pwd() + "files/" + username + ".mp4")
		log.Println("Create user caused failure\n" + create_user_response.Message)
		w.WriteHeader(403)
		return
	}

	app.DB.Do("HSET", "logins", username+":userid", create_user_response.UserID)

	// Create new video enrollment for user for given group
	create_user_video_enrollment_response := structs.CreateUserVideoEnrollmentResponse{}

	json.Unmarshal(
		app.VoiceIt.CreateVideoEnrollment(
			create_user_response.UserID,
			"en-US",
			utils.Pwd()+"files/"+username+".mp4").Bytes(),
		&create_user_video_enrollment_response)

	// Process first enrollment
	if create_user_video_enrollment_response.ResponseCode != "SUCC" {
		app.DB.Do("SREM", "users", username)
		app.DB.Do("HDEL", "logins", username+":userid")
		out.Close()
		os.Remove(utils.Pwd() + "files/" + username + ".mp4")
		log.Println(create_user_video_enrollment_response.Message)
		log.Println("Creating user video enrollment #1 failed.")
		w.WriteHeader(403)
		return
	}

	os.Remove(utils.Pwd() + "files/" + username + ".mp4")

	// Process second enrollment
	file2, header, err := r.FormFile("file2")
	if err != nil {
		app.DB.Do("SREM", "users", username)
		app.DB.Do("HDEL", "logins", username+":userid")
		log.Println(err.Error())
		w.WriteHeader(403)
		return
	}
	defer file2.Close()

	out2, err3 := os.Create(utils.Pwd() + "files/" + username + "2.mp4")
	if err3 != nil {
		app.DB.Do("SREM", "users", username)
		app.DB.Do("HDEL", "logins", username+":userid")
		os.Remove(utils.Pwd() + "files/" + username + ".mp4")
		log.Println("Failed to create file " + username + "2.mp4")
		w.WriteHeader(403)
		return
	}

	_, err5 := io.Copy(out2, file2)
	if err5 != nil {
		app.DB.Do("SREM", "users", username)
		app.DB.Do("HDEL", "logins", username+":userid")
		out.Close()
		os.Remove(utils.Pwd() + "files/" + username + ".mp4")
		log.Println("Failed to io.Copy enrollment 2")
		w.WriteHeader(403)
		return
	}

	create_user_video_enrollment_response2 := structs.CreateUserVideoEnrollmentResponse{}

	json.Unmarshal(
		app.VoiceIt.CreateVideoEnrollment(
			create_user_response.UserID,
			"en-US",
			utils.Pwd()+"files/"+username+"2.mp4").Bytes(),
		&create_user_video_enrollment_response2)

	if create_user_video_enrollment_response2.ResponseCode != "SUCC" {
		app.DB.Do("SREM", "users", username)
		app.DB.Do("HDEL", "logins", username+":userid")
		out.Close()
		os.Remove(utils.Pwd() + "files/" + username + ".mp4")
		log.Println(create_user_video_enrollment_response2.Message)
		log.Println("Creating user video enrollment #2 failed.")
		w.WriteHeader(403)
		return
	}
	out2.Close()
	os.Remove(utils.Pwd() + "files/" + username + "2.mp4")

	// Process third enrollment
	file3, header, err := r.FormFile("file3")
	if err != nil {
		app.DB.Do("SREM", "users", username)
		app.DB.Do("HDEL", "logins", username+":userid")
		log.Println(err.Error())
		w.WriteHeader(403)
		return
	}
	defer file3.Close()

	// Check is username is saved in the database
	out3, err4 := os.Create(utils.Pwd() + "files/" + username + "3.mp4")
	if err4 != nil {
		app.DB.Do("SREM", "users", username)
		app.DB.Do("HDEL", "logins", username+":userid")
		os.Remove(utils.Pwd() + "files/" + username + ".mp4")
		log.Println("Failed to create file " + username + "3.mp4")
		w.WriteHeader(403)
		return
	}

	_, err6 := io.Copy(out3, file3)
	if err6 != nil {
		app.DB.Do("SREM", "users", username)
		app.DB.Do("HDEL", "logins", username+":userid")
		out.Close()
		os.Remove(utils.Pwd() + "files/" + username + ".mp4")
		w.WriteHeader(403)
		log.Println("Failed to io.Copy enrollment 3")
		return
	}

	create_user_video_enrollment_response3 := structs.CreateUserVideoEnrollmentResponse{}

	json.Unmarshal(
		app.VoiceIt.CreateVideoEnrollment(
			create_user_response.UserID,
			"en-US",
			utils.Pwd()+"files/"+username+"3.mp4").Bytes(),
		&create_user_video_enrollment_response3)

	if create_user_video_enrollment_response3.ResponseCode != "SUCC" {
		app.DB.Do("SREM", "users", username)
		app.DB.Do("HDEL", "logins", username+":userid")
		out.Close()
		os.Remove(utils.Pwd() + "files/" + username + ".mp4")
		log.Println(create_user_video_enrollment_response3.Message)
		log.Println("Creating user video enrollment #3 failed.")
		w.WriteHeader(403)
		return
	}
	out3.Close()
	os.Remove(utils.Pwd() + "files/" + username + "3.mp4")
	w.WriteHeader(302)
}
