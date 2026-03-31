package model

import (
	"gorm.io/gorm"
)

type Article struct {
	gorm.Model
	Title    string `gorm:"type:varchar(200);not null;index" json:"title"` // 标题，加索引方便搜索
	Summary  string `gorm:"type:varchar(500)" json:"summary"`              // 摘要，用于列表展示，避免加载全文
	Content  string `gorm:"type:longtext;not null" json:"content"`         // 富文本内容（1万字以内）
	AuthorID uint   `gorm:"index" json:"author_id"`                        // 作者ID
	//Season   string `gorm:"type:varchar(20)" json:"season"`               // 季节：春/夏/秋/冬
	//Mood     string `gorm:"type:varchar(50)" json:"mood"`                 // 心情标签
	ViewCount int `gorm:"default:0" json:"view_count"` // 阅读量
	IsPublic  int `gorm:"default:0" json:"is_public"`  //0为私有 1为可见
}
