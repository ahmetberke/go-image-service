package files

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func setup(t *testing.T) (*localStore, string, func()) {
	dir, err := ioutil.TempDir("./", "files")
	if err != nil {
		t.Fatal(err)
	}

	l, err := NewLocalStore(dir, 10*Megabyte)
	if err != nil {
		t.Fatal(err)
	}

	return l, dir, func() {
		err := os.RemoveAll(dir)
		if err != nil {
			fmt.Printf("error : %s", err.Error())
		}
	}

}

func TestLocalStore_Save(t *testing.T) {
	savePath := "/1/test.png"
	fileContents := "Hello World"
	l, dir, cf := setup(t)
	defer cf()

	err := l.Save(savePath, bytes.NewBuffer([]byte(fileContents)))
	assert.NoError(t, err)

	f, err := os.Open(filepath.Join(dir, savePath))
	assert.NoError(t, err)
	defer f.Close()

	d, err := ioutil.ReadAll(f)
	assert.NoError(t, err)
	assert.Equal(t, fileContents, string(d))

}

func TestLocalStore_Get(t *testing.T) {
	savePath := "/1/test.png"
	fileContents := "Hello World"
	l, _, cf := setup(t)
	defer cf()

	err := l.Save(savePath, bytes.NewBuffer([]byte(fileContents)))
	assert.NoError(t, err)

	r, err := l.Get(savePath)
	assert.NoError(t, err)
	defer r.Close()

	d, err := ioutil.ReadAll(r)
	assert.NoError(t, err)
	assert.Equal(t, string(d), fileContents)
}
