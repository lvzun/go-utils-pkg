package imgUtils

import (
	"errors"
	"fmt"
	"github.com/Comdex/imgo"
	"github.com/golang/freetype"
	"github.com/lvzun/go-utils-pkg/imgUtils/textTransImage"
	"github.com/sirupsen/logrus"
	"golang.org/x/image/bmp"
	"golang.org/x/image/webp"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"math"
	"net/http"
	"os"
	"sort"
	"strings"
)

func New2DSlice(x int, y int) (theSlice [][]uint8) {
	theSlice = make([][]uint8, y, y)
	for i := 0; i < y; i++ {
		s2 := make([]uint8, x, x)
		theSlice[i] = s2
	}
	return
}

func New3DSlice(x int, y int, z int) (theSlice [][][]uint8) {
	theSlice = make([][][]uint8, x, x)
	for i := 0; i < x; i++ {
		s2 := make([][]uint8, y, y)
		for j := 0; j < y; j++ {
			s3 := make([]uint8, z, z)
			s2[j] = s3
		}
		theSlice[i] = s2
	}
	return
}

func LoadImage(filePath string) (img image.Image, err error) {
	f1Src, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}
	defer f1Src.Close()

	buff := make([]byte, 512) // why 512 bytes ? see http://golang.org/pkg/net/http/#DetectContentType
	_, err = f1Src.Read(buff)

	if err != nil {
		return nil, err
	}

	fileType := http.DetectContentType(buff)

	fmt.Println(fileType)

	fSrc, err := os.Open(filePath)
	defer fSrc.Close()

	switch fileType {
	case "image/jpeg", "image/jpg":
		img, err = jpeg.Decode(fSrc)
		if err != nil {
			fmt.Println("jpeg error")
			return nil, err
		}

	case "image/gif":
		img, err = gif.Decode(fSrc)
		if err != nil {
			return nil, err
		}

	case "image/png":
		img, err = png.Decode(fSrc)
		if err != nil {
			return nil, err
		}
	case "image/webp":
		img, err = webp.Decode(fSrc)
		if err != nil {
			return nil, err
		}
	case "image/bmp":
		img, err = bmp.Decode(fSrc)
		if err != nil {
			return nil, err
		}
	default:
		return nil, err
	}
	return img, nil
}

func ImageSave(fileName string, rgba *image.NRGBA) error {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logrus.Errorf("Open err:%v", err)
		return err
	}
	file.Name()
	// 将图片和扩展名分离
	stringSlice := strings.Split(file.Name(), ".")

	// 根据图片的扩展名来运用不同的处理
	switch stringSlice[len(stringSlice)-1] {
	case "jpg":
		return jpeg.Encode(file, rgba, nil)
	case "jpeg":
		return jpeg.Encode(file, rgba, nil)
	case "gif":
		return gif.Encode(file, rgba, nil)
	case "png":
		return png.Encode(file, rgba)
	default:
		panic("不支持的图片类型")
	}
}
func GrayingImage(m *image.NRGBA) *image.NRGBA {
	bounds := m.Bounds()

	dx := bounds.Dx()

	dy := bounds.Dy()

	newRgba := image.NewNRGBA(bounds)

	for x := 0; x < dx; x++ {

		for y := 0; y < dy; y++ {

			colorRgb := m.At(x, y)

			r, g, b, _ := colorRgb.RGBA()
			//avg := float64(0.3)*float64(uint8(r)) + float64(0.59)*float64(uint8(g)) + float64(0.11)*float64(uint8(b))
			avg := (uint8(r) + uint8(g) + uint8(b)) / 3

			// 将每个点的设置为灰度值
			newRgba.SetNRGBA(x, y, color.NRGBA{R: uint8(avg), G: uint8(avg), B: uint8(avg), A: 255})
		}
	}

	return newRgba
}

