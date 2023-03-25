package snok

import (
	"errors"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/schollz/progressbar/v3"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
	"wpull/internal/printline"
)

var Debug = false
var Threads = 5
var ChunkSize int64 = 1024 * 1024
var _chunks []*Chunk
var wg sync.WaitGroup
var progressBar *progressbar.ProgressBar = nil
var OverWrite = false
var Username, Password string = "", ""

// Filename finds a free filename to find or overwrite
func Filename(url string) string {
	filename := path.Base(url)

	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) || OverWrite {
		return filename
	}
	return getFilenameIncrement(filename, 0)
}

func getFilenameIncrement(filename string, increment int) string {
	increment++
	cFilename := filename + "." + strconv.Itoa(increment)
	if _, err := os.Stat(cFilename); errors.Is(err, os.ErrNotExist) {
		return cFilename
	}
	return getFilenameIncrement(filename, increment)
}

// Snok downloads a file from a URL
func Snok(url string) error {
	printline.Debug = Debug
	start := time.Now()
	allowsChunks, size, err := head(url)
	if err != nil {
		panic(err)
	}

	if !allowsChunks || Threads <= 1 {
		Threads = 1
		currentChunk := &Chunk{
			Index:  0,
			Url:    url,
			Offset: 0,
			Size:   size,
		}
		_chunks = append(_chunks, currentChunk)
	} else {
		_chunks = calculateChunks(url, size)
	}

	printline.Printf(false, "HTTP request sent. Length: %d(%s), Threads: %d\n", size, humanize.Bytes(uint64(size)), Threads)
	printline.Printf(false, "Saving to: %s\n", Filename(url))

	work := make(chan *Chunk)
	for i := 0; i < Threads; i++ {
		wg.Add(1)
		go worker(i, work, &wg)
	}

	progressBar = progressbar.NewOptions64(size,
		progressbar.OptionSetDescription("downloading"),
		progressbar.OptionSetWidth(10),
		progressbar.OptionFullWidth(),
		progressbar.OptionShowBytes(true),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	for _, chunk := range _chunks {
		work <- chunk
	}

	close(work)
	wg.Wait()

	writeChunks()
	elapsed := time.Since(start)
	printline.Printf(false, "Download took %s\n", elapsed)
	printline.Printf(false, "In Chunks: %v, Total Size: %v\n", allowsChunks, humanize.Bytes(uint64(size)))
	return nil
}

func worker(id int, work chan *Chunk, wg *sync.WaitGroup) {
	for job := range work {
		printline.Printf(true, "Downloading chunk %d:%+v\n", id, job)
		err := downloadRange(job)
		if err != nil {
			printline.Printf(false, "Error downloading chunk: %s", err)
			os.Exit(1)
		}
	}
	wg.Done()
}

func writeChunks() {
	filename := Filename(_chunks[0].Url)
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644)

	defer closeCheck(file)
	if err != nil {
		printline.Printf(false, "Error opening file: %s", err)
		return
	}
	for _, chunk := range _chunks {
		printline.Printf(true, "Writing chunk %+v\n", chunk.Index)
		_, err := file.Write(chunk.Data)
		if err != nil {
			printline.Printf(false, "Error writing chunk: %s", err)
			os.Exit(1)
		}
	}
}

func calculateChunks(url string, size int64) []*Chunk {
	count := int(size / ChunkSize)
	if size%ChunkSize != 0 {
		count++
	}
	chunks := make([]*Chunk, count)
	for i := 0; i < count; i++ {
		currentOffset := ChunkSize * int64(i)
		currentSize := ChunkSize
		if currentOffset+ChunkSize > size {
			currentSize = size - currentOffset
		}
		chunks[i] = &Chunk{
			Index:  i,
			Url:    url,
			Offset: currentOffset,
			Size:   currentSize,
		}
	}
	return chunks
}

func head(url string) (bool, int64, error) {
	var fileBytes int64 = 0
	var allowsChunks = false
	resp, err := http.Head(url)
	if err != nil {
		return allowsChunks, fileBytes, err
	}

	for name, value := range resp.Header {
		if strings.ToLower(name) == "accept-ranges" && strings.ToLower(value[0]) == "bytes" {
			allowsChunks = true
		}
		if strings.ToLower(name) == "content-length" {
			fileBytes, err = strconv.ParseInt(value[0], 10, 64)
			if err != nil {
				return allowsChunks, fileBytes, err
			}
		}
	}

	return allowsChunks, fileBytes, nil
}

func downloadRange(chunk *Chunk) error {
	req, err := http.NewRequest("GET", chunk.Url, nil)
	if err != nil {
		return err
	}

	if chunk.Offset > 0 || chunk.Size > 0 {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", chunk.Offset, chunk.Offset+chunk.Size-1))
	}

	if Username != "" || Password != "" {
		req.SetBasicAuth(Username, Password)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	chunk.Data = make([]byte, chunk.Size)
	_, err = io.Copy(io.MultiWriter(chunk, progressBar), resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func closeCheck(c io.Closer) {
	err := c.Close()
	if err != nil {
		printline.Printf(false, "Error closing: %s", err)
	}
}
