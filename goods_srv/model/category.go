package model

type Category struct {
	ID               uint32
	Name             string `gorm:"type:varchar(20);default:'';comment:'名称'"`
	ParentCategoryId uint32 `gorm:"index:category_parent_category_id;comment:'父类别';default:0"`
	Level            uint8  `gorm:"type:tinyint(1) UNSIGNED;comment:'级别';default:0"`
	IsTab            uint8  `gorm:"type:tinyint(1) UNSIGNED;comment:'是否显示在首页tab,1:是,0:否';default:0"`
	CreatedAt        uint32 `gorm:"comment:'添加时间';default:0"`
	UpdatedAt        uint32 `gorm:"comment:'更新时间';default:0"`
	IsDeleted        uint8  `gorm:"type:tinyint(1) UNSIGNED;comment:'是否删除,1:是 0:否';default:0"`
	DeletedAt        uint32 `gorm:"comment:'删除时间';default:0"`
}

type Goodscategorybrand struct {
	ID         uint32
	CategoryId uint32 `gorm:"index:category_id;index:category_id_brand_id,unique;comment:'类别';default:0"`
	BrandId    uint32 `gorm:"index:brand_id;index:category_id_brand_id,unique;comment:'品牌';default:0"`
	CreatedAt  uint32 `gorm:"comment:'添加时间';default:0"`
	UpdatedAt  uint32 `gorm:"comment:'更新时间';default:0"`
	IsDeleted  uint8  `gorm:"type:tinyint(1) UNSIGNED;comment:'是否删除,1:是 0:否';default:0"`
	DeletedAt  uint32 `gorm:"comment:'删除时间';default:0"`
}