func Copy(src [][][]uint8, dst [][][]uint8, sx, sy, dx, dy, sw, sh int) {
	for i := 0; i < sh; i++ {
		for j := 0; j < sw; j++ {
			dst[dy+i][dx+j][0] = src[sy+i][sx+j][0]
			dst[dy+i][dx+j][1] = src[sy+i][sx+j][1]
			dst[dy+i][dx+j][2] = src[sy+i][sx+j][2]
			dst[dy+i][dx+j][3] = src[sy+i][sx+j][3]
		}
	}
}

func GetHistGram(src [][][]uint8) []int {
	height := len(src)
	width := len(src[0])
	histGram := make([]int, 256)
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			histGram[int(src[i][j][0])]++
		}
	}
	return histGram
}

func BinaryThreshold(image *image.NRGBA) uint8 {
	bounds := image.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	pixCount := 0
	rgbNum := uint64(0)
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			r, g, b, _ := image.At(j, i).RGBA()
			var rgb = (uint8(r) + uint8(g) + uint8(b)) / 3
			pixCount++
			rgbNum += uint64(rgb)
		}
	}
	return uint8(rgbNum / uint64(pixCount))
}

func SplitHorizontal(theSlices [][][]uint8, threshold uint8) [][][]uint8 {
	height := len(theSlices)
	width := len(theSlices[0])
	sameList := make([]int, width)
	for i := 1; i < width; i++ {
		sameValue := 0
		for j := 0; j < height; j++ {
			if theSlices[j][i-1][0] == theSlices[j][i][0] {
				sameValue++
			}
		}
		sameList[i] = sameValue
	}
	//fmt.Printf("sameList:%v", sameList)
	return theSlices
}

func RGB2Gray(fileName string) {
	imgMatrix := imgo.MustRead(fileName) //获取一个[][][]uint8对象
	//binaryzation process of image matrix , threshold can use 127 to test
	//func Binaryzation(src [][][]uint8, threshold int) [][][]uint8{}
	imgMatrix_gray := imgo.RGB2Gray(imgMatrix)

	gray_rgb := "gray_rgb.png"
	err := imgo.SaveAsPNG(gray_rgb, imgMatrix_gray)
	if err != nil {
		fmt.Println(err)
	}
}

