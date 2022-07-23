package ytdl

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/briandowns/spinner"
	"github.com/kkdai/youtube/v2"
	"github.com/schollz/progressbar/v3"
)

func DownloadVideo(id string, path string) {
	videoID := id
	client := youtube.Client{}
	video, err := client.GetVideo(videoID)
	videoFileName := path + "/" + video.Title + ".mp4"
	if err != nil {
		log.Fatalln(err)
	}
	formats := video.Formats
	videoFile, err := os.Create(path + "/" + video.ID + ".mp4")
	if err != nil {
		log.Fatalln(err)
	}
	defer videoFile.Close()
	audioFile, err := os.Create(path + "/" + video.ID + ".mp3")
	if err != nil {
		log.Fatalln(err)
	}
	defer audioFile.Close()

	videoStream, _, err := client.GetStream(video, &formats[0])
	audioStream, _, err := client.GetStream(video, &formats.WithAudioChannels().Type("audio/mp4")[0])

	videoBar := progressbar.DefaultBytes(formats[0].ContentLength, fmt.Sprintf("Downloading video %s", video.Title))
	io.Copy(io.MultiWriter(videoFile, videoBar), videoStream)

	audioBar := progressbar.DefaultBytes(formats.WithAudioChannels().Type("audio/mp4")[0].ContentLength, fmt.Sprintf("Downloading audio for video %s", video.Title))
	io.Copy(io.MultiWriter(audioFile, audioBar), audioStream)

	if err != nil {
		log.Fatalln(err)
	}

	_, err = io.Copy(videoFile, videoStream)
	if err != nil {
		log.Fatalln(err)
	}
	_, err = io.Copy(audioFile, audioStream)
	if err != nil {
		log.Fatalln(err)
	}
	s := spinner.New(spinner.CharSets[6], 100*time.Millisecond)
	s.Prefix = "Merging video and audio "
	s.Start()
	cmd := exec.Command("ffmpeg", "-i", videoFile.Name(), "-i", audioFile.Name(), "-c:v", "copy", "-map", "0:v", "-map", "1:a", "-y", videoFileName)
	err = cmd.Run()
	s.Stop()
	if err != nil {
		log.Fatalf("Error merging video and audio with ffmpeg: %v", err)
	}
	os.Remove(videoFile.Name())
	os.Remove(audioFile.Name())
}
