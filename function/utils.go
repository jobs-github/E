package function

import (
	"errors"
	"fmt"
	"os"
	"unsafe"
)

func FileSize(file string) (int64, error) {
	f, e := os.Stat(file)
	if e != nil {
		return 0, e
	}
	return f.Size(), nil
}

func LoadFile(path string) ([]byte, error) {
	filesize, err := FileSize(path)
	if nil != err {
		return nil, err
	}
	if filesize < 1 {
		return nil, errors.New("empty file")
	}
	buf := make([]byte, filesize)
	f, err := os.Open(path)
	if nil != err {
		return nil, err
	}
	defer f.Close()

	n, err := f.Read(buf)
	if nil != err {
		return nil, err
	}
	if filesize != int64(n) {
		return nil, fmt.Errorf("filesize mismatch, %v, but read %v", filesize, n)
	}
	return buf, nil
}

func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