func ConvertToNRGBA(src image.Image) *image.NRGBA {
	srcBounds := src.Bounds()
	dstBounds := srcBounds.Sub(srcBounds.Min)

	dst := image.NewNRGBA(dstBounds)

	dstMinX := dstBounds.Min.X
	dstMinY := dstBounds.Min.Y

	srcMinX := srcBounds.Min.X
	srcMinY := srcBounds.Min.Y
	srcMaxX := srcBounds.Max.X
	srcMaxY := srcBounds.Max.Y

	switch src0 := src.(type) {

	case *image.NRGBA:
		rowSize := srcBounds.Dx() * 4
		numRows := srcBounds.Dy()

		i0 := dst.PixOffset(dstMinX, dstMinY)
		j0 := src0.PixOffset(srcMinX, srcMinY)

		di := dst.Stride
		dj := src0.Stride

		for row := 0; row < numRows; row++ {
			copy(dst.Pix[i0:i0+rowSize], src0.Pix[j0:j0+rowSize])
			i0 += di
			j0 += dj
		}

	case *image.NRGBA64:
		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				j := src0.PixOffset(x, y)

				dst.Pix[i+0] = src0.Pix[j+0]
				dst.Pix[i+1] = src0.Pix[j+2]
				dst.Pix[i+2] = src0.Pix[j+4]
				dst.Pix[i+3] = src0.Pix[j+6]

			}
		}

	case *image.RGBA:
		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				j := src0.PixOffset(x, y)
				a := src0.Pix[j+3]
				dst.Pix[i+3] = a

				switch a {
				case 0:
					dst.Pix[i+0] = 0
					dst.Pix[i+1] = 0
					dst.Pix[i+2] = 0
				case 0xff:
					dst.Pix[i+0] = src0.Pix[j+0]
					dst.Pix[i+1] = src0.Pix[j+1]
					dst.Pix[i+2] = src0.Pix[j+2]
				default:
					dst.Pix[i+0] = uint8(uint16(src0.Pix[j+0]) * 0xff / uint16(a))
					dst.Pix[i+1] = uint8(uint16(src0.Pix[j+1]) * 0xff / uint16(a))
					dst.Pix[i+2] = uint8(uint16(src0.Pix[j+2]) * 0xff / uint16(a))
				}
			}
		}

	case *image.RGBA64:
		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				j := src0.PixOffset(x, y)
				a := src0.Pix[j+6]
				dst.Pix[i+3] = a

				switch a {
				case 0:
					dst.Pix[i+0] = 0
					dst.Pix[i+1] = 0
					dst.Pix[i+2] = 0
				case 0xff:
					dst.Pix[i+0] = src0.Pix[j+0]
					dst.Pix[i+1] = src0.Pix[j+2]
					dst.Pix[i+2] = src0.Pix[j+4]
				default:
					dst.Pix[i+0] = uint8(uint16(src0.Pix[j+0]) * 0xff / uint16(a))
					dst.Pix[i+1] = uint8(uint16(src0.Pix[j+2]) * 0xff / uint16(a))
					dst.Pix[i+2] = uint8(uint16(src0.Pix[j+4]) * 0xff / uint16(a))
				}
			}
		}

	case *image.Gray:
		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				j := src0.PixOffset(x, y)
				c := src0.Pix[j]
				dst.Pix[i+0] = c
				dst.Pix[i+1] = c
				dst.Pix[i+2] = c
				dst.Pix[i+3] = 0xff

			}
		}

	case *image.Gray16:
		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				j := src0.PixOffset(x, y)
				c := src0.Pix[j]
				dst.Pix[i+0] = c
				dst.Pix[i+1] = c
				dst.Pix[i+2] = c
				dst.Pix[i+3] = 0xff

			}
		}

	case *image.YCbCr:
		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				yj := src0.YOffset(x, y)
				cj := src0.COffset(x, y)
				r, g, b := color.YCbCrToRGB(src0.Y[yj], src0.Cb[cj], src0.Cr[cj])

				dst.Pix[i+0] = r
				dst.Pix[i+1] = g
				dst.Pix[i+2] = b
				dst.Pix[i+3] = 0xff

			}
		}

	default:
		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				c := color.NRGBAModel.Convert(src.At(x, y)).(color.NRGBA)

				dst.Pix[i+0] = c.R
				dst.Pix[i+1] = c.G
				dst.Pix[i+2] = c.B
				dst.Pix[i+3] = c.A

			}
		}
	}

	return dst
}

// binaryzation process of image matrix , threshold can use 127 to test
func Binaryzation(src [][][]uint8, threshold int) [][][]uint8 {
	height := len(src)
	width := len(src[0])
	imgMatrix := imgo.NewRGBAMatrix(height, width)
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			var rgb int = (int(src[i][j][0]) + int(src[i][j][1]) + int(src[i][j][2])) / 3
			if rgb > threshold {
				rgb = 255
			} else {
				rgb = 0
			}
			imgMatrix[i][j][0] = uint8(rgb)
			imgMatrix[i][j][1] = uint8(rgb)
			imgMatrix[i][j][2] = uint8(rgb)
			imgMatrix[i][j][3] = src[i][j][3]
		}
	}
	return imgMatrix
}

func DeNoise(imgMatrix [][][]uint8) [][][]uint8 {
	height := len(imgMatrix)
	width := len(imgMatrix[0])
	for i := 0; i < height; i += 3 {
		for j := 0; j < width; j += 3 {
			count := D8(imgMatrix, j, i)
			if count <= 2 && count > 0 {
				clearRect(imgMatrix, Rect{X: j - 1, Y: i - 1, Width: 3, Height: 3})
			}
		}
	}
	return imgMatrix
}

