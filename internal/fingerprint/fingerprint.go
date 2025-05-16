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
	fuzzFactor int64 = 2
	dsFactor   int   = 4
	frameLen   int   = 1024
	frameShift int   = 512
)

var bins = []struct{ start, end int }{{40, 80}, {80, 120}, {120, 300}, {300, 600}}

func Fingerprint(input []float64, duration float64, sampleRate uint32, songID string) ([]models.SongPoint, error) {
	downsampled, err := downsample(input, int(sampleRate), int(sampleRate)/dsFactor)
	if err != nil {
		return nil, err
	}

	spectogram := stft(downsampled)
	chunkDuration := duration / float64(len(spectogram))
	songPoints := make([]models.SongPoint, 0)

	for chunkIdx, chunk := range spectogram {
		highscores := make([]float64, len(bins))
		peaks := make([]int64, len(bins))

		for i, bin := range bins {
			for freq := bin.start; freq < bin.end; freq++ {
				mag := math.Log(cmplx.Abs(chunk[freq]) + 1)

				if mag > highscores[i] {
					highscores[i] = mag
					peaks[i] = int64(freq)
				}
			}
		}

		fp := hash(peaks[0], peaks[1], peaks[2], peaks[3])
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

func stft(audio []float64) [][]complex128 {
	numFrames := int(float64(len(audio)-frameLen)/float64(frameShift)) + 1
	spectogram := make([][]complex128, numFrames)

	windows := make([]float64, frameLen)
	arg := 2.0 * math.Pi / float64(frameLen-1)
	for i := range windows {
		windows[i] = 0.5 - 0.5*math.Cos(arg*float64(i))
	}

	frames := make([][]float64, numFrames)
	for i := range numFrames {
		frames[i] = audio[i*frameShift : i*frameShift+frameLen]
	}

	for i, frame := range frames {
		windowed := make([]float64, frameLen)
		for j, window := range windows {
			windowed[j] = frame[j] * window
		}

		spectogram[i] = fft.FFTReal(windowed)
	}

	return spectogram
}

func hash(p1, p2, p3, p4 int64) int64 {
	return (p4-(p4%fuzzFactor))*10000000 + (p3-(p3%fuzzFactor))*100000 + (p2-(p2%fuzzFactor))*100 + (p1 - (p1 % fuzzFactor))
}
