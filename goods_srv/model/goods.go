package model

type Goods struct {
	ID              uint64
	CategoryId      uint32  `gorm:"index:category_id;comment:'商品类目id';default:0"`
	BrandId         uint32  `gorm:"index:brand_id;comment:'品牌类目id';default:0"`
	OnSale          uint8   `gorm:"type:tinyint(2) UNSIGNED;comment:'是否上架  1:是  0:否';default:1"`
	GoodsSn         string  `gorm:"type:varchar(50);default:'';comment:'商品唯一货号'"`
	Name            string  `gorm:"type:varchar(100);default:'';comment:'商品名称'"`
	ClickNum        uint32  `gorm:"comment:'点击数';default:0"`
	SoldNum         uint32  `gorm:"comment:'商品销售量';default:0"`
	FavNum          uint32  `gorm:"comment:'收藏数';default:0"`
	MarketPrice     float32 `gorm:"type:decimal(10,2);comment:'市场价格';default:0"`
	ShopPrice       float32 `gorm:"type:decimal(10,2);comment:'本店价格';default:0"`
	GoodsBrief      string  `gorm:"type:varchar(200);default:'';comment:'商品简短描述'"`
	ShipFree        uint8   `gorm:"type:tinyint(1) UNSIGNED;comment:'是否承担运费,1:是 0:否';default:1"`
	Images          string  `gorm:"type:json;comment:'商品轮播图'"`
	DescImages      string  `gorm:"type:json;comment:'详情页图片'"`
	GoodsFrontImage string  `gorm:"type:varchar(200);default:'';comment:'封面图'"`
	IsNew           uint8   `gorm:"type:tinyint(1) UNSIGNED;comment:'是否新品,1:是 0:否';default:0"`
	IsHot           uint8   `gorm:"type:tinyint(1) UNSIGNED;comment:'是否热销,1:是 0:否';default:0"`
	CreatedAt       uint32  `gorm:"comment:'添加时间';default:0"`
	UpdatedAt       uint32  `gorm:"comment:'更新时间';default:0"`
	IsDeleted       uint8   `gorm:"type:tinyint(1) UNSIGNED;comment:'是否删除,1:是 0:否';default:0"`
	DeletedAt       uint32  `gorm:"comment:'删除时间';default:0"`
}
