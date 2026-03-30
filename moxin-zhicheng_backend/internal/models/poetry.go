package model

import "gorm.io/gorm"

// Poetry 诗词表模型
type Poetry struct {
	gorm.Model
	// Title 标题，增加索引方便按标题搜索
	Title string `gorm:"type:varchar(255);index;comment:标题"`

	// Author 作者，增加索引方便按作者搜索
	Author string `gorm:"type:varchar(100);index;comment:作者"`

	// Paragraphs 正文内容，使用 text 类型存储长文本
	Paragraphs string `gorm:"type:text;comment:正文内容"`

	// Type 区分类型，比如：tang (唐诗), song (宋词)
	Type string `gorm:"type:varchar(50);index;comment:诗词类型"`

	// Dynasty 朝代，虽然 type 能区分，但有时需要更细的维度
	Dynasty string `gorm:"type:varchar(20);comment:朝代"`

	Chapter string `gorm:"type:varchar(32);comment:章节"`

	// Rhythmic 词牌名
	Rhythmic string `gorm:"type:varchar(32);comment:词牌名"`

	// Translation 译文，存储爬取到的译文数据
	Translation string `gorm:"type:text;comment:译文"`

	// Annotation 注释，存储爬取到的注释数据
	Annotation string `gorm:"type:text;comment:注释"`
}

// TableName 指定表名（可选，不指定默认是 poetries）
func (Poetry) TableName() string {
	return "poetry"
}
