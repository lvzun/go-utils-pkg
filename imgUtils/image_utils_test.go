package imgUtils

import (
	"context"
	"io/ioutil"
	"reflect"
	"testing"
)

func Test_baiduImageCute(t *testing.T) {
	file, err := ioutil.ReadFile("a.bmp")
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		ctx      context.Context
		imgBytes []byte
		name     string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "test1",
			args: args{
				ctx:      context.Background(),
				imgBytes: file,
				name:     "8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//file, _ := ioutil.ReadFile("upload/20221021/1234_170803.png")
			//toString := base64.StdEncoding.EncodeToString(file)
			//t.Logf("img:%s", toString)
			if got := baiduImageCute(tt.args.ctx, tt.args.imgBytes, tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("baiduImageCute() = %v, want %v", got, tt.want)
			}
		})
	}
}
