package imgUtils

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Comdex/imgo"
	"github.com/sirupsen/logrus"
	"image"
)

func baiduImageCute(ctx context.Context, imgBytes []byte, name string) []byte {
	var imgMatrixResult [][][]uint8
	var err error
	if imgBytes != nil && len(imgBytes) > 0 {
		imgMatrixResult, err = decodeImageBytes(imgBytes)
		if err != nil {
			logrus.WithContext(ctx).Errorf("decodeImageBytes failed: %v", err)
		}
	}
	if len(imgMatrixResult) == 0 && len(name) > 0 {
		imgMatrixResult = imgo.MustRead("img/" + name + ".png")
	}
	//imgMatrixResult = imgo.RGB2Gray(imgMatrixResult)
	//imgo.SaveAsPNG("img/"+name+"_gray.png", imgMatrixResult)

	GetHistGram := GetHistGram(imgMatrixResult)
	fmt.Printf("GetHistGram:%v\n", GetHistGram)

	//resultHistGram := make([]int, len(GetHistGram))
	var imgMatrixBinary [][][]uint8
	threshold := GetMinimumThreshold(GetHistGram)
	fmt.Printf("threshold:%v\n", threshold)

	imgMatrixBinary = Binaryzation(imgMatrixResult, threshold)
	//imgMatrixBinary = DeNoise(imgMatrixBinary)
	imgo.SaveAsPNG(fmt.Sprintf("img/%s_binary.png", name), imgMatrixBinary)
	//imgMatrixBinary = BinaryReverse(imgMatrixBinary)
	//imgo.SaveAsPNG(fmt.Sprintf("img/%s_binary_reverse.png", name), imgMatrixBinary)

	search := SubImageSearch(imgMatrixBinary, 2)
	for _, rect := range search {
		if rect.Width < 100 {
			continue
		}
		if rect.Height < 800 {
			continue
		}
		logrus.Infof("rect:x:%d,y:%d:w:%d,h:%d", rect.X, rect.Y, rect.Width, rect.Height)
	}

	return nil
}

func decodeImageBytes(data []byte) (imgMatrix [][][]uint8, err error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	bounds := img.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y
	src := ConvertToNRGBA(img)
	imgMatrix = imgo.NewRGBAMatrix(height, width)
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
	return imgMatrix, nil
}
