package ytdl

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/kkdai/youtube/v2"
	"github.com/schollz/progressbar/v3"
)

func DownloadAudio(id string, path string) {
	exec.Command("rm", "-rf", "./*.mp3").Run()
	exec.Command("rm", "-rf", "./*.mp4").Run()
	videoID := id
	client := youtube.Client{}
	video, err := client.GetVideo(videoID)
	if err != nil {
		log.Fatalln(err)
	}
	formats := video.Formats
	audioFile, err := os.Create(path + "/" + video.Title + ".mp3")
	if err != nil {
		log.Fatalln(err)
	}
	defer audioFile.Close()
	audioStream, _, err := client.GetStream(video, &formats.WithAudioChannels().Type("audio/mp4")[0])
	if err != nil {
		log.Fatalln(err)
	}
	audioBar := progressbar.DefaultBytes(formats.WithAudioChannels().Type("audio/mp4")[0].ContentLength, fmt.Sprintf("Downloading audio for video %s", video.Title))
	io.Copy(io.MultiWriter(audioFile, audioBar), audioStream)
	_, err = io.Copy(audioFile, audioStream)
	if err != nil {
		log.Fatalln(err)
	}
}
