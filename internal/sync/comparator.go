package sync

import (
	"crypto/sha256"
	"io"
	"os"
	"time"
)

type FileRecord struct {
	Path    string
	ModTime time.Time
	Size    int64
}

func NeedsSync(src, dst FileRecord) bool {
	return src.ModTime.After(dst.ModTime) || src.Size != dst.Size
}

func ChecksumMatch(srcPath, dstPath string) bool {
	s, err := os.Open(srcPath)
	if err != nil {
		return false
	}
	defer s.Close()
	d, err := os.Open(dstPath)
	if err != nil {
		return false
	}
	defer d.Close()
	hs := sha256.New()
	hd := sha256.New()
	_, _ = io.Copy(hs, s)
	_, _ = io.Copy(hd, d)
	return string(hs.Sum(nil)) == string(hd.Sum(nil))
}
