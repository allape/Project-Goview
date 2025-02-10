package model

import "github.com/allape/gocrud"

type Tag struct {
	gocrud.Base
	Name  string `json:"name"`
	Color string `json:"color"`
}
