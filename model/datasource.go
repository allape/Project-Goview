package model

import (
	"crypto/tls"
	"errors"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"

	"github.com/allape/gocrud"
	"github.com/allape/gohtvfs"
	"github.com/allape/goview/env"
)

type LocalFS struct {
	DatasourceFS
	wd string
}

func (f *LocalFS) Open(name string) (File, error) {
	return os.Open(path.Join(f.wd, name))
}

func (f *LocalFS) ReadDir(name string) ([]fs.DirEntry, error) {
	return os.ReadDir(path.Join(f.wd, name))
}

type DuFS struct {
	DatasourceFS
	dufs *gohtvfs.DufsVFS
}

func (f *DuFS) Open(name string) (File, error) {
	file, err := f.dufs.Open(name)
	if err != nil {
		return nil, err
	}
	return file.(*gohtvfs.DufsFile), nil
}

func (f *DuFS) ReadDir(name string) ([]fs.DirEntry, error) {
	return f.dufs.ReadDir(name)
}

type File interface {
	io.WriterTo
	Stat() (os.FileInfo, error)
}

type DatasourceFS interface {
	Open(name string) (File, error)
	ReadDir(name string) ([]fs.DirEntry, error)
}

type Type string

const (
	DUFS  Type = "dufs"
	LOCAL Type = "local"
)

type Datasource struct {
	gocrud.Base
	Name string `json:"name"`
	Type Type   `json:"type"`
	Cwd  string `json:"cwd"`
}

func GetFS(datasource Datasource) (DatasourceFS, error) {
	switch datasource.Type {
	case DUFS:
		caCertPool, err := env.TrustedCertsPoolFromEnv()
		if err != nil {
			return nil, err
		}
		tlsConfig := &tls.Config{
			RootCAs: caCertPool,
		}
		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsConfig,
			},
		}

		dufs, err := gohtvfs.NewDufsVFS(datasource.Cwd)
		if err != nil {
			return nil, err
		}
		dufs.SetHttpClient(client)

		return &DuFS{dufs: dufs}, nil
	case LOCAL:
		return &LocalFS{wd: datasource.Cwd}, nil
	default:
		return nil, errors.New("datasource not supported")
	}
}
