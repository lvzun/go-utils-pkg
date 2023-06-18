package winApi

import "testing"

func TestMessageBox(t *testing.T) {
	type args struct {
		hwnd    uintptr
		caption string
		text    string
		flags   uint32
	}
	tests := []struct {
		name string
		args args
		want int32
	}{
		// TODO: Add test cases.
		{
			name: "test",
			args: args{
				hwnd:    0,
				caption: "test",
				text:    "test",
				flags:   0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MessageBox(tt.args.hwnd, tt.args.caption, tt.args.text, tt.args.flags); got != tt.want {
				t.Errorf("MessageBox() = %v, want %v", got, tt.want)
			}
		})
	}
}
