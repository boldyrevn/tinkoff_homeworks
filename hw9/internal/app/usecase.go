package app

import (
	"homework/internal/domain/file"
	"io"
)

type useCase struct {
	fileRepo file.Repository
}

func (u *useCase) GetFileData(name string) (io.ReadCloser, error) {
	return u.fileRepo.GetData(name)
}

func (u *useCase) GetFileInfo(name string) (file.Info, error) {
	return u.fileRepo.GetInfo(name)
}

func (u *useCase) GetFileList() ([]file.Info, error) {
	return u.fileRepo.GetList()
}

func NewUseCase(repo file.Repository) UseCase {
	return &useCase{fileRepo: repo}
}
