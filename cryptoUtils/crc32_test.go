package cryptoUtils

import "testing"

func TestCrc32(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want uint32
	}{
		{
			name: "test crc32",
			args: args{
				s: "hello world",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Crc32(tt.args.s); got != tt.want {
				t.Errorf("Crc32() = %v, want %v", got, tt.want)
			}
		})
	}
}
