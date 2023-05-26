package cryptoUtils

import (
	"encoding/hex"
	"reflect"
	"testing"
)

func TestSm2Encrypt(t *testing.T) {
	type args struct {
		plaintext      []byte
		pubKeyBytes    []byte
		cipherTextType int
	}

	text, _ := hex.DecodeString("123456")

	//lsy, _ := hex.DecodeString("f288aaa6cc611bcb86bf4c963bf15eea13ecf96f9455209ecd2d3a2d698bb9f200fe839f8609bd8d0425c14fe024cdb8")
	//t.Logf("lsy:%s", string(lsy))
	//uni, _ := hex.DecodeString("662d648cfa88c52c195899d7982d418a0df4e73f46caceaabd79e41ef066dd21")
	//t.Logf("uni:%s", string(uni))

	want, _ := hex.DecodeString("04e65a89361648432cc5c58b75c3b344695a3c803c2b7e2992e82ed0d69d3c590cff503227f81bf5e0fa40bd29921a9945282133f667a153288a9707fd912060f3662d648cfa88c52c195899d7982d418a0df4e73f46caceaabd79e41ef066dd21877393")
	decodeString, err := hex.DecodeString("a915adfb0c6513f6c31398c9a3b72ed915e9f6fe141688a1caab09f8d1d54463b23b3136916b7b7b4c866d0a9de85f9af17fef7f01a6243ad231086e1ece95c8")
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "Test Sm2 Encrypt",
			args: args{
				//plaintext:      []byte(`{"password":"04278a9bda05e61a5cb0d8297817bb4b82cdd733c0cb39e6de7d418142ae755b80694c28aee61776e631f6f51e7fea13c251bc70300d8bb73f2a1757da04386899bb134b7f9853aef0401aaa44c637d4132c1c71f492584c1a8e30173e1843674c13b8e1d9ce229c14","smsCode":"916254"}`),
				plaintext:      text,
				pubKeyBytes:    decodeString,
				cipherTextType: 2,
			},
			want: want,
		},
	}
	//0629ad5b803a1ec85b250d35932fedb0d9c8cd326ec16136eb0667d244bdc334
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Sm2Encrypt(tt.args.plaintext, tt.args.pubKeyBytes, tt.args.cipherTextType)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sm2Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			toString := hex.EncodeToString(got)
			t.Logf("Sm2Encrypt() string = %s", toString)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sm2Encrypt() got = %v, want %v", got, tt.want)
			}
		})
	}
}
