package app

import (
	"homework/internal/domain/file"
	"io"
)

type UseCase interface {
	GetFileData(name string) (io.ReadCloser, error)
	GetFileInfo(name string) (file.Info, error)
	GetFileList() ([]file.Info, error)
}
