package ytdl

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"time"

	"github.com/briandowns/spinner"
	"github.com/kkdai/youtube/v2"
	"github.com/schollz/progressbar/v3"
)

func DownloadVideo(id string, path string) {
	videoID := id
	client := youtube.Client{}
	video, err := client.GetVideo(videoID)
	fileNameRegex := regexp.MustCompile(`[<>:"\/\|?*^]`)
	videoName := fileNameRegex.ReplaceAllString(video.Title, "")
	videoPathID := path + "/" + video.ID + ".mp4"
	videoPath := path + "/" + videoName + ".mp4"
	audioPathID := path + "/" + video.ID + ".mp3"
	if err != nil {
		log.Fatalln(err)
	}
	formats := video.Formats
	videoFile, err := os.Create(videoPathID)
	if err != nil {
		log.Fatalln(err)
	}
	audioFile, err := os.Create(audioPathID)
	if err != nil {
		log.Fatalln(err)
	}

	videoStream, _, err := client.GetStream(video, &formats[0])
	audioStream, _, err := client.GetStream(video, &formats.WithAudioChannels().Type("audio/mp4")[0])

	videoBar := progressbar.DefaultBytes(formats[0].ContentLength, fmt.Sprintf("Downloading video %s", video.Title))
	io.Copy(io.MultiWriter(videoFile, videoBar), videoStream)

	audioBar := progressbar.DefaultBytes(formats.WithAudioChannels().Type("audio/mp4")[0].ContentLength, fmt.Sprintf("Downloading audio for video %s", videoPath))
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
	cmd := exec.Command("ffmpeg", "-i", videoFile.Name(), "-i", audioFile.Name(), "-c:v", "copy", "-map", "0:v", "-map", "1:a", "-y", videoPath)
	err = cmd.Run()
	s.Stop()
	if err != nil {
		log.Fatalf("Error merging video and audio with ffmpeg: %v", err)
	}
	time.Sleep(5)
	videoFile.Close()
	audioFile.Close()
	err = os.Remove(videoPathID)
	if err != nil {
		log.Fatalln(err)
	}
	err = os.Remove(audioPathID)
	if err != nil {
		log.Fatalln(err)
	}
}