func D8(imgMatrix [][][]uint8, x, y int) int {
	height := len(imgMatrix)
	width := len(imgMatrix[0])

	count := 0
	for i := 0; i < 2; i++ {
		if y+i-1 < 0 || y+i-1 >= height {
			continue
		}
		for j := 0; j < 2; j++ {
			if x+j-1 < 0 || x+j-1 >= width {
				continue
			}
			if imgMatrix[y+i-1][x+j-1][0] != 0 {
				count++
			}
		}
	}
	return count
}
func BinaryReverse(theSlices [][][]uint8) [][][]uint8 {
	height := len(theSlices)
	width := len(theSlices[0])
	for i := 1; i < width; i++ {
		for j := 0; j < height; j++ {
			if theSlices[j][i][0] == 0 {
				theSlices[j][i][0] = 255
				theSlices[j][i][1] = 255
				theSlices[j][i][2] = 255
			} else {
				theSlices[j][i][0] = 0
				theSlices[j][i][1] = 0
				theSlices[j][i][2] = 0
			}
		}
	}
	return theSlices
}
func BinaryProportion(theSlices [][][]uint8) int {
	height := len(theSlices)
	width := len(theSlices[0])
	white := 0
	for i := 1; i < width; i++ {
		for j := 0; j < height; j++ {
			if theSlices[j][i][0] != 0 {
				white++
			}
		}
	}
	return int(float64(white) * 100.0 / float64(height*width))
}

type Rect struct {
	X      int
	Y      int
	Width  int
	Height int
}

func SubImageSearch(theSlices [][][]uint8, threshold int) []Rect {
	height := len(theSlices)
	width := len(theSlices[0])
	rects := make([]Rect, 0)
	imgMatrix := imgo.NewRGBAMatrix(height, width)
	Copy(theSlices, imgMatrix, 0, 0, 0, 0, width, height)

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			if imgMatrix[i][j][0] != 0 {
				rect := GetRect(imgMatrix, i, j, threshold) //获取非透明像素点所在的矩形区域
				if rect.Width > 5 && rect.Height > 5 {      //剔除尺寸小于10x10的子图区域
					rects = append(rects, rect)
				}
				clearRect(imgMatrix, rect)
			}
		}
	}

	_ = imgMatrix
	return rects
}

func GetRect(theSlices [][][]uint8, i int, j, threshold int) Rect {
	r := Rect{X: j, Y: i, Width: 1, Height: 1}

	flag := true
	for {
		if !flag {
			break
		}
		flag = false
		for {
			if findRight(theSlices, r, threshold) {
				r.Width++
				flag = true
			} else {
				break
			}
		}
		for {
			if findDown(theSlices, r, threshold) {
				r.Height++
				flag = true
			} else {
				break
			}
		}
		for {
			if findLeft(theSlices, r, threshold) {
				r.Width++
				r.X--
				flag = true
			} else {
				break
			}
		}
		for {
			if findUp(theSlices, r, threshold) {
				r.Height++
				r.Y--
				flag = true
			} else {
				break
			}
		}
	}

	r.Width++
	r.Height++
	return r
}

func findRight(uint8s [][][]uint8, rect Rect, threshold int) bool {
	if rect.X+rect.Width >= len(uint8s[0]) || rect.X < 0 {
		return false
	}
	for i := -1; i <= rect.Height; i++ {
		if rect.Y+i < 0 || rect.Y+i > len(uint8s) {
			continue
		}
		for j := 1; j <= threshold; j++ {
			if rect.X+rect.Width+j >= len(uint8s[0]) {
				continue
			}
			if uint8s[rect.Y+i][rect.X+rect.Width+j][0] != 0 {
				return true
			}
		}
	}
	return false
}

