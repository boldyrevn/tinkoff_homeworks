package osfilesys

import (
	"homework/internal/domain/file"
	"io"
	"os"
	"path"
)

// fileRepository implements file.Repository interface using OS file system
type fileRepository struct {
	rootDir string
}

func (fr *fileRepository) GetData(name string) (io.ReadCloser, error) {
	f, err := os.Open(path.Join(fr.rootDir, name))
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fr *fileRepository) GetInfo(name string) (file.Info, error) {
	stat, err := os.Stat(path.Join(fr.rootDir, name))
	if err != nil {
		return file.Info{}, err
	}
	return file.Info{
		Name:    stat.Name(),
		Size:    stat.Size(),
		ModTime: stat.ModTime(),
	}, nil
}

func (fr *fileRepository) GetList() ([]file.Info, error) {
	entries, err := os.ReadDir(fr.rootDir)
	if err != nil {
		return nil, err
	}
	infos := make([]file.Info, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil { // may occur when file is deleted or modified
			continue
		}
		infos = append(infos, file.Info{
			Name:    info.Name(),
			Size:    info.Size(),
			ModTime: info.ModTime(),
		})
	}
	return infos, nil
}

// NewFileRepository creates new file repository, checks if rootDir exists. Otherwise, returns an error
func NewFileRepository(rootDir string) (file.Repository, error) {
	if _, err := os.Stat(rootDir); err != nil {
		return nil, err
	}
	return &fileRepository{rootDir: rootDir}, nil
}
