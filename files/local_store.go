package files

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

const Kilobyte = 1024
const Megabyte = 1048576

type localStore struct {
	path        string
	maxFileSize int
}

func NewLocalStore(path string, maxFileSize int) (*localStore, error) {
	fPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	return &localStore{
		path:        fPath,
		maxFileSize: maxFileSize,
	}, nil
}

func (l *localStore) Save(path string, content io.Reader) error {
	fp := l.FilePath(path)
	dir := filepath.Dir(fp)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return errors.New("unable to create directory")
	}

	_, err = os.Stat(fp)
	if err == nil {
		err = os.Remove(fp)
		if err != nil {
			return errors.New("unable to delete file")
		}
	} else if !os.IsNotExist(err) {
		return errors.New("unable to get file info")
	}

	f, err := os.Create(fp)
	if err != nil {
		return errors.New("unable to create file")
	}
	defer f.Close()

	_, err = io.Copy(f, content)
	if err != nil {
		return errors.New("unable to write to file")
	}

	return nil
}

func (l *localStore) Get(path string) (*os.File, error) {
	fp := l.FilePath(path)

	f, err := os.Open(fp)
	if err != nil {
		return nil, errors.New("unable to open file")
	}

	return f, nil
}

// FilePath returns the path of the specific file
func (l *localStore) FilePath(path string) string {
	return filepath.Join(l.path, path)
}
