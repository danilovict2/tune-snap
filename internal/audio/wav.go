package audio

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	channels   = "2"
	sampleRate = "48000"
)

func convertToWav(inputPath string) error {
	ext := filepath.Ext(inputPath)
	outputPath := strings.TrimRight(inputPath, ext) + ".wav"

	tmpPath := filepath.Join(os.Getenv("SONGS_DIR"), "tmp_"+filepath.Base(outputPath))
	defer os.Remove(tmpPath)

	comm := exec.Command("ffmpeg", "-y", "-i", inputPath, "-c", "pcm_s16le", "-ar", sampleRate, "-ac", channels, tmpPath)

	if err := comm.Run(); err != nil {
		return err
	}

	tmpFile, err := os.Open(tmpPath)
	if err != nil {
		return err
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	if _, err := io.Copy(outputFile, tmpFile); err != nil {
		return err
	}

	return nil
}
