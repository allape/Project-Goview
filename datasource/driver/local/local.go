package local

import (
	"github.com/allape/goview/datasource/driver"
	"io"
	"os"
	"path"
)

type Driver struct {
	driver.Driver
}

func (d *Driver) PathJoin(segments ...string) string {
	return path.Join(segments...)
}

func (d *Driver) Status(file string) (*driver.File, error) {
	stat, err := os.Stat(file)
	if err != nil {
		return nil, err
	}
	return &driver.File{
		Name:  stat.Name(),
		Size:  stat.Size(),
		IsDir: stat.IsDir(),
	}, nil
}

func (d *Driver) List(wd string) ([]driver.File, error) {
	files, err := os.ReadDir(wd)
	if err != nil {
		return nil, err
	}

	result := make([]driver.File, len(files))

	for index, file := range files {
		stat, err := file.Info()
		if err != nil {
			return nil, err
		}
		result[index] = driver.File{
			Name:  stat.Name(),
			Size:  stat.Size(),
			IsDir: stat.IsDir(),
		}
	}

	return result, nil
}

func (d *Driver) Concatenate(file string, writer io.Writer) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}

	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	_, err = io.Copy(writer, f)

	return err
}
