package ytdl

import (
	"errors"

	"github.com/kkdai/youtube/v2"
)

func ValidateID(id string) error {
	_, err := youtube.ExtractVideoID(id)
	if err != nil {
		return errors.New("Invalid Youtube-URL")
	}
	return nil
}
