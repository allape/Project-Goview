package dufs

import (
	"encoding/json"
	"github.com/allape/goview/datasource/driver"
	"github.com/allape/goview/datasource/driver/httpbased"
	"io"
	"net/url"
	"path"
	"strconv"
	"strings"
)

type PathType string

const (
	TypeDir PathType = "Dir"
	//TypeFile PathType = "File"
)

type File struct {
	Mtime    int64    `json:"mtime"`
	Name     string   `json:"name"`
	PathType PathType `json:"path_type"`
	Size     *int64   `json:"size"`
}

type Data struct {
	AllowArchive bool    `json:"allow_archive"`
	AllowDelete  bool    `json:"allow_delete"`
	AllowSearch  bool    `json:"allow_search"`
	AllowUpload  bool    `json:"allow_upload"`
	Auth         bool    `json:"auth"`
	DirExists    bool    `json:"dir_exists"`
	Href         string  `json:"href"`
	Kind         string  `json:"kind"`
	Paths        []File  `json:"paths"`
	UriPrefix    string  `json:"uri_prefix"`
	User         *string `json:"user"`
}

type Driver struct {
	httpbased.AbstractDriver
}

func (d *Driver) PathJoin(segments ...string) string {
	length := len(segments)
	refinedSegments := make([]string, length)
	for i, segment := range segments {
		if i == 0 {
			refinedSegments[i] = strings.TrimSuffix(segment, "/")
		} else if i == length-1 {
			refinedSegments[i] = strings.TrimPrefix(segment, "/")
		} else {
			refinedSegments[i] = strings.Trim(segment, "/")
		}
	}
	return strings.Join(refinedSegments, "/")
}

func (d *Driver) Status(file string) (*driver.File, error) {
	_, header, err := d.AbstractDriver.NewRequest(file, httpbased.HeaderOnlyOption{})
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(file)
	if err != nil {
		return nil, err
	}

	f := driver.File{
		Name: path.Base(u.Path),
	}

	if header.Get("etag") != "" {
		f.IsDir = false

		contentLength := header.Get("content-length")
		if contentLength != "" {
			f.Size, err = strconv.ParseInt(contentLength, 10, 64)
			if err == nil {
				return &f, nil
			}
		}
	} else {
		f.IsDir = true
	}

	return &f, nil
}

func (d *Driver) List(wd string) ([]driver.File, error) {
	content, _, err := d.AbstractDriver.NewRequest(wd + "?json")
	if err != nil {
		return nil, err
	}
	var data Data
	err = json.Unmarshal(content, &data)

	var files []driver.File
	for _, file := range data.Paths {
		size := int64(0)
		if file.Size != nil {
			size = *file.Size
		}
		files = append(files, driver.File{
			IsDir: file.PathType == TypeDir,
			Name:  file.Name,
			Size:  size,
		})
	}

	return files, nil
}

func (d *Driver) Concatenate(file string, writer io.Writer) error {
	_, _, err := d.AbstractDriver.NewRequest(file, httpbased.RedirectWriterOption{Writer: writer})
	return err
}
