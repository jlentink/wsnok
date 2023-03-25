package stringtoint

import (
	"testing"
)

func TestParse(t *testing.T) {
	type args struct {
		s string
	}

	type want struct {
		bytes int64
		err   bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Test 1 Byte",
			args: args{
				s: "1b",
			},
			want: want{
				bytes: 1,
				err:   false,
			},
		},
		{
			name: "Test 1 KiloByte",
			args: args{
				s: "1k",
			},
			want: want{
				bytes: 1024,
				err:   false,
			},
		},
		{
			name: "Test 1 MegaByte",
			args: args{
				s: "1M",
			},
			want: want{
				bytes: 1024 * 1024,
				err:   false,
			},
		},
		{
			name: "Test 1 GigaByte",
			args: args{
				s: "1G",
			},
			want: want{
				bytes: 1024 * 1024 * 1024,
				err:   false,
			},
		},
		{
			name: "Test 1 TeraByte",
			args: args{
				s: "1T",
			},
			want: want{
				bytes: 1024 * 1024 * 1024 * 1024,
				err:   false,
			},
		},
		{
			name: "simple bytes",
			args: args{
				s: "512",
			},
			want: want{
				bytes: 512,
				err:   false,
			},
		},
		{
			name: "Unknown unit",
			args: args{
				s: "512R",
			},
			want: want{
				bytes: -1,
				err:   true,
			},
		},
		{
			name: "Unknown unit",
			args: args{
				s: "abc",
			},
			want: want{
				bytes: -1,
				err:   true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.s)
			if got != tt.want.bytes {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
			if (err != nil) != tt.want.err {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.want.err)
			}
		})
	}
}

//
//func Test_calculateBytes(t *testing.T) {
//	type args struct {
//		s          string
//		multiplier int
//	}
//	tests := []struct {
//		name string
//		args args
//		want int64
//	}{
//		{
//			name: "Test 1 Byte",
//			args: args{
//				s:          "1b",
//				multiplier: 1,
//			},
//			want: 1,
//		},
//		{
//			name: "Test 1 KiloByte",
//			args: args{
//				s:          "1k",
//				multiplier: 1024,
//			},
//			want: 1024,
//		},
//		{
//			name: "Test 1 MegaByte",
//			args: args{
//				s:          "1m",
//				multiplier: 1024 * 1024,
//			},
//			want: 1024 * 1024,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := calculateBytes(tt.args.s, tt.args.multiplier); got != tt.want {
//				t.Errorf("calculateBytes() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
