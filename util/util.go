package util

import (
	"io"
	"os"
	"path"
	"strings"
)

// CopyToDir copies the src file to dstdir.
func CopyToDir(src, dstdir string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	_, filename := path.Split(src)
	dst := strings.Join([]string{dstdir, filename}, "/")
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return nil
}
