package main

// import (
// 	"bufio"
// 	"context"
// 	"os/exec"
// 	"path/filepath"
// 	"strings"
// )

// func downloadMP3WithProgress(
// 	ctx context.Context,
// 	url string,
// 	onProgress func(string),
// ) (string, error) {

// 	cmd := exec.CommandContext(
// 		ctx,
// 		"yt-dlp",
// 		"-x",
// 		"--audio-format", "mp3",
// 		"--audio-quality", "0",
// 		"--progress-template", "%(progress._percent_str)s",
// 		"-o", "%(title)s.%(ext)s",
// 		url,
// 	)

// 	stdout, err := cmd.StdoutPipe()
// 	if err != nil {
// 		return "", err
// 	}

// 	cmd.Stderr = cmd.Stdout

// 	if err := cmd.Start(); err != nil {
// 		return "", err
// 	}

// 	scanner := bufio.NewScanner(stdout)
// 	var lastPercent string

// 	for scanner.Scan() {
// 		line := strings.TrimSpace(scanner.Text())

// 		if strings.HasSuffix(line, "%") && line != lastPercent {
// 			lastPercent = line
// 			onProgress(line)
// 		}
// 	}

// 	if err := cmd.Wait(); err != nil {
// 		return "", err
// 	}

// 	return findDownloadedMP3(), nil
// }

// func findDownloadedMP3() string {
// 	files, err := filepath.Glob("*.mp3")
// 	if err != nil {
// 		return ""
// 	}

// 	if len(files) == 0 {
// 		return ""
// 	}

// 	// берём последний (yt-dlp обычно создаёт один файл)
// 	return files[len(files)-1]
// }
