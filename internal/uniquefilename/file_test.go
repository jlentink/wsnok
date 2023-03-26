package uniquefilename

import (
	"reflect"
	"testing"
)

func TestFile_NextName(t *testing.T) {
	type fields struct {
		Increment int
		Filename  string
		Name      string
		Extension string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "TestFile_NextName_1",
			fields: fields{
				Increment: 0,
				Filename:  "test.txt",
				Name:      "test",
				Extension: "txt",
			},
			want: "test1.txt",
		},
		{
			name: "TestFile_NextName_2",
			fields: fields{
				Increment: 1,
				Filename:  "test.txt",
				Name:      "test",
				Extension: "txt",
			},
			want: "test2.txt",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{
				Increment: tt.fields.Increment,
				Filename:  tt.fields.Filename,
				Name:      tt.fields.Name,
				Extension: tt.fields.Extension,
			}
			if got := f.NextName(); got != tt.want {
				t.Errorf("NextName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewFile(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		name string
		args args
		want *File
	}{
		{
			name: "TestNewFile",
			args: args{fileName: "test.txt"},
			want: &File{
				Increment: 0,
				Filename:  "test.txt",
				Name:      "test",
				Extension: "txt",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFile(tt.args.fileName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
