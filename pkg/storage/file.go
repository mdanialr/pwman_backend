package storage

import (
	"io"
	"os"

	help "github.com/mdanialr/pwman_backend/pkg/helper"

	"go.uber.org/zap"
)

// NewFile return implementation of storage.I that use local file system as the
// storage.
func NewFile(zap *zap.Logger) Port {
	return &fileStorage{zap}
}

type fileStorage struct {
	zap *zap.Logger
}

func (f *fileStorage) Store(rc io.ReadCloser, s string) {
	fl, err := os.Create(s)
	if err != nil {
		f.zap.Error(help.Pad("failed to create file with name", s, ":", err.Error()))
		return
	}
	defer fl.Close()
	defer rc.Close()

	// copy from rc to fl
	if _, err = io.Copy(fl, rc); err != nil {
		f.zap.Error(help.Pad("failed to copy file to", s, ":", err.Error()))
	}
}

func (f *fileStorage) Remove(s string) {
	if err := os.Remove(s); err != nil {
		f.zap.Error(help.Pad("failed to remove", s, ":", err.Error()))
	}
}
