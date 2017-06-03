package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	videodir = flag.String("d", "/tmp/video", "Video source dirctory")
	port     = flag.String("p", "8888", "Port number")
	num      = flag.Int("w", 4, "Number of encoder workers")
)

func run() {
	list, err := getVideoListFromDir(*videodir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error ReadDir(%s): %s", *videodir, err.Error())
		os.Exit(-1)
	}

	ch := make(chan string)
	for i := 0; i < *num; i++ {
		go worker(ch)
	}

	for _, s := range list {
		ch <- s
	}

	direvent := make(chan string)
	go watchDirForNewVideo(*videodir, direvent)

	go func() {
		for {
			select {
			case file := <-direvent:
				log.Printf("Detected new video file: %s", file)
				ch <- file
			}
		}
	}()
}

func main() {
	flag.Parse()

	if !dirExists(*videodir) {
		fmt.Fprintf(os.Stderr, "Can't find directory: %s\n", *videodir)
		os.Exit(-1)
	}

	run()
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
