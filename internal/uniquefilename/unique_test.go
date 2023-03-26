package uniquefilename

import (
	"fmt"
	"github.com/spf13/afero"
	"os"
	"testing"
)

func setup() {
	fs = afero.NewMemMapFs()
}

func teardown() {
	fs = afero.NewOsFs()
}

func CreateFiles(files []string) {
	cwd, _ := os.Getwd()
	fs = afero.NewMemMapFs()
	afs := &afero.Afero{Fs: fs}
	fps, _ := afs.ReadDir(cwd)
	for _, file := range fps {
		afs.Remove(file.Name()) // nolint:errcheck
	}
	for _, file := range files {
		afs.WriteFile(file, []byte("test"), 0644) // nolint:errcheck
	}
	fps, _ = afs.ReadDir(cwd)
	for _, file := range fps {
		fmt.Println(file.Name())
	}
}

func TestGetUniqueFilename(t *testing.T) {
	type args struct {
		filename    string
		OverWrite   bool
		createFiles []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestGetUniqueFilename_Override",
			args: args{
				filename:    "test.txt",
				OverWrite:   true,
				createFiles: []string{"test.txt"},
			},
			want: "test.txt",
		},
		{
			name: "TestGetUniqueFilename_1",
			args: args{
				filename:    "test.txt",
				OverWrite:   false,
				createFiles: []string{"test.txt"},
			},
			want: "test1.txt",
		},
		{
			name: "TestGetUniqueFilename_4",
			args: args{
				filename:    "test.txt",
				OverWrite:   false,
				createFiles: []string{"test.txt", "test1.txt", "test2.txt", "test3.txt"},
			},
			want: "test4.txt",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup()
			CreateFiles(tt.args.createFiles)
			if got := GetUniqueFilename(tt.args.filename, tt.args.OverWrite); got != tt.want {
				t.Errorf("GetUniqueFilename() = %v, want %v", got, tt.want)
			}
			teardown()
		})
	}
}

func TestGetUniqueFilenameFromUrl(t *testing.T) {
	type args struct {
		url         string
		OverWrite   bool
		createFiles []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "URLTestGetUniqueFilename_Override",
			args: args{
				url:         "https://www.someurl.nl/test.txt",
				OverWrite:   true,
				createFiles: []string{"test.txt"},
			},
			want: "test.txt",
		},
		{
			name: "URLTestGetUniqueFilename_1",
			args: args{
				url:         "https://www.someurl.nl/test.txt",
				OverWrite:   false,
				createFiles: []string{"test.txt"},
			},
			want: "test1.txt",
		},
		{
			name: "URLTestGetUniqueFilename_4",
			args: args{
				url:         "https://www.someurl.nl/test.txt?bla=ddd",
				OverWrite:   false,
				createFiles: []string{"test.txt", "test1.txt", "test2.txt", "test3.txt"},
			},
			want: "test4.txt",
		},
		{
			name: "Domain no file name",
			args: args{
				url:         "https://www.someurl.nl/",
				OverWrite:   false,
				createFiles: []string{"test.txt", "test1.txt", "test2.txt", "test3.txt"},
			},
			want: "www.someurl.nl",
		},
		{
			name: "Folder name",
			args: args{
				url:         "https://www.someurl.nl/iamafolder/",
				OverWrite:   false,
				createFiles: []string{"test.txt", "test1.txt", "test2.txt", "test3.txt"},
			},
			want: "iamafolder",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup()
			CreateFiles(tt.args.createFiles)
			if got := GetUniqueFilenameFromUrl(tt.args.url, tt.args.OverWrite); got != tt.want {
				t.Errorf("GetUniqueFilenameFromUrl() = %v, want %v", got, tt.want)
			}
			teardown()
		})
	}
}
