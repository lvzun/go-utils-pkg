package sobel

import "C"
import (
	"image"
	"image/color"
	"math"
	"reflect"
	"unsafe"
)

var (
	sobelX = [3][3]int{
		{-1, 0, 1},
		{-2, 0, 2},
		{-1, 0, 1},
	}

	sobelY = [3][3]int{
		{-1, -2, -1},
		{0, 0, 0},
		{1, 2, 1},
	}

	sobelXF = [9]float64{
		-1, 0, 1,
		-2, 0, 2,
		-1, 0, 1,
	}

	sobelYF = [9]float64{
		-1, -2, -1,
		0, 0, 0,
		1, 2, 1,
	}

	sobelXL = [9]int{
		-1, 0, 1,
		-2, 0, 2,
		-1, 0, 1,
	}

	sobelYL = [9]int{
		-1, -2, -1,
		0, 0, 0,
		1, 2, 1,
	}

	sharaX = [3][3]int{
		{-3, 0, 3},
		{-10, 0, 10},
		{-3, 0, 3},
	}
	sharaY = [3][3]int{
		{3, 10, 3},
		{0, 0, 0},
		{-3, -10, -3},
	}

	laplasianX = [3][3]int{
		{1, 1, 1},
		{1, -8, 1},
		{1, 1, 1},
	}
	laplasianY = [3][3]int{
		{1, 1, 1},
		{1, -8, 1},
		{1, 1, 1},
	}

	sharpenX = [3][3]int{
		{0, -1, 0},
		{-1, 5, -1},
		{0, -1, 0},
	}
	sharpenY = [3][3]int{
		{0, -1, 0},
		{-1, 5, -1},
		{0, -1, 0},
	}
)

type FilterType int

const kernelSize = 3

const (
	Sobel FilterType = iota
	SobelFast
	Laplasian
	Shara
	Sharpen
)

type filterFunc func(img *image.Gray, x int, y int) (uint32, uint32)

func Filter(img image.Image, flt FilterType) *image.Gray {
	grayImg := ToGrayscale(img)
	return FilterGrayFast(grayImg, flt)
}

func FilterMath(img image.Image, flt FilterType) *image.Gray {
	grayImg := ToGrayscale(img)
	return FilterGrayMath(grayImg)
}

func FilterSimd(img image.Image, flt FilterType) *image.Gray {
	grayImg := ToGrayscale(img)
	return FilterGraySimd(grayImg)
}

func getFilterFunc(flt FilterType) (res filterFunc) {
	switch flt {
	case Sobel:
		res = applySobelFilter
		break
	case SobelFast:
		res = applySobelFilterFast
		break
	case Laplasian:
		res = applyLaplasianFilter
		break
	case Shara:
		res = applySharaFilter
		break
	case Sharpen:
		res = applySharpenFilter
	}
	return res
}

//for better optimization in case of input gray image
func FilterGray(grayImg *image.Gray, flt FilterType) (filtered *image.Gray) {
	max := grayImg.Bounds().Max
	min := grayImg.Bounds().Min
	/* filtered image must be two pixels shorter, because
	there must be a row of pixels on each side of a pixel for the sobel operator
	to work*/
	filtered = image.NewGray(image.Rect(max.X-2, max.Y-2, min.X, min.Y))
	width := max.X - 1 //to provide 1 pixel "border"
	height := max.Y - 1
	applay := getFilterFunc(flt)
	for x := 1; x < width; x++ {
		for y := 1; y < height; y++ {
			fX, fY := applay(grayImg, x, y)
			v := ISqrt((fX*fX)+(fY*fY)) + 1 // +1 to make it ceil
			pixel := color.Gray{Y: uint8(v)}
			filtered.SetGray(x, y, pixel)
		}
	}
	return filtered
}

//Benchmark_FilterGray	   27336411 ns/op
//Benchmark_FilterGrayFast 19521755 ns/op
//for better optimization in case of input gray image
func FilterGrayFast(grayImg *image.Gray, flt FilterType) (filtered *image.Gray) {
	max := grayImg.Bounds().Max
	min := grayImg.Bounds().Min

	/* filtered image must be two pixels shorter, because
	there must be a row of pixels on each side of a pixel for the sobel operator
	to work*/
	filtered = image.NewGray(image.Rect(max.X-2, max.Y-2, min.X, min.Y))
	width := max.X - 1 //to provide a "border" of 1 pixel
	height := max.Y - 1

	var v uint32
	applay := getFilterFunc(flt)

	for x := 1; x < width; x++ {
		for y := 1; y < height; y++ {
			fX, fY := applay(grayImg, x, y)
			v = FloorSqrt((fX*fX)+(fY*fY)) + 1 // +1 to make it ceil
			filtered.SetGray(x, y, color.Gray{Y: uint8(v)})
		}
	}

	return filtered
}

