package cryptoUtils

import (
	"encoding/hex"
	"reflect"
	"testing"
)

func TestSm4CbcPkcs5Decode(t *testing.T) {
	iv, err := hex.DecodeString("4F723F7349774F063C0C477A367B3278")
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
	key, err := hex.DecodeString("78494784947849478494784947849479")
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
	type args struct {
		data []byte
		key  []byte
		iv   []byte
	}
	data, err := hex.DecodeString("cd4242b70dce901ec1d998b077ef9c8e")
	_ = err
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "decode",
			args:    args{data, key, iv},
			want:    []byte("1"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Sm4CbcPkcs5Decode(tt.args.data, tt.args.key, tt.args.iv)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sm4CbcPkcs5Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sm4CbcPkcs5Decode() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSm4CbcPkcs5Encode(t *testing.T) {
	iv, err := hex.DecodeString("4F723F7349774F063C0C477A367B3278")
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
	key, err := hex.DecodeString("78494784947849478494784947849479")
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
	type args struct {
		data []byte
		key  []byte
		iv   []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "encdoe",
			args:    args{[]byte("12345678"), key, iv},
			want:    []byte("9e6663b6bff82b5c94e78e1d1367b56f"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Sm4CbcPkcs5Encode(tt.args.data, tt.args.key, tt.args.iv)
			t.Logf("got:%s", hex.EncodeToString(got))
			got = []byte(hex.EncodeToString(got))
			if (err != nil) != tt.wantErr {
				t.Errorf("Sm4CbcPkcs5Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sm4CbcPkcs5Encode() got = %v, want %v", got, tt.want)
			}
		})
	}
}
