package file

import "io"

type Repository interface {
	GetData(name string) (io.ReadCloser, error)
	GetInfo(name string) (Info, error)
	GetList() ([]Info, error)
}
