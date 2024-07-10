package driver

import (
	"io"
)

type File struct {
	IsDir bool   `json:"isDir"`
	Name  string `json:"name"`
	Size  int64  `json:"size"`
}

type Driver interface {
	PathJoin(segments ...string) string
	Status(file string) (*File, error)
	List(wd string) ([]File, error)
	Concatenate(file string, writer io.Writer) error
}
