package hls

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	"github.com/smallinsky/hls/pkg/fs"
	"github.com/smallinsky/hls/pkg/notify"
	"github.com/smallinsky/hls/pkg/video"
)

const (
	defaultNumOfSegmentationJobs = 5
	defaultJobQeueSize           = 20
)

func NewService(videoSrc, outSegmentDir, port string) *HLS {
	hls := &HLS{
		srcDir:      videoSrc,
		segmentsDir: outSegmentDir,
		port:        port,
		jobQueue:    make(chan string, defaultJobQeueSize),
		stop:        make(chan struct{}),
	}

	hls.startSegmentationJobs(defaultNumOfSegmentationJobs)
	return hls
}

type HLS struct {
	srcDir      string
	segmentsDir string
	port        string
	jobQueue    chan string
	stop        chan struct{}
}

func (h *HLS) AddVideoFromDir(dir string) error {
	files, err := fs.ListVideoFiles(dir)
	if err != nil {
		return errors.Wrapf(err, "Failed to list video files in %s directory", dir)
	}
	for _, file := range files {
		h.AddVideoFromFile(file)
	}
	return nil
}

func (h *HLS) AddVideoFromFile(videoFile string) {
	h.jobQueue <- videoFile
}

func (h *HLS) Close() {
	close(h.stop)
}

func (h *HLS) addVideoEndpoint(endpoint string) {
	http.HandleFunc(fmt.Sprintf("/%s/", endpoint), h.httpHandler)
	log.Printf("Added video %s enpoint", endpoint)
}

func (h *HLS) watchForNewVideoFiles() {
	c, err := notify.WatchDir(h.segmentsDir)
	if err != nil {
		log.Print("Failed to start dir watch: %s", err)
		return
	}
	go func() {
		select {
		case file := <-c:
			if fs.IsVideoFile(file) {
				log.Print("Dectected new video file %s", file)
				h.jobQueue <- file
			}
		case <-h.stop:
			return
		}
	}()
}

func (h *HLS) startSegmentationJobs(n int) {
	for i := 0; i < n; i++ {
		go func() {
			h.segmentationJob(context.Background())
		}()
	}
}

func (h *HLS) segmentationJob(ctx context.Context) {
	for {
		select {
		case file := <-h.jobQueue:
			if err := video.Segmentation(h.srcDir+"/"+file, h.segmentsDir); err != nil {
				// TODO propage error
				log.Printf("Failed to process %s file", file)
				continue
			}
			h.addVideoEndpoint(fs.FileName(file))
		case <-ctx.Done():
			return
		case <-h.stop:
			return
		}
	}
}

func (h *HLS) StartVideoStreaming() error {
	h.AddVideoFromDir(h.srcDir)

	return http.ListenAndServe(fmt.Sprintf(":%s", h.port), nil)
}

func (h *HLS) httpHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	path := urlToFilePath(r.URL.Path)
	data, err := ioutil.ReadFile(h.segmentsDir + "/" + path)
	if err != nil {
		log.Printf("HttpHandler failed to read file %s", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header()["Access-Control-Allow-Orgin"] = []string{"*"}
	w.Write(data)
}

func isIdxReq(url string) bool {
	tab := strings.Split(url, "/")
	return len(tab) == 2 && tab[1] == ""
}

func urlToFilePath(path string) string {
	a := strings.Split(path, "/")
	if len(a) == 3 && a[2] == "" {
		return fmt.Sprintf("%s/index.m3u8", a[1])
	}
	return path
}
