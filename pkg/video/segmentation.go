package video

import (
	"fmt"
	"os/exec"

	"github.com/pkg/errors"

	"github.com/smallinsky/hls/pkg/fs"
)

func Segmentation(srcFile, dstDir string) error {
	if !fs.DirExist(dstDir) {
		if err := fs.Mk(dstDir); err != nil {
			return errors.Wrapf(err, "failed to create dir '%s'", dstDir)
		}
	}
	outDir := fmt.Sprintf("%s/%s", dstDir, fs.FileName(srcFile))
	if !fs.FileExist(srcFile) {
		return errors.Errorf("failed to find src video file '%s'", srcFile)
	}

	if err := fs.Mk(outDir); err != nil {
		return errors.Wrapf(err, "failed to creater dir '%s'", outDir)
	}

	if err := execute("ffmpeg", encodeArgs(srcFile, dstDir)); err != nil {
		return errors.Wrap(err, "ffmpeg encoding failed")
	}
	return nil
}

func execute(command string, args []string) error {
	cmd := exec.Command(command, args...)
	if err := cmd.Start(); err != nil {
		return errors.Wrap(err, "failed to start command")
	}
	if err := cmd.Wait(); err != nil {
		return errors.Wrap(err, "cmd.Wait() call failed")
	}
	return nil
}

func encodeArgs(srcFile, dstDir string) []string {
	prefix := fs.FileName(srcFile)
	return []string{
		"-i", fmt.Sprintf("%s", srcFile),
		"-codec", "copy",
		"-f", "segment",
		"-vbsf", "h264_mp4toannexb",
		"-segment_time", "10",
		"-segment_format", "mpegts",
		"-segment_list", fmt.Sprintf("%s/%s/index.m3u8", dstDir, prefix),
		"-segment_list_type", "m3u8",
		segmentFormat(dstDir, prefix),
	}
}

func segmentFormat(dir, prefix string) string {
	return fmt.Sprintf("%s/%s/%%d.ts", dir, prefix)
}
