package notify

import (
	"strings"
	"unsafe"

	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
)

func WatchDir(dir string) (chan string, error) {
	fd, err := unix.InotifyInit()
	if fd == -1 {
		return nil, errors.Wrap(err, "inotify syscall failed")
	}
	wd, err := unix.InotifyAddWatch(fd, dir, unix.IN_ALL_EVENTS)
	if wd == -1 {
		return nil, errors.Wrap(err, "inotify_add_watch syscall failed")
	}
	var buf [(unix.SizeofInotifyEvent + unix.PathMax) * 128]byte
	ch := make(chan string)
	go func() {
		for {
			n, err := unix.Read(fd, buf[:])
			if err == unix.EINTR {
				continue
			}

			if n < unix.SizeofInotifyEvent {
				continue
			}

			var offset uint32
			for offset <= uint32(n) {
				raw := (*unix.InotifyEvent)(unsafe.Pointer(&buf[offset]))

				mask := uint32(raw.Mask)
				len := uint32(raw.Len)

				tmp := (*[unix.PathMax]byte)(unsafe.Pointer(&buf[offset+unix.SizeofInotifyEvent]))
				file := strings.TrimRight(string(tmp[0:len]), "\000")

				offset += unix.SizeofInotifyEvent + len
				if (mask&unix.IN_CLOSE_WRITE != unix.IN_CLOSE_WRITE) && (mask&unix.IN_MODIFY != unix.IN_MODIFY) {
					continue
				}
				ch <- file
			}
		}
	}()
	return ch, nil
}
