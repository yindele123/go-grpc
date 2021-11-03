package model

import "project/order_srv/global"

type Orderinfo struct {
	ID           uint64
	User         uint32  `gorm:"comment:'用户id';default:0"`
	OrderSn      string  `gorm:"unique;type:varchar(30);default:'';comment:'订单号'"`
	PayType      uint32  `gorm:"type:tinyint(2) UNSIGNED;comment:'支付方式  1:支付宝  2:微信...';default:1"`
	Status       uint32  `gorm:"type:tinyint(2) UNSIGNED;comment:'订单状态  1:成功  2:超时关闭 3:交易创建 4:交易结束';default:1"`
	TradeNo      string  `gorm:"unique;type:varchar(100);default:'';comment:'交易号'"`
	OrderMount   float32 `gorm:"type:decimal(10,2);comment:'订单金额';default:0"`
	PayTime      uint32  `gorm:"comment:'支付时间';default:0"`
	Address      string  `gorm:"type:varchar(100);default:'';comment:'收货地址'"`
	SignerName   string  `gorm:"type:varchar(20);default:'';comment:'签收人'"`
	SingerMobile string  `gorm:"type:varchar(11);default:'';comment:'联系电话'"`
	Post         string  `gorm:"type:varchar(200);default:'';comment:'留言'"`
	CreatedAt    uint32  `gorm:"comment:'添加时间';default:0"`
	UpdatedAt    uint32  `gorm:"comment:'更新时间';default:0"`
	IsDeleted    bool    `gorm:"type:bool;comment:'是否删除,1:是 0:否';default:false"`
	DeletedAt    uint32  `gorm:"comment:'删除时间';default:0"`
}

func GetOrderinfoList(whereSql string, vals []interface{}, fields string, Offset int, limit int) (resOrderinfo []Orderinfo, rows int64, err error) {
	mod := global.MysqlDb.Limit(limit).Offset(Offset)
	if len(fields) != 0 {
		mod.Select(fields)
	}
	if len(whereSql) != 0 && len(vals) != 0 {
		mod.Where(whereSql, vals...)
	}
	result := mod.Debug().Find(&resOrderinfo)
	return resOrderinfo, result.RowsAffected, result.Error
}

func GetOrderinfoCount(whereSql string, vals []interface{}) (count int64, err error) {
	mod := global.MysqlDb.Model(&Orderinfo{})
	if len(whereSql) != 0 && len(vals) != 0 {
		mod.Where(whereSql, vals...)
	}
	result := mod.Count(&count)
	return count, result.Error
}

func GetOrderinfoFirst(whereSql string, vals []interface{}, fields string) (orderinfoFirst Orderinfo, rows int64, err error) {
	mod := global.MysqlDb.Limit(1)
	if len(fields) != 0 {
		mod.Select(fields)
	}
	if len(whereSql) != 0 && len(vals) != 0 {
		mod.Where(whereSql, vals...)
	}
	result := mod.Find(&orderinfoFirst)
	return orderinfoFirst, result.RowsAffected, result.Error
}

func UpdateOrderinfo(data interface{}, whereSql string, vals []interface{}) (err error) {
	if data == nil || len(whereSql) == 0 || len(vals) == 0 {
		return
	}
	result := global.MysqlDb.Model(&Orderinfo{}).Where(whereSql, vals...).Updates(data)
	return result.Error
}
