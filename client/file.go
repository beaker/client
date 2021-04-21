package client

import (
	"context"
	"io"
	"os"
	"path/filepath"
)

// FileHandle provides operations on a file within a dataset.
type FileHandle struct {
	dataset *DatasetHandle
	file    string
}

// FileRef creates an actor for an existing file within a dataset.
// This call doesn't perform any network operations.
func (h *DatasetHandle) FileRef(filePath string) *FileHandle {
	return &FileHandle{h, filePath}
}

// Download gets a file from a datastore.
func (h *FileHandle) Download(ctx context.Context) (io.ReadCloser, error) {
	return h.dataset.Storage.ReadFile(ctx, h.file)
}

// DownloadRange reads a range of bytes from a file.
// If length is negative, the file is read until the end.
func (h *FileHandle) DownloadRange(ctx context.Context, offset, length int64) (io.ReadCloser, error) {
	return h.dataset.Storage.ReadFileRange(ctx, h.file, offset, length)
}

// DownloadTo downloads a file and writes it to disk.
func (h *FileHandle) DownloadTo(ctx context.Context, filePath string) error {
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}

	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer safeClose(f)

	var written int64
	for {
		var r io.ReadCloser
		var err error
		if written == 0 {
			r, err = h.Download(ctx)
		} else {
			r, err = h.DownloadRange(ctx, written, -1)
		}
		if err != nil {
			return err
		}

		n, err := io.Copy(f, r)
		safeClose(r)
		if err == nil {
			return nil
		}
		written += n
	}
}

// Upload creates or overwrites a file.
func (h *FileHandle) Upload(ctx context.Context, source io.Reader, length int64) error {
	return h.dataset.Storage.WriteFile(ctx, h.file, source, length)
}

// Delete removes a file from an uncommitted datastore.
func (h *FileHandle) Delete(ctx context.Context) error {
	return h.dataset.Storage.DeleteFile(ctx, h.file)
}