func FilterGrayMath(grayImg *image.Gray) (filtered *image.Gray) {
	max := grayImg.Bounds().Max
	min := grayImg.Bounds().Min

	/* filtered image must be two pixels shorter, because
	there must be a row of pixels on each side of a pixel for the sobel operator
	to work*/
	filtered = image.NewGray(image.Rect(max.X-2, max.Y-2, min.X, min.Y))
	width := max.X - 1 //to provide a "border" of 1 pixel
	height := max.Y - 1

	for x := 1; x < width; x++ {
		for y := 1; y < height; y++ {
			fX, fY := applySobelFilterMath(grayImg, x, y)
			fS := (fX * fX) + (fY * fY)
			//filtered.SetGray(x, y, color.Gray{Y: uint8(v)})
			filtered.Pix[filtered.PixOffset(x-1, y-1)] = uint8(math.Sqrt(fS))
		}
	}

	return filtered
}

func applyLaplasianFilter(img *image.Gray, x int, y int) (uint32, uint32) {
	var fX, fY int
	curX := x - 1
	curY := y - 1
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			pixel := int(img.GrayAt(curX, curY).Y)
			fX += laplasianX[i][j] * pixel
			fY += laplasianY[i][j] * pixel
			curX++
		}
		curX = x - 1
		curY++
	}
	return Abs(fX), Abs(fY)
}

func applySobelFilter(img *image.Gray, x int, y int) (uint32, uint32) {
	var fX, fY, pixel int
	curX := x - 1
	curY := y - 1
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			pixel = int(img.GrayAt(curX, curY).Y)
			fX += sobelX[i][j] * pixel
			fY += sobelY[i][j] * pixel
			curX++
		}
		curX = x - 1
		curY++
	}
	return Abs(fX), Abs(fY)
}

//Benchmark_applySobelFilter	   	33.8 ns/op
//Benchmark_applySobelFilterFast   	19.0 ns/op
func applySobelFilterFast(img *image.Gray, x int, y int) (uint32, uint32) {
	var fX, fY, pixel, index int
	curX := x - 1
	curY := y - 1
	for i := 0; i < kernelSize; i++ {
		//index = i * kernelSize
		for j := 0; j < kernelSize; j++ {
			//it is unsafe but faster on 10% or so
			pixel = int(img.Pix[img.PixOffset(curX, curY)])
			fX += sobelXL[index] * pixel
			fY += sobelYL[index] * pixel
			curX++
			index++
		}
		curX = x - 1
		curY++
	}
	return Abs(fX), Abs(fY)
}

func applySobelFilterMath(img *image.Gray, x int, y int) (float64, float64) {
	var fX, fY, pixel, index int
	curX := x - 1
	curY := y - 1
	for i := 0; i < kernelSize; i++ {
		//index = i * kernelSize
		for j := 0; j < kernelSize; j++ {
			//it is unsafe but faster on 10% or so
			pixel = int(img.Pix[img.PixOffset(curX, curY)])
			fX += sobelXL[index] * pixel
			fY += sobelYL[index] * pixel
			curX++
			index++
		}
		curX = x - 1
		curY++
	}
	return math.Abs(float64(fX)), math.Abs(float64(fY))
}

func applySharaFilter(img *image.Gray, x int, y int) (uint32, uint32) {
	var fX, fY int
	curX := x - 1
	curY := y - 1
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			pixel := int(img.GrayAt(curX, curY).Y)
			fX += sharaX[i][j] * pixel
			fY += sharaY[i][j] * pixel
			curX++
		}
		curX = x - 1
		curY++
	}
	return Abs(fX), Abs(fY)
}

func applySharpenFilter(img *image.Gray, x int, y int) (uint32, uint32) {
	var fX, fY int
	curX := x - 1
	curY := y - 1
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			pixel := int(img.GrayAt(curX, curY).Y)
			fX += sharpenX[i][j] * pixel
			fY += sharpenY[i][j] * pixel
			curX++
		}
		curX = x - 1
		curY++
	}
	return Abs(fX), Abs(fY)
}

func ToGrayscale(img image.Image) *image.Gray {
	max := img.Bounds().Max
	min := img.Bounds().Min

	var filtered = image.NewGray(image.Rect(max.X, max.Y, min.X, min.Y))
	for x := 0; x < max.X; x++ {
		for y := 0; y < max.Y; y++ {
			grayColor := color.GrayModel.Convert(img.At(x, y))
			filtered.Set(x, y, grayColor)
		}
	}
	return filtered
}

// ISqrt returns floor(sqrt(n)). Typical run time is few hundreds of ns.
//https://gitlab.com/cznic/mathutil/-/blob/master/mathutil.go
func ISqrt(n uint32) (x uint32) {
	if n == 0 {
		return
	}

	if n >= math.MaxUint16*math.MaxUint16 {
		return math.MaxUint16
	}
	var px, nx uint32
	for x = n; ; px, x = x, nx {
		nx = (x + n/x) / 2
		if nx == x || nx == px {
			break
		}
	}
	return
}

