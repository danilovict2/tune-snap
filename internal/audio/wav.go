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
)

const (
	channels   = "2"
	sampleRate = "44100"
)

type Wav struct {
	Duration   float64
	Audio      []float64
	SampleRate uint32
}

type wavHeader struct {
	ChunkID       [4]byte
	ChunkSize     uint32
	Format        [4]byte
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

	tmpPath := filepath.Join(".", "tmp_"+filepath.Base(outputPath))
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

func ReadWav(wavFile io.Reader) (*Wav, error) {
	data, err := io.ReadAll(wavFile)
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

	if string(header.ChunkID[:]) != "RIFF" || string(header.Format[:]) != "WAVE" || header.AudioFormat != 1 || header.BitsPerSample != 16 {
		return nil, fmt.Errorf("unsupported WAV file format")
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
	min, max := -32768.0, 32768.0

	for i := 0; i < len(data); i += 2 {
		sample := binary.LittleEndian.Uint16(data[i : i+2])
		samples[i/2] = float64(sample)
	}

	// Normalize sample to [-1, 1]
	for i, sample := range samples {
		samples[i] = 2*((sample-min)/(max-min)) - 1
	}

	return samples, nil
}
