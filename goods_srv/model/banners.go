package model

type Banners struct {
	ID        uint32
	Image     string `gorm:"type:varchar(200);default:'';comment:'图片url'"`
	Url       string `gorm:"type:varchar(200);default:'';comment:'访问url'"`
	Index     uint32 `gorm:"comment:'轮播顺序';default:0"`
	CreatedAt uint32 `gorm:"comment:'添加时间';default:0"`
	UpdatedAt uint32 `gorm:"comment:'更新时间';default:0"`
	IsDeleted uint8  `gorm:"type:tinyint(1) UNSIGNED;comment:'是否删除,1:是 0:否';default:0"`
	DeletedAt uint32 `gorm:"comment:'删除时间';default:0"`
}
