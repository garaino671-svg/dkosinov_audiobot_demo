package main

import (
	"encoding/json"
	"errors"
	"os/exec"
)

const maxDurationSeconds = 2 * 60 * 60 // 2 часа

type ytMeta struct {
	Duration *int `json:"duration"`
}

func checkVideoDuration(url string) error {
	cmd := exec.Command("yt-dlp", "--dump-json", url)

	output, err := cmd.Output()
	if err != nil {
		return err
	}

	var meta ytMeta
	if err := json.Unmarshal(output, &meta); err != nil {
		return err
	}

	// live-стримы: duration == nil
	if meta.Duration == nil {
		return errors.New("live stream or unknown duration")
	}

	if *meta.Duration > maxDurationSeconds {
		return errors.New("video is longer than 2 hours")
	}

	return nil
}