//https://www.geeksforgeeks.org/square-root-of-an-integer/
//Time Complexity: O(Log x)
func FloorSqrt(x uint32) (ans uint32) {
	// Base Cases
	if x == 0 || x == 1 {
		return x
	}

	// Do Binary Search for floor(sqrt(x))
	var (
		start uint32 = 1
		mid   uint32
	)
	end := x

	for start <= end {
		mid = (start + end) / 2

		// If x is a perfect square
		if mid*mid == x {
			return mid
		}

		// Since we need floor, we update answer when mid*mid is
		// smaller than x, and move closer to sqrt(x)
		if mid*mid < x {

			start = mid + 1
			ans = mid
		} else { // If mid*mid is greater than x
			end = mid - 1
		}
	}
	return ans
}

//Note: The Binary Search can be further optimized to start with ‘start’ = 0 and ‘end’ = x/2.
//Floor of square root of x cannot be more than x/2 when x > 1.
//Benchmark_FloorSqrt-8       25.7 ns/op
//Benchmark_FloorSqrtFast-8   20.5 ns/op

func FloorSqrtFast(x uint32) (ans uint32) {
	// Base Cases
	if x == 0 || x == 1 {
		return x
	}

	// Do Binary Search for floor(sqrt(x))
	var (
		start uint32 = 0
		mid   uint32
	)
	end := x / 2

	for start <= end {
		mid = (start + end) / 2

		// If x is a perfect square
		if mid*mid == x {
			return mid
		}

		// Since we need floor, we update answer when mid*mid is
		// smaller than x, and move closer to sqrt(x)
		if mid*mid < x {

			start = mid + 1
			ans = mid
		} else { // If mid*mid is greater than x
			end = mid - 1
		}
	}
	return ans
}

// Abs returns the absolute value of the given int.
func Abs(x int) uint32 {
	if x < 0 {
		return uint32(-x)
	} else {
		return uint32(x)
	}
}

//BenchmarkIT/Benchmark_FilterGraySimd-2         	    1032	   1499811 ns/op malloc
//
//BenchmarkIT/Benchmark_FilterGraySimd-2         	     501	   2000184 ns/op
// dstX := make([]uint16, dstSize) dstXC := (*C.uint8_t)(unsafe.Pointer(&dstX[0]))
//
//BenchmarkIT/Benchmark_FilterGraySimd-2         	     276	   4853208 ns/op
//add flited filling
func FilterGraySimd(grayImg *image.Gray) (filtered *image.Gray) {
	max := grayImg.Bounds().Max
	min := grayImg.Bounds().Min
	filtered = image.NewGray(image.Rect(max.X, max.Y, min.X, min.Y))
	imSize := len(grayImg.Pix)

	// void SimdSobelDyAbs (
	//    [in]	src	- a pointer to pixels data of the input image.
	//    [in]	srcStride	- a row size of the input image.
	//    [in]	width	- an image width.
	//    [in]	height	- an image height.
	//    [out]	dst	- a pointer to pixels data of the output image.
	//    [in]	dstStride	- a row size of the output image (in bytes).
	// )
	//All images must have the same width and height. Input image
	// must has 8-bit gray format, output image must has 16-bit integer format.
	src := (*C.uint8_t)(unsafe.Pointer(&grayImg.Pix[0]))
	srcStride := C.size_t(grayImg.Stride)
	width := C.size_t(max.X)
	height := C.size_t(max.Y)
	dstStride := srcStride * 2
	dstSize := C.size_t(imSize * 2)

	dstXC := (*C.uint8_t)(unsafe.Pointer(C.malloc(dstSize)))
	defer C.free(unsafe.Pointer(dstXC))
	C.SimdSobelDxAbs(src, srcStride, width, height, dstXC, dstStride)

	dstYC := (*C.uint8_t)(unsafe.Pointer(C.malloc(dstSize)))
	defer C.free(unsafe.Pointer(dstYC))

	C.SimdSobelDyAbs(src, srcStride, width, height, dstYC, dstStride)

	dstX := nonCopyGoUint16(uintptr(unsafe.Pointer(dstXC)), imSize)
	dstY := nonCopyGoUint16(uintptr(unsafe.Pointer(dstYC)), imSize)

	var fX, fY uint
	var fS float64
	var pix uint8
	if len(dstX) == imSize && len(dstY) == imSize && len(filtered.Pix) == imSize {
		for i := 0; i < imSize; i++ {
			fX = uint(dstX[i])
			fY = uint(dstY[i])
			fS = math.Sqrt(float64(fX*fX + fY*fY))
			//clipping
			if fS > 255.0 {
				pix = 255
			} else {
				pix = uint8(fS)
			}
			filtered.Pix[i] = pix
		}
	}
	return filtered
}

func nonCopyGoUint16(ptr uintptr, length int) []uint16 {
	var slice []uint16
	header := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	header.Cap = length
	header.Len = length
	header.Data = ptr
	return slice
}
