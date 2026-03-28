package model

import "gorm.io/gorm"

type PoetryTag struct {
	gorm.Model
	Tag string `gorm:"type:varchar(255);uniqueIndex;not null;comment:关键词"`
	// Count 记录该标签被提取到的总次数
	Count int64 `gorm:"type:int;default:1;comment:提取次数"`
}

func (PoetryTag) TableName() string {
	return "poetry_tag"
}
