package hamming

import (
	"github.com/Comdex/imgo"
	"math"
	"strings"
)

type HammingConfig struct {
	Size int
	C    []float64
}

func New(size int) *HammingConfig {
	h := &HammingConfig{}
	h.Size = size
	h.C = make([]float64, size)
	for i := 1; i < size; i++ {
		h.C[i] = 1
	}
	h.C[0] = 1.0 / math.Sqrt(2.0)
	return h

}

func (h *HammingConfig) HammingDistance(a, b [][][]uint8) int {

	hasha := h.getHash(a)
	hashb := h.getHash(b)
	distance := h.distance(hasha, hashb)
	return distance
}

func (h *HammingConfig) applyDCT(f [][][]uint8) [][][]uint8 {
	N := h.Size
	F := imgo.NewRGBAMatrix(len(f), len(f[0]))
	for u := 0; u < N; u++ {
		for v := 0; v < N; v++ {
			sum := 0.0
			for i := 0; i < N; i++ {
				for j := 0; j < N; j++ {
					sum += math.Cos(((2*float64(i)+1)/(2.0*float64(N)))*float64(u)*math.Pi) * math.Cos(((2.0*float64(j)+1.0)/(2.0*float64(N)))*float64(v)*math.Pi) * float64(f[i][j][0])
				}
			}
			sum *= (h.C[u] * h.C[v]) / 4.0
			F[u][v][0] = uint8(sum)
			F[u][v][1] = uint8(sum)
			F[u][v][2] = uint8(sum)
			F[u][v][3] = uint8(255)
		}
	}
	return F
}

func (h *HammingConfig) distance(s1, s2 string) int {
	counter := 0
	for k := 0; k < len(s1); k++ {
		if !strings.EqualFold(s1[k:k+1], s2[k:k+1]) {
			counter++
		}
	}
	return counter
}

// Returns a 'binary string' (like. 001010111011100010) which is easy to do a hamming distance on.
func (h *HammingConfig) getHash(imgBytes [][][]uint8) string {
	height := len(imgBytes)
	width := len(imgBytes[0])
	/* 1. Reduce size.
	 * Like Average Hash, pHash starts with a small image.
	 * However, the image is larger than 8x8; 32x32 is a good size.
	 * This is really done to simplify the DCT computation and not
	 * because it is needed to reduce the high frequencies.
	 */
	//img = resize(img, size, size);

	/* 2. Reduce color.
	 * The image is reduced to a grayscale just to further simplify
	 * the number of computations.
	 */
	//img = grayscale(img);

	/* 3. Compute the DCT.
	 * The DCT separates the image into a collection of frequencies
	 * and scalars. While JPEG uses an 8x8 DCT, this algorithm uses
	 * a 32x32 DCT.
	 */
	//start := time.Now()
	dctVals := h.applyDCT(imgBytes)
	//fmt.Printf("DCT use: %d ms\n", time.Now().Sub(start).Milliseconds())

	/* 4. Reduce the DCT.
	 * This is the magic step. While the DCT is 32x32, just keep the
	 * top-left 8x8. Those represent the lowest frequencies in the
	 * picture.
	 */
	/* 5. Compute the average value.
	 * Like the Average Hash, compute the mean DCT value (using only
	 * the 8x8 DCT low-frequency values and excluding the first term
	 * since the DC coefficient can be significantly different from
	 * the other values and will throw off the average).
	 */
	total := 0
	for x := 0; x < height; x++ {
		for y := 0; y < width; y++ {
			total += int(dctVals[x][y][0])
		}
	}
	total -= int(dctVals[0][0][0])
	avg := float64(total) / float64((height*width)-1)

	/* 6. Further reduce the DCT.
	 * This is the magic step. Set the 64 hash bits to 0 or 1
	 * depending on whether each of the 64 DCT values is above or
	 * below the average value. The result doesn't tell us the
	 * actual low frequencies; it just tells us the very-rough
	 * relative scale of the frequencies to the mean. The result
	 * will not vary as long as the overall structure of the image
	 * remains the same; this can survive gamma and color histogram
	 * adjustments without a problem.
	 */
	hash := ""
	for x := 0; x < height; x++ {
		for y := 0; y < width; y++ {
			if x != 0 && y != 0 {
				if float64(dctVals[x][y][0]) > avg {
					hash += "1"
				} else {
					hash += "0"
				}
			}
		}
	}

	return hash
}
