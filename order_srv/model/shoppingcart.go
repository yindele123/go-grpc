package model

import "project/order_srv/global"

type Shoppingcart struct {
	ID        uint32
	User      uint32 `gorm:"comment:'用户id';default:0"`
	Goods     uint64 `gorm:"comment:'商品id';default:0"`
	Nums      uint32 `gorm:"comment:'购买数量';default:0"`
	Checked   bool   `gorm:"type:bool;comment:'是否选中,1:是 0:否';default:false"`
	CreatedAt uint32 `gorm:"comment:'添加时间';default:0"`
	UpdatedAt uint32 `gorm:"comment:'更新时间';default:0"`
	IsDeleted bool   `gorm:"type:bool;comment:'是否删除,1:是 0:否';default:false"`
	DeletedAt uint32 `gorm:"comment:'删除时间';default:0"`
}

func GetShoppingcartList(whereSql string, vals []interface{}, fields string, Offset int, limit int) (resShoppingcart []Shoppingcart, rows int64, err error) {
	mod := global.MysqlDb.Limit(limit).Offset(Offset)
	if len(fields) != 0 {
		mod.Select(fields)
	}
	if len(whereSql) != 0 && len(vals) != 0 {
		mod.Where(whereSql, vals...)
	}
	result := mod.Debug().Find(&resShoppingcart)
	return resShoppingcart, result.RowsAffected, result.Error
}

func GetShoppingcartCount(whereSql string, vals []interface{}) (count int64, err error) {
	mod := global.MysqlDb.Model(&Shoppingcart{})
	if len(whereSql) != 0 && len(vals) != 0 {
		mod.Where(whereSql, vals...)
	}
	result := mod.Count(&count)
	return count, result.Error
}

func UpdateShoppingcart(data interface{}, whereSql string, vals []interface{}) (err error) {
	if data == nil || len(whereSql) == 0 || len(vals) == 0 {
		return
	}
	result := global.MysqlDb.Model(&Shoppingcart{}).Where(whereSql, vals...).Updates(data)
	return result.Error
}

func CreateShoppingcart(shoppingcart Shoppingcart) (data Shoppingcart, err error) {
	result := global.MysqlDb.Create(&shoppingcart)
	return shoppingcart, result.Error
}

func GetShoppingcartFirst(whereSql string, vals []interface{}, fields string) (shoppingcartFirst Shoppingcart, rows int64, err error) {
	mod := global.MysqlDb.Limit(1)
	if len(fields) != 0 {
		mod.Select(fields)
	}
	if len(whereSql) != 0 && len(vals) != 0 {
		mod.Where(whereSql, vals...)
	}
	result := mod.Find(&shoppingcartFirst)
	return shoppingcartFirst, result.RowsAffected, result.Error
}