func findDown(uint8s [][][]uint8, rect Rect, threshold int) bool {
	if rect.Y+rect.Height >= len(uint8s) || rect.Y < 0 {
		return false
	}
	for i := -1; i <= rect.Width; i++ {
		if rect.X+i < 0 || rect.X+i > len(uint8s[0]) {
			continue
		}
		for j := 1; j <= threshold; j++ {
			if rect.Y+rect.Height+j >= len(uint8s) {
				continue
			}
			if uint8s[rect.Y+rect.Height+j][rect.X+i][0] != 0 {
				return true
			}
		}

	}
	return false
}
func findLeft(uint8s [][][]uint8, rect Rect, threshold int) bool {
	if rect.X+rect.Width >= len(uint8s[0]) || rect.X < 0 {
		return false
	}
	for i := -1; i <= rect.Height; i++ {
		if rect.Y+i < 0 || rect.Y+i > len(uint8s) {
			continue
		}
		for j := 1; j <= threshold; j++ {
			if rect.X-j < 0 {
				continue
			}
			if uint8s[rect.Y+i][rect.X-j][0] != 0 {
				return true
			}
		}
	}
	return false
}
func findUp(uint8s [][][]uint8, rect Rect, threshold int) bool {
	if rect.Y+rect.Height >= len(uint8s) || rect.Y < 0 {
		return false
	}
	for i := 0; i < rect.Width; i++ {
		if rect.X+i < 0 || rect.X+i > len(uint8s[0]) {
			continue
		}
		for j := 1; j <= threshold; j++ {
			if rect.Y-j < 0 {
				continue
			}
			if uint8s[rect.Y-j][rect.X+i][0] != 0 {
				return true
			}
		}
	}
	return false
}

func clearRect(theSlices [][][]uint8, rect Rect) {
	height := len(theSlices)
	width := len(theSlices[0])
	for i := rect.Y; i < rect.Y+rect.Height; i++ {
		if i < 0 || i >= height {
			continue
		}
		for j := rect.X; j < rect.X+rect.Width; j++ {
			if j < 0 || j >= width {
				continue
			}
			theSlices[i][j][0] = 0
			theSlices[i][j][1] = 0
			theSlices[i][j][2] = 0
			theSlices[i][j][3] = 0
		}
	}
}

/**
 * 创建任意角度的旋转图像
 * @param image
 * @param theta
 * @param backgroundColor
 * @return
 */
func RotateImage(image [][][]uint8, theta float64, backgroundRgb uint8) [][][]uint8 {
	height := len(image)
	width := len(image[0])
	angle := theta * math.Pi / 180 // 度转弧度
	xCoords := getX(width/2, height/2, angle)
	yCoords := getY(width/2, height/2, angle)
	newWidth := int(xCoords[3] - xCoords[0])
	newHeight := int(yCoords[3] - yCoords[0])
	resultImage := imgo.NewRGBAMatrix(newHeight, newWidth)

	for i := 0; i < newWidth; i++ {
		for j := 0; j < newHeight; j++ {
			x := i - newWidth/2
			y := newHeight/2 - j
			radius := math.Sqrt(float64(x*x + y*y))
			angle1 := float64(0)
			if y > 0 {
				angle1 = math.Acos(float64(x) / radius)
			} else {
				angle1 = 2*math.Pi - math.Acos(float64(x)/radius)
			}
			x = int(radius * math.Cos(angle1-angle))
			y = int(radius * math.Sin(angle1-angle))
			if x < (width/2) && x > -(width/2) && y < (height/2) && y > -(height/2) {
				resultImage[j][i][0] = image[height/2-y][x+width/2][0]
				resultImage[j][i][1] = image[height/2-y][x+width/2][1]
				resultImage[j][i][2] = image[height/2-y][x+width/2][2]
				resultImage[j][i][3] = image[height/2-y][x+width/2][3]
			} else {

				resultImage[j][i][0] = backgroundRgb
				resultImage[j][i][1] = backgroundRgb
				resultImage[j][i][2] = backgroundRgb
				resultImage[j][i][3] = backgroundRgb
			}
		}
	}
	return resultImage
}

// 获取四个角点旋转后Y方向坐标
func getY(i, j int, angle float64) []float64 {
	results := make([]float64, 4)
	radius := math.Sqrt(float64(i*i + j*j))
	angle1 := math.Asin(float64(j) / radius)
	results[0] = radius * math.Sin(angle1+angle)
	results[1] = radius * math.Sin(math.Pi-angle1+angle)
	results[2] = -results[0]
	results[3] = -results[1]
	sort.Sort(sort.Float64Slice(results))
	return results
}

