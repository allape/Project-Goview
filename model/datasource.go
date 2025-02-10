package model

import (
	"github.com/allape/gocrud"
)

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
