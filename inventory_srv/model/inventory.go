package model

import (
	"project/inventory_srv/global"
)

type Inventory struct {
	ID        uint64
	Goods     uint64 `gorm:"index:goods_id;index:unique_goods_id,unique;comment:'商品id';default:0"`
	Stocks    uint32 `gorm:"comment:'库存数量';default:0"`
	Version   uint32 `gorm:"comment:'版本号';default:0"`
	CreatedAt uint32 `gorm:"comment:'添加时间';default:0"`
	UpdatedAt uint32 `gorm:"comment:'更新时间';default:0"`
	DeletedAt uint32 `gorm:"comment:'删除时间';default:0"`
}

type Inventoryhistory struct {
	ID             uint64
	OrderSn        string `gorm:"index:order_sn;index:unique_order_sn,unique;type:varchar(50);default:'';comment:'订单编号'"`
	OrderInvDetail string `gorm:"type:varchar(200);default:'';comment:'订单详情'"`
	Status         uint32 `gorm:"type:tinyint(2) UNSIGNED;comment:'出库状态  1:已扣减  2:已归还';default:1"`
	CreatedAt      uint32 `gorm:"comment:'添加时间';default:0"`
	UpdatedAt      uint32 `gorm:"comment:'更新时间';default:0"`
	DeletedAt      uint32 `gorm:"comment:'删除时间';default:0"`
}

func GetInventoryFirst(whereSql string, vals []interface{}, fields string) (InventoryFirst Inventory, rows int64, err error) {
	mod := global.MysqlDb.Limit(1)
	if len(fields) != 0 {
		mod.Select(fields)
	}
	if len(whereSql) != 0 && len(vals) != 0 {
		mod.Where(whereSql, vals...)
	}
	result := mod.Find(&InventoryFirst)
	return InventoryFirst, result.RowsAffected, result.Error
}

func UpdateInventory(data interface{}, whereSql string, vals []interface{}) (err error) {
	if data == nil || len(whereSql) == 0 || len(vals) == 0 {
		return
	}

	result := global.MysqlDb.Model(&Inventory{}).Where(whereSql, vals...).Updates(data)
	return result.Error
}

func CreateInventory(inventory Inventory) (data Inventory, err error) {
	result := global.MysqlDb.Create(&inventory)
	return inventory, result.Error
}

func GetInventoryList(whereSql string, vals []interface{}, fields string, Offset int, limit int) (inventory []Inventory, rows int64, err error) {
	mod := global.MysqlDb.Limit(limit).Offset(Offset)
	if len(fields) != 0 {
		mod.Select(fields)
	}
	if len(whereSql) != 0 && len(vals) != 0 {
		mod.Where(whereSql, vals...)
	}
	result := mod.Debug().Find(&inventory)
	return inventory, result.RowsAffected, result.Error
}
