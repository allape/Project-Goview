package model

import "github.com/allape/gocrud"

type Tag struct {
	gocrud.Base
	Name  string `json:"name"`
	Key   string `json:"key" gorm:"uniqueIndex;type:varchar(200)"`
	Color string `json:"color"`
}
