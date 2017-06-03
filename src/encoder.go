package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

func execute(command string, args []string) error {
	cmd := exec.Command(command, args...)
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("%s", err.Error())
	}

	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("%s", err.Error())
	}
	return nil
}

func encodeArgs(videofile string) []string {
	dir := strings.TrimSuffix(videofile, path.Ext(videofile))
	return []string{
		"-i", fmt.Sprintf("%s/%s", *videodir, videofile),
		"-codec", "copy",
		"-f", "segment",
		"-vbsf", "h264_mp4toannexb",
		"-segment_time", "10",
		"-segment_format", "mpegts",
		"-segment_list", fmt.Sprintf("%s/index.m3u8", dir),
		"-segment_list_type", "m3u8",
		fmt.Sprintf("%s/%s%%d.ts", dir, dir)}
}

func worker(ch chan string) {
	for {
		select {
		case file := <-ch:
			dirname := strings.TrimSuffix(file, path.Ext(file))
			if !dirExists(dirname) {
				err := os.Mkdir(dirname, 0777)
				if err != nil {
					continue
				}

				args := encodeArgs(file)
				log.Printf("Encoding %s started", file)
				err = execute("ffmpeg", args)
				if err != nil {
					log.Println("Faile to encode video file: ", file)
					os.RemoveAll(dirname)
					continue
				}
			}
			addHandler(dirname)
		}
	}
}
