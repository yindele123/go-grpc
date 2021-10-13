package model

type Brands struct {
	ID        uint32
	Name      string `gorm:"uniqueIndex:brands_name;type:varchar(50);default:'';comment:'名称'"`
	Logo      string `gorm:"type:varchar(200);default:'';comment:'图标'"`
	CreatedAt uint32 `gorm:"comment:'添加时间';default:0"`
	UpdatedAt uint32 `gorm:"comment:'更新时间';default:0"`
	IsDeleted uint8  `gorm:"type:tinyint(1) UNSIGNED;comment:'是否删除,1:是 0:否';default:0"`
	DeletedAt uint32 `gorm:"comment:'删除时间';default:0"`
}
