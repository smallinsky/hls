package main

import (
	"flag"
	"log"

	"github.com/smallinsky/hls/pkg/fs"
	"github.com/smallinsky/hls/pkg/hls"
)

var (
	videoSrcDir = flag.String("d", "/tmp/video", "Video source dirctory")
	segmentsDir = flag.String("s", "/tmp/segments", "Video segment dirctory")
	port        = flag.String("p", "8888", "Port number")
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flag.Parse()
	if !fs.DirExist(*videoSrcDir) {
		log.Fatalf("Can't find directory %s\n", *videoSrcDir)
	}

	hlsService := hls.NewService(*videoSrcDir, *segmentsDir, *port)
	defer hlsService.Close()

	log.Print("starting video streaming")
	hlsService.StartVideoStreaming()

}
