package model

import "project/order_srv/global"

type Ordergoods struct {
	ID         uint64
	Order      uint64  `gorm:"comment:'订单id';default:0"`
	Goods      uint64  `gorm:"comment:'商品id';default:0"`
	OnSale     bool    `gorm:"type:bool;comment:'是否上架  1:是  0:否';default:true"`
	GoodsSn    string  `gorm:"type:varchar(50);default:'';comment:'商品唯一货号'"`
	GoodsName  string  `gorm:"type:varchar(100);default:'';comment:'商品名称'"`
	GoodsImage string  `gorm:"type:varchar(200);default:'';comment:'商品图片'"`
	GoodsPrice float32 `gorm:"type:decimal(10,2);comment:'价格';default:0"`
	Nums       uint32  `gorm:"comment:'商品数量';default:0"`
	CreatedAt  uint32  `gorm:"comment:'添加时间';default:0"`
	UpdatedAt  uint32  `gorm:"comment:'更新时间';default:0"`
	IsDeleted  bool    `gorm:"type:bool;comment:'是否删除,1:是 0:否';default:false"`
	DeletedAt  uint32  `gorm:"comment:'删除时间';default:0"`
}

func GetOrdergoodsList(whereSql string, vals []interface{}, fields string, Offset int, limit int) (resOrdergoods []Ordergoods, rows int64, err error) {
	mod := global.MysqlDb.Limit(limit).Offset(Offset)
	if len(fields) != 0 {
		mod.Select(fields)
	}
	if len(whereSql) != 0 && len(vals) != 0 {
		mod.Where(whereSql, vals...)
	}
	result := mod.Debug().Find(&resOrdergoods)
	return resOrdergoods, result.RowsAffected, result.Error
}
