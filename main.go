package main

import (
	"log"
	"os/user"
	"path"

	"github.com/arno4000/ytdl/pkg/ytdl"
	"github.com/manifoldco/promptui"
)

func main() {
	videoPromt := promptui.Prompt{
		Label:    "Youtube-URL",
		Validate: ytdl.ValidateID,
	}
	videoID, err := videoPromt.Run()
	if err != nil {
		log.Fatalln(err)
	}
	actionPromt := promptui.Select{
		Label: "Audio or Video",
		Items: []string{"Audio", "Video"},
	}
	_, result, err := actionPromt.Run()
	if err != nil {
		log.Fatalln(err)
	}
	usr, err := user.Current()
	if err != nil {
		log.Fatalln(err)
	}
	videoPath := path.Join(usr.HomeDir, "Videos")
	audioPath := path.Join(usr.HomeDir, "Music")
	if result == "Audio" {
		storagePromt := promptui.Prompt{
			Label:   "Download path",
			Default: audioPath,
		}
		audioPath, err := storagePromt.Run()
		if err != nil {
			log.Fatalln(err)
		}
		ytdl.DownloadAudio(videoID, audioPath)

	} else if result == "Video" {
		storagePromt := promptui.Prompt{
			Label:   "Download path",
			Default: videoPath,
		}
		videoPath, err := storagePromt.Run()
		if err != nil {
			log.Fatalln(err)
		}
		ytdl.DownloadVideo(videoID, videoPath)

	}
}
