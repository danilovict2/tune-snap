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
	fuzz_factor int64 = 2
	dsFactor    int   = 4
	frameLen    int   = 1024
	frameShift  int   = frameLen / 32
)

var ranges = [...]int{40, 80, 120, 180, 300}

func Fingerprint(input []float64, duration float64, sampleRate uint32, songID string) ([]models.SongPoint, error) {
	downsampled, err := downsample(input, int(sampleRate), int(sampleRate)/dsFactor)
	if err != nil {
		return nil, err
	}

	spectogram := stft(downsampled)
	chunkDuration := duration / float64(len(spectogram))
	songPoints := make([]models.SongPoint, 0)

	for winIdx, window := range spectogram {
		highscores := make([]float64, len(ranges))
		points := make([]int64, len(ranges))

		for freq := 40; freq < 300; freq++ {
			mag := math.Log(cmplx.Abs(window[freq]) + 1)
			index := getFreqRangeIndex(freq)

			if mag > highscores[index] {
				highscores[index] = mag
				points[index] = int64(freq)
			}
		}

		fp := hash(points[0], points[1], points[2], points[3])
		chunkTime := float64(winIdx) * chunkDuration
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

func stft(audio []float64) [][]complex128 {
	numFrames := int(float64(len(audio)-frameLen)/float64(frameShift)) + 1
	spectogram := make([][]complex128, numFrames)

	windows := make([]float64, numFrames)
	arg := 2.0 * math.Pi / float64(numFrames-1)
	for i := range windows {
		windows[i] = 0.5 - 0.5*math.Cos(arg*float64(i))
	}

	frames := make([][]float64, numFrames)
	for i := range numFrames {
		frames[i] = audio[i*frameShift : i*frameShift+frameLen]
	}

	for i, frame := range frames {
		windowed := make([]float64, len(frame))
		for _, window := range windows {
			windowed[i] = frame[i] * window
		}

		spectogram[i] = fft.FFTReal(windowed)
	}

	return spectogram
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
