package video

import (
	"log"
	"os"
	"os/exec"
)

func ConvertToH264MP4(path, username string) {
	var convert = exec.Command("ffmpeg", "-i", path+username+".mp4", "-vcodec", "libx264", "-preset", "veryslow", "-acodec", "aac", path+username+"2.mp4")
	// var convert = exec.Command("ffmpeg", "-framerate", "24", "-i", path+username+".mp4", "-preset", "slow", "-acodec", "aac", path+username+"2.mp4")
	if err := convert.Run(); err != nil {
		log.Println("Failed to start ffmpeg command")
	}
	if err2 := os.Remove(path + username + ".mp4"); err2 != nil {
		log.Println("Failed to delete old file")
	}
	if err3 := os.Rename(path+username+"2.mp4", path+username+".mp4"); err3 != nil {
		log.Println("Failed to rename file " + username + "2.mp4 to " + username + ".mp4")
	}
}
