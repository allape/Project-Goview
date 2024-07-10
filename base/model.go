package base

import (
	"gorm.io/gorm"
	"time"
)

type Model struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time      `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"autoUpdateTime:milli"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

type R[T any] struct {
	Code    string `json:"c"`
	Message string `json:"m"`
	Data    T      `json:"d"`
}
