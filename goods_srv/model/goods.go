package model

import (
	"project/goods_srv/global"
)

type Goods struct {
	ID              uint64
	CategoryId      uint32   `gorm:"index:category_id;comment:'商品类目id';default:0"`
	BrandId         uint32   `gorm:"index:brand_id;comment:'品牌类目id';default:0"`
	OnSale          bool     `gorm:"type:bool;comment:'是否上架  1:是  0:否';default:true"`
	GoodsSn         string   `gorm:"type:varchar(50);default:'';comment:'商品唯一货号'"`
	Name            string   `gorm:"type:varchar(100);default:'';comment:'商品名称'"`
	ClickNum        uint32   `gorm:"comment:'点击数';default:0"`
	SoldNum         uint32   `gorm:"comment:'商品销售量';default:0"`
	FavNum          uint32   `gorm:"comment:'收藏数';default:0"`
	MarketPrice     float32  `gorm:"type:decimal(10,2);comment:'市场价格';default:0"`
	ShopPrice       float32  `gorm:"type:decimal(10,2);comment:'本店价格';default:0"`
	GoodsBrief      string   `gorm:"type:varchar(200);default:'';comment:'商品简短描述'"`
	ShipFree        bool     `gorm:"type:bool;comment:'是否承担运费,1:是 0:否';default:true"`
	Images          []string `gorm:"type:json;comment:'商品轮播图'"`
	DescImages      []string `gorm:"type:json;comment:'详情页图片'"`
	GoodsFrontImage string   `gorm:"type:varchar(200);default:'';comment:'封面图'"`
	IsNew           bool     `gorm:"type:bool;comment:'是否新品,1:是 0:否';default:false"`
	IsHot           bool     `gorm:"type:bool;comment:'是否热销,1:是 0:否';default:false"`
	CreatedAt       uint32   `gorm:"comment:'添加时间';default:0"`
	UpdatedAt       uint32   `gorm:"comment:'更新时间';default:0"`
	IsDeleted       bool     `gorm:"type:bool;comment:'是否删除,1:是 0:否';default:false"`
	DeletedAt       uint32   `gorm:"comment:'删除时间';default:0"`
}

func GetGoodsList(whereSql string,vals []interface{}, fields string, Offset int, limit int) (resGoods []Goods, rows int64, err error) {
	mod := global.MysqlDb.Limit(limit).Offset(Offset)
	if len(fields) != 0 {
		mod.Select(fields)
	}
	if len(whereSql) != 0 && len(vals) != 0{
		mod.Where(whereSql,vals...)
	}
	result := mod.Debug().Find(&resGoods)
	return resGoods, result.RowsAffected, result.Error
}

func GetGoodsCount(whereSql string,vals []interface{}) (count int64, err error) {
	mod := global.MysqlDb.Model(&Goods{})
	if len(whereSql) != 0 && len(vals) != 0{
		mod.Where(whereSql,vals...)
	}
	result := mod.Count(&count)
	return count, result.Error
}

func UpdateGoods(data interface{}, whereSql string,vals []interface{}) (err error) {
	if data == nil || len(whereSql) == 0 || len(vals) == 0{
		return
	}
	result := global.MysqlDb.Model(&Goods{}).Where(whereSql,vals...).Updates(data)
	return result.Error
}

func CreateGoods(goods Goods) (data Goods, err error) {
	result := global.MysqlDb.Create(&goods)
	return goods, result.Error
}

func GetGoodsFirst(whereSql string,vals []interface{}, fields string) (goodsFirst Goods, rows int64, err error) {
	mod := global.MysqlDb.Limit(1)
	if len(fields) != 0 {
		mod.Select(fields)
	}
	if len(whereSql) != 0 && len(vals) != 0{
		mod.Where(whereSql,vals...)
	}
	result := mod.Find(&goodsFirst)
	return goodsFirst, result.RowsAffected, result.Error
}
