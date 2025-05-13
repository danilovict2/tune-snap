package audio

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/kashifkhan0771/utils/math"
)

const (
	channels   = "2"
	sampleRate = "48000"
)

type Wav struct {
	Duration   float64
	Audio      []float64
	SampleRate uint32
}

type wavHeader struct {
	ChunkID       uint32
	ChunkSize     uint32
	Format        uint32
	Subchunk1ID   uint32
	Subchunk1Size uint32
	AudioFormat   uint16
	NumChannels   uint16
	SampleRate    uint32
	ByteRate      uint32
	BlockAlign    uint16
	BitsPerSample uint16
	Subchunk2ID   uint32
	Subchunk2Size uint32
}

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

func readWav(path string) (*Wav, error) {
	if filepath.Ext(path) != ".wav" {
		return nil, fmt.Errorf("invalid file format: expected a .wav file, got %s", filepath.Ext(path))
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if len(data) < 44 {
		return nil, fmt.Errorf("invalid wav file: file size is too small to contain a valid header")
	}

	var header wavHeader
	if err := binary.Read(bytes.NewReader(data[:44]), binary.LittleEndian, &header); err != nil {
		return nil, err
	}

	sample, err := BytesToSamples(data[44:])
	if err != nil {
		return nil, err
	}

	info := Wav{
		Duration:   float64(len(data[44:])) / float64(header.ByteRate),
		Audio:      sample,
		SampleRate: header.SampleRate,
	}

	return &info, nil
}

func BytesToSamples(data []byte) ([]float64, error) {
	if len(data)%2 != 0 {
		return nil, fmt.Errorf("invalid data length: expected even number of bytes, got %d", len(data))
	}

	sampleLength := len(data) / 2
	samples := make([]float64, sampleLength)
	min, max := 1e9, -1e9

	for i := 0; i < len(data); i += 2 {
		sample := binary.LittleEndian.Uint16(data[i : i+2])
		samples[i/2] = float64(sample)

		min = math.Min(min, samples[i/2])
		max = math.Max(max, samples[i/2])
	}

	// Normalize sample to [-1, 1]
	for i, sample := range samples {
		samples[i] = 2*((sample-min)/(max-min)) - 1
	}

	return samples, nil
}
