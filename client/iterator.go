package client

import (
	"time"

	fileheap "github.com/beaker/fileheap/client"
	"github.com/pkg/errors"
)

// FileIterator is an iterator over files within a dataset.
type FileIterator interface {
	Next() (*FileHandle, *FileInfo, error)
}

// FileInfo describes a single file within a dataset.
type FileInfo struct {
	// Path of the file relative to its dataset root.
	Path string `json:"path"`

	// Size of the file in bytes.
	Size int64 `json:"size"`

	// Time at which the file was last updated.
	Updated time.Time `json:"updated"`
}

// ErrDone indicates an iterator is expended.
var ErrDone = errors.New("no more items in iterator")

// fileHeapIterator is an iterator over files within a FileHeap dataset.
type fileHeapIterator struct {
	dataset  *DatasetHandle
	iterator *fileheap.FileIterator
}

func (i *fileHeapIterator) Next() (*FileHandle, *FileInfo, error) {
	info, err := i.iterator.Next()
	if err == fileheap.ErrDone {
		return nil, nil, ErrDone
	}
	if err != nil {
		return nil, nil, err
	}
	return &FileHandle{
			dataset: i.dataset,
			file:    info.Path,
		}, &FileInfo{
			Path:    info.Path,
			Size:    info.Size,
			Updated: info.Updated,
		}, nil
}
