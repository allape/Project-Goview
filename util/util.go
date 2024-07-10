package util

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"strings"
)

func Sha256(reader io.ReadSeeker) (string, error) {
	_, err := reader.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}

	hasher := sha256.New()
	_, err = io.Copy(hasher, reader)
	if err != nil {
		return "", err
	}

	return strings.ToUpper(hex.EncodeToString(hasher.Sum(nil))), nil
}
