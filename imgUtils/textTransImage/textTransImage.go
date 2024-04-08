package textTransImage

import (
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"io/ioutil"
	"log"
)

var Font *truetype.Font

func init() {
	file := New("3.ttf")
	Font = file.GetFont()
}

type FontFile struct {
	FontFileName string
}

func New(fontName string) *FontFile {
	f := &FontFile{FontFileName: fontName}
	return f
}
func (f *FontFile) GetFont() *truetype.Font {

	fontBytes, err := ioutil.ReadFile(f.FontFileName)
	if err != nil {
		log.Println(err)
		return nil
	}

	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		return nil
	}
	return font
}
