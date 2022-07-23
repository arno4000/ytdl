package ytdl

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"

	"github.com/kkdai/youtube/v2"
	"github.com/schollz/progressbar/v3"
)

func DownloadAudio(id string, path string) {

	videoID := id
	client := youtube.Client{}
	video, err := client.GetVideo(videoID)
	fileNameRegex := regexp.MustCompile(`[<>:"\/\|?*^]`)
	videoName := fileNameRegex.ReplaceAllString(video.Title, "")
	audioPath := path + "/" + videoName + ".mp3"
	if err != nil {
		log.Fatalln(err)
	}
	formats := video.Formats
	audioFile, err := os.Create(audioPath)
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
