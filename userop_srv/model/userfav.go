package model

type Userfav struct {
	ID        uint32
	User      uint32 `gorm:"comment:'用户id';default:0"`
	Goods     uint64 `gorm:"comment:'商品id';default:0"`
	CreatedAt uint32 `gorm:"comment:'添加时间';default:0"`
	UpdatedAt uint32 `gorm:"comment:'更新时间';default:0"`
	IsDeleted bool   `gorm:"type:bool;comment:'是否删除,1:是 0:否';default:false"`
	DeletedAt uint32 `gorm:"comment:'删除时间';default:0"`
}
