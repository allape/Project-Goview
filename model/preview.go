package model

import (
	"errors"
	"fmt"
	"image"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/allape/gocrud"
	"github.com/allape/gogger"
	"github.com/allape/goview/util"
	"github.com/h2non/filetype"
)

var l = gogger.New("model")

type FileKey string

type Preview struct {
	gocrud.Base
	DatasourceID gocrud.ID `json:"datasourceId"`
	Key          FileKey   `json:"key"`
	Digest       string    `json:"digest" gorm:"type:varchar(64)"`
	Cover        string    `json:"cover"`
	MIME         string    `json:"mime"`
	FFProbeInfo  string    `json:"ffprobeInfo"`
}

func BuildPreviewKey(datasource Datasource, file string) FileKey {
	return FileKey(fmt.Sprintf("goview://%d%s", datasource.ID, file))
}

var locker = &sync.Mutex{}

func GeneratePreview(datasource Datasource, srcFile, dstFolder string, finder func(digest string) (*Preview, error)) (*Preview, error) {
	key := BuildPreviewKey(datasource, srcFile)

	locker.Lock()
	defer locker.Unlock()

	dfs, err := GetFS(datasource)
	if err != nil {
		return nil, err
	}

	file, err := dfs.Open(srcFile)
	if err != nil {
		return nil, err
	}

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	} else if stat.IsDir() {
		return nil, errors.New("can not preview a directory")
	}

	tmpFile, err := os.CreateTemp(os.TempDir(), fmt.Sprintf("goview_*%s", path.Ext(stat.Name())))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tmpFile.Close()
		_ = os.Remove(tmpFile.Name())
	}()

	n, err := file.WriteTo(tmpFile)
	if err != nil {
		return nil, err
	} else if n != stat.Size() {
		return nil, fmt.Errorf("unable to read the whole file, expected %d, got %d", stat.Size(), n)
	}

	digest, err := util.Sha256(tmpFile)
	if err != nil {
		return nil, err
	}

	l.Info().Printf("digest of %s = %s", srcFile, digest)

	found, err := finder(digest)
	if err == nil {
		l.Info().Printf("found preview %s", found.Key)
		found.ID = 0
		found.CreatedAt = time.Now()
		found.UpdatedAt = time.Now()
		found.DeletedAt = nil
		found.DatasourceID = datasource.ID
		found.Key = key
		return found, nil
	}

	l.Info().Printf("generating preview for %s", key)

	fileType, err := filetype.MatchFile(tmpFile.Name())
	if err != nil {
		return nil, err
	}

	prev := Preview{
		DatasourceID: datasource.ID,
		Key:          key,
		MIME:         fileType.MIME.Value,
		Digest:       digest,
	}
	dstFile := fmt.Sprintf("%s/%s.%s", digest[0:4], digest, "jpg")

	prev.FFProbeInfo, err = util.FFProbeInfo(tmpFile.Name())
	if err != nil {
		return nil, err
	}

	fullDstFilePath := path.Join(dstFolder, dstFile)

	err = os.MkdirAll(path.Dir(fullDstFilePath), 0755)
	if err != nil {
		return nil, err
	}

	coverStat, err := os.Stat(fullDstFilePath)
	if err == nil && coverStat.Size() > 0 {
		l.Info().Printf("cover %s already exists", fullDstFilePath)
		prev.Cover = dstFile
	} else {
		l.Info().Printf("generating cover %s", fullDstFilePath)
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
	}

	return &prev, nil
}
