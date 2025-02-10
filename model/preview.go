package model

import (
	"crypto/tls"
	"errors"
	"fmt"
	vfs "github.com/allape/go-http-vfs"
	"github.com/allape/gocrud"
	"github.com/allape/goview/env"
	"github.com/allape/goview/util"
	"github.com/h2non/filetype"
	"gorm.io/gorm"
	"image"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

type PreviewableFile interface {
	io.ReaderFrom
	Stat() (os.FileInfo, error)
}

type FileKey string

type Preview struct {
	gocrud.Base
	DatasourceID gocrud.ID `json:"datasourceId"`
	Key          FileKey   `json:"key" gorm:"uniqueIndex"`
	Digest       string    `json:"digest" gorm:"uniqueIndex;type:varchar(64)"`
	Cover        string    `json:"cover"`
	FFProbeInfo  string    `json:"ffprobeInfo"`
}

func BuildPreviewKey(datasource Datasource, file string) FileKey {
	return FileKey(fmt.Sprintf("goview://%d?file=%s", datasource.ID, url.QueryEscape(file)))
}

func GetPreview(repo *gorm.DB, datasource Datasource, file string) (FileKey, *Preview, error) {
	key := BuildPreviewKey(datasource, file)
	var pre Preview
	err := repo.First(&pre, "`key` = ?", key).Error
	if err != nil {
		return key, nil, err
	}
	if pre.ID == 0 {
		return key, nil, errors.New("preview not found")
	}
	return key, &pre, nil
}

var locker = &sync.Mutex{}

func GeneratePreview(datasource Datasource, srcFile, dstFolder string, finder func(digest string) (*Preview, error)) (*Preview, error) {
	key := BuildPreviewKey(datasource, srcFile)

	locker.Lock()
	defer locker.Unlock()

	var file PreviewableFile

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

		dufs, err := vfs.NewDufsVFS(datasource.Cwd)
		if err != nil {
			return nil, err
		}
		dufs.SetHttpClient(client)

		dufsFile, err := dufs.Open(srcFile)
		if err != nil {
			return nil, err
		}

		file = dufsFile.(*vfs.DufsFile)
	case LOCAL:
		f, err := os.Open(srcFile)
		if err != nil {
			return nil, err
		}

		file = f
	default:
		return nil, errors.New("datasource not supported")
	}

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	} else if stat.IsDir() {
		return nil, errors.New("can not preview a directory")
	}

	tmpFile, err := os.CreateTemp(os.TempDir(), "goview_*_"+stat.Name())
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tmpFile.Close()
	}()

	n, err := file.ReadFrom(tmpFile)
	if err != nil {
		return nil, err
	} else if n != stat.Size() {
		return nil, fmt.Errorf("unable to read the whole file, expected %d, got %d", stat.Size(), n)
	}

	digest, err := util.Sha256(tmpFile)
	if err != nil {
		return nil, err
	}

	found, err := finder(digest)
	if err == nil {
		found.ID = 0
		found.CreatedAt = time.Now()
		found.UpdatedAt = time.Now()
		found.DeletedAt = nil
		found.DatasourceID = datasource.ID
		found.Key = key
		return found, nil
	}

	prev := Preview{
		DatasourceID: datasource.ID,
		Key:          key,
		Digest:       digest,
	}

	fileType, err := filetype.MatchFile(tmpFile.Name())
	if err != nil {
		return nil, err
	}
	dstFile := fmt.Sprintf("%s/%s.%s", digest[0:4], digest, "jpg")

	prev.FFProbeInfo, err = util.FFProbeInfo(tmpFile.Name())
	if err != nil {
		return nil, err
	}

	fullDstFilePath := path.Join(dstFolder, dstFile)

	folder := path.Dir(fullDstFilePath)
	err = os.MkdirAll(folder, 0755)
	if err != nil {
		return nil, err
	}

	switch fileType.MIME.Type {
	case "image":
		ext := strings.ToLower(path.Ext(tmpFile.Name()))
		switch ext {
		case ".gif":
			_, err = util.FFMpegVideoSampleImage(tmpFile.Name(), fullDstFilePath, 0.5, image.Point{X: 2, Y: 2})
			if err != nil {
				return nil, err
			}
		case ".raw":
			fallthrough
		case ".arw":
			err = util.ExifToolPreview(fullDstFilePath, tmpFile.Name())
			if err != nil {
				return nil, err
			}
		default:
			_, err = util.FFMpegScaleImage(fullDstFilePath, tmpFile.Name(), 0.1)
			if err != nil {
				return nil, err
			}
		}
	case "video":
		_, err = util.FFMpegVideoSampleImage(tmpFile.Name(), fullDstFilePath, 0.25, image.Point{X: 10, Y: 10})
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("filetype %s is not supported", fileType.MIME.Type)
	}

	prev.Cover = dstFile

	return &prev, nil
}