// 获取四个角点旋转后X方向坐标
func getX(i, j int, angle float64) []float64 {
	results := make([]float64, 4)
	radius := math.Sqrt(float64(i*i + j*j))
	angle1 := math.Acos(float64(i) / radius)
	results[0] = radius * math.Cos(angle1+angle)
	results[1] = radius * math.Cos(math.Pi-angle1+angle)
	results[2] = -results[0]
	results[3] = -results[1]
	sort.Sort(sort.Float64Slice(results))
	return results
}

func Text2Image(ttfName string, text string, width, height int) ([][][]uint8, int) {

	size := width
	if size > height {
		size = height
	}
	rgbaMatrix := imgo.NewRGBAMatrix(height, width)
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	file := textTransImage.New(ttfName)
	font := file.GetFont()

	if font == nil {
		logrus.Error("加载字体出错")
		return rgbaMatrix, 0
	}
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(font)
	c.SetFontSize(float64(size))
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.White)
	//设置字体显示位置
	pt := freetype.Pt(0, (height+1)/2+int(c.PointToFixed(float64(size))>>8))
	_, err := c.DrawString(text, pt)
	if err != nil {
		logrus.Error(err)
		return rgbaMatrix, 0
	}

	count := 0
	src := ConvertToNRGBA(img)
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			c := src.At(j, i)
			r, g, b, _ := c.RGBA()
			if r+g+b > 0 {
				count++
				rgbaMatrix[i][j][0] = uint8(255)
				rgbaMatrix[i][j][1] = uint8(255)
				rgbaMatrix[i][j][2] = uint8(255)
				rgbaMatrix[i][j][3] = uint8(255)
			}
		}
	}

	return rgbaMatrix, count
}

func MatrixToNRGB(imgMatrix [][][]uint8) (*image.NRGBA, error) {
	height := len(imgMatrix)
	width := len(imgMatrix[0])

	if height == 0 || width == 0 {
		return nil, errors.New("The input of matrix is illegal!")
	}

	nrgba := image.NewNRGBA(image.Rect(0, 0, width, height))

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			nrgba.SetNRGBA(j, i, color.NRGBA{imgMatrix[i][j][0], imgMatrix[i][j][1], imgMatrix[i][j][2], imgMatrix[i][j][3]})
		}
	}
	return nrgba, nil
}

func NRGBAToMatrix(src *image.NRGBA) [][][]uint8 {
	bounds := src.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y
	imgMatrix := imgo.NewRGBAMatrix(height, width)

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			c := src.At(j, i)
			r, g, b, a := c.RGBA()
			imgMatrix[i][j][0] = uint8(r)
			imgMatrix[i][j][1] = uint8(g)
			imgMatrix[i][j][2] = uint8(b)
			imgMatrix[i][j][3] = uint8(a)

		}
	}
	return imgMatrix
}

//func Resize(matrix [][][]uint8,width,height int)[][][]uint8{
//	nrgba, err := MatrixToNRGB(matrix)
//	if err != nil {
//		fmt.Printf("MatrixToNRGB err: %s", err)
//		return matrix
//	}
//	resize := imgo.Resize(nrgba, width, height)
//	toMatrix := NRGBAToMatrix(resize)
//
//	return toMatrix
//}

func Resize(matrix [][][]uint8, dWidth, dHeight int) [][][]uint8 {
	sHeight := len(matrix)
	sWidth := len(matrix[0])

	imgMatrix := imgo.NewRGBAMatrix(dHeight, dWidth)

	for i := 0; i < dHeight; i++ {
		for j := 0; j < dWidth; j++ {
			si := sHeight * i / dHeight
			sj := sWidth * j / dWidth
			dst := matrix[si][sj][0]
			imgMatrix[i][j][0] = dst
			imgMatrix[i][j][1] = dst
			imgMatrix[i][j][2] = dst
			imgMatrix[i][j][3] = 255
		}
	}
	return imgMatrix
}
