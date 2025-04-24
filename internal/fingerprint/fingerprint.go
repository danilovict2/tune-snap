package fingerprint

import (
	"fmt"
	"math"
	"math/cmplx"

	"github.com/scientificgo/fft"
)

const (
	chunkSize   int   = 4000
	fuzz_factor int64 = 2
)

var ranges = [...]int{40, 80, 120, 180, 300}

func Fingerprint(audio []byte) {
	chunks := partition(audio)
	for _, chunk := range chunks {
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

		h := hash(points[0], points[1], points[2], points[3])
		fmt.Println(h)
	}

}

func partition(audio []byte) [][]complex128 {
	chunks := make([][]complex128, 0)
	sampleSize := len(audio) / chunkSize

	for i := range sampleSize {
		chunk := make([]complex128, 0)
		for j := range chunkSize {
			chunk = append(chunk, complex(float64(audio[(i*chunkSize)+j]), 0))
		}

		chunks = append(chunks, fft.Fft(chunk, false))
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
