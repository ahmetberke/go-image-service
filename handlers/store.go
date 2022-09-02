package handlers

import (
	"io"
	"os"
)

type store interface {
	Save(path string, content io.Reader) error
	Get(path string) (*os.File, error)
	FilePath(path string) string
}
