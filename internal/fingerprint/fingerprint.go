package fingerprint

import (
	"math"
	"math/cmplx"

	"github.com/danilovict2/shazam-clone/models"
	"github.com/mattetti/audio"
	"github.com/mattetti/audio/transforms"
	"github.com/mjibson/go-dsp/fft"
)

const (
	chunkSize   int   = 4000
	fuzz_factor int64 = 2
	dsFactor    int   = 4
)

var ranges = [...]int{40, 80, 120, 180, 300}

func Fingerprint(input []float64, duration float64, sampleRate uint32, songID string) ([]models.SongPoint, error) {
	downsampled, err := downsample(input, int(sampleRate), int(sampleRate)/dsFactor)
	if err != nil {
		return nil, err
	}

	chunks := partition(downsampled)
	chunkDuration := duration / float64(len(chunks))
	songPoints := make([]models.SongPoint, 0)

	for chunkIdx, chunk := range chunks {
		highscores := make([]float64, len(ranges))
		points := make([]int64, len(ranges))

		for freq := 40; freq < 300; freq++ {
			mag := math.Log(cmplx.Abs(chunk[freq]) + 1)
			index := getFreqRangeIndex(freq)

			if mag > highscores[index] {
				highscores[index] = mag
				points[index] = int64(freq)
			}
		}

		fp := hash(points[0], points[1], points[2], points[3])
		chunkTime := float64(chunkIdx) * chunkDuration
		songPoints = append(songPoints, models.SongPoint{SongID: songID, Fingerprint: fp, TimeMS: chunkTime * 1000})
	}

	return songPoints, nil
}

func downsample(input []float64, sampleRate int, targetRate int) ([]float64, error) {
	buf := audio.NewPCMFloatBuffer(input, &audio.Format{SampleRate: sampleRate})

	if err := transforms.Resample(buf, float64(targetRate)); err != nil {
		return nil, err
	}

	return buf.Floats, nil
}

func partition(audio []float64) [][]complex128 {
	sampleSize := len(audio) / chunkSize
	chunks := make([][]complex128, 0)

	for i := range sampleSize {
		chunk := make([]complex128, 0)
		for j := range chunkSize {
			chunk = append(chunk, complex(audio[(i*chunkSize)+j], 0))
		}

		chunks = append(chunks, fft.FFT(chunk))
	}

	return chunks
}

func getFreqRangeIndex(freq int) int {
	i := 0
	for ranges[i] < freq {
		i++
	}

	return i
}

func hash(p1, p2, p3, p4 int64) int64 {
	return (p4-(p4%fuzz_factor))*10000000 + (p3-(p3%fuzz_factor))*100000 + (p2-(p2%fuzz_factor))*100 + (p1 - (p1 % fuzz_factor))
}
