package informer

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

type File interface {
	Name() string

	Stat() os.FileInfo

	Content() []byte
	ContentSum256() string
}

type localFile struct {
	name    string
	content []byte
	stat    os.FileInfo
}

var (
	ErrIsDirectory = errors.New("is a directory")
)

func NewFile(fileName string) (File, error) {
	stat, err := os.Stat(fileName)
	if err != nil {
		return nil, err
	}
	if stat.IsDir() {
		return nil, ErrIsDirectory
	}
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return &localFile{
		name:    fileName,
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

func (f *localFile) ContentSum256() string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(string(f.Content()))))
}
