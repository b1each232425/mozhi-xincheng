package model

type PoetryTagRelation struct {
	ID       uint `gorm:"primary_key"`
	PoetryID uint `gorm:"index:idx_poetry_id"`
	TagID    uint `gorm:"index:idx_tag_id"` // 对应 PoetryTag 表 ID
}
