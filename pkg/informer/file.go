package informer

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type File interface {
	Name() string

	Stat() os.FileInfo
	Content() []byte
}

type localFile struct {
	name    string
	baseDir string
	content []byte
	stat    os.FileInfo
}

var (
	ErrIsDirectory = errors.New("is a directory")
)

func NewFile(baseDir, name string) (File, error) {
	absPath := filepath.Join(baseDir, name)
	stat, err := os.Stat(absPath)
	if err == os.ErrNotExist {
		return &localFile{name: name, baseDir: baseDir}, nil
	}
	if err != nil {
		return nil, err
	}
	if stat.IsDir() {
		return nil, ErrIsDirectory
	}
	content, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, err
	}
	return &localFile{
		name:    stat.Name(),
		baseDir: baseDir,
		content: content,
		stat:    stat,
	}, nil
}

func (f *localFile) Name() string {
	return f.name
}

func (f *localFile) Stat() os.FileInfo {
	return f.stat
}

func (f *localFile) Content() []byte {
	return f.content
}
