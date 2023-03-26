package snok

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/schollz/progressbar/v3"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"wsnok/internal/printline"
	"wsnok/internal/uniquefilename"
)

var Debug = false
var Threads = 5
var ChunkSize int64 = 1024 * 1024
var _chunks []*Chunk
var wgDown sync.WaitGroup
var wgWrite sync.WaitGroup
var progressBar *progressbar.ProgressBar = nil
var OverWrite = false
var Username = ""
var Password = ""

// Snok downloads a file from a URL
func Snok(url string) error {
	destinationFile := uniquefilename.GetUniqueFilenameFromUrl(url, OverWrite)
	printline.Debug = Debug
	start := time.Now()
	allowsChunks, size, contentType, err := head(url)
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

	printline.Printf(false, "HTTP request sent. Length: %d(%s), Threads: %d, ChunkSize: %d,Content-Type: %s\n", size, humanize.Bytes(uint64(size)), Threads, ChunkSize, contentType)
	printline.Printf(false, "Saving to: %s\n", destinationFile)

	work := make(chan *Chunk)
	writeChannel := make(chan *Chunk)

	go writeChunkChannel(destinationFile, writeChannel, &wgWrite)
	for i := 0; i < Threads; i++ {
		wgDown.Add(1)
		go worker(i, work, &wgDown, writeChannel)
	}

	progressBar = progressbar.NewOptions64(size,
		progressbar.OptionSetDescription("downloading"),
		progressbar.OptionSetWidth(10),
		progressbar.OptionFullWidth(),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetVisibility(!Debug),
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
	wgDown.Wait()
	close(writeChannel)
	wgWrite.Wait()
	elapsed := time.Since(start)
	printline.Printf(false, "Download took %s\n", elapsed)
	printline.Printf(false, "In Chunks: %v, Total Size: %v\n", allowsChunks, humanize.Bytes(uint64(size)))
	return nil
}

func worker(id int, work chan *Chunk, wg *sync.WaitGroup, write chan *Chunk) {
	for job := range work {
		printline.Printf(true, "Downloading chunk %d:%+v\n", id, job)
		err := downloadRange(job)
		if err != nil {
			printline.Printf(false, "Error downloading chunk: %s", err)
			os.Exit(1)
		}
		write <- job
	}
	wg.Done()
}

func writeChunkChannel(filename string, chunk chan *Chunk, wg *sync.WaitGroup) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		printline.Printf(false, "Error opening file: %s", err)
		os.Exit(1)
	}
	defer closeCheck(file)
	wg.Add(1)

	go func() {
		for chunk := range chunk {
			printline.Printf(true, "Writing chunk %d\n", chunk.Index)
			_, err := file.WriteAt(chunk.Data, chunk.Offset)
			if err != nil {
				printline.Printf(false, "Error writing chunk: %s", err)
				os.Exit(1)
			}
			if Debug {
				err := os.WriteFile(filename+"."+strconv.Itoa(chunk.Index), chunk.Data, 0644)
				if err != nil {
					printline.Printf(false, "Error writing chunk test file: %s", err)
					os.Exit(1)
				}
			}
			chunk.Data = nil

		}
		wg.Done()
	}()
	wg.Wait()
}

func calculateChunks(url string, size int64) []*Chunk {
	count := int(size / ChunkSize)
	var offset int64 = 0
	if size%ChunkSize != 0 {
		count++
	}
	chunks := make([]*Chunk, count)
	for i := 0; i < count; i++ {
		currentSize := ChunkSize
		if offset+ChunkSize > size {
			currentSize = size - ChunkSize
		}
		chunks[i] = &Chunk{
			Index:  i,
			Url:    url,
			Offset: offset,
			Size:   currentSize,
		}
		offset += currentSize

	}
	return chunks
}

func head(url string) (allowsChunks bool, fileBytes int64, contentType string, err error) {
	fileBytes = 0
	contentType = ""
	allowsChunks = false
	resp, err := http.Head(url)
	if err != nil {
		return
	}

	for name, value := range resp.Header {
		switch strings.ToLower(name) {
		case "accept-ranges":
			if strings.ToLower(value[0]) == "bytes" {
				allowsChunks = true
			}
		case "content-length":
			fileBytes, err = strconv.ParseInt(value[0], 10, 64)
			if err != nil {
				return
			}
		case "content-type":
			contentType = value[0]
		}
	}

	return
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
