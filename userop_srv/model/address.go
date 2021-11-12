package model

import "project/userop_srv/global"

type Address struct {
	ID           uint32
	User         uint32 `gorm:"comment:'用户id';default:0"`
	Province     string `gorm:"type:varchar(100);default:'';comment:'省份'"`
	City         string `gorm:"type:varchar(100);default:'';comment:'城市'"`
	District     string `gorm:"type:varchar(100);default:'';comment:'区域'"`
	Address      string `gorm:"type:varchar(100);default:'';comment:'详细地址'"`
	SignerName   string `gorm:"type:varchar(100);default:'';comment:'签收人'"`
	SignerMobile string `gorm:"type:varchar(11);default:'';comment:'电话'"`
	CreatedAt    uint32 `gorm:"comment:'添加时间';default:0"`
	UpdatedAt    uint32 `gorm:"comment:'更新时间';default:0"`
	IsDeleted    bool   `gorm:"type:bool;comment:'是否删除,1:是 0:否';default:false"`
	DeletedAt    uint32 `gorm:"comment:'删除时间';default:0"`
}

func GetAddressList(whereSql string, vals []interface{}, fields string, Offset int, limit int) (resAddress []Address, rows int64, err error) {
	mod := global.MysqlDb.Limit(limit).Offset(Offset)
	if len(fields) != 0 {
		mod.Select(fields)
	}
	if len(whereSql) != 0 && len(vals) != 0 {
		mod.Where(whereSql, vals...)
	}
	result := mod.Debug().Find(&resAddress)
	return resAddress, result.RowsAffected, result.Error
}

func GetAddressCount(whereSql string, vals []interface{}) (count int64, err error) {
	mod := global.MysqlDb.Model(&Address{})
	if len(whereSql) != 0 && len(vals) != 0 {
		mod.Where(whereSql, vals...)
	}
	result := mod.Count(&count)
	return count, result.Error
}

func CreateAddress(address Address) (data Address, err error) {
	result := global.MysqlDb.Create(&address)
	return address, result.Error
}

func GetAddressFirst(whereSql string, vals []interface{}, fields string) (addressFirst Address, rows int64, err error) {
	mod := global.MysqlDb.Limit(1)
	if len(fields) != 0 {
		mod.Select(fields)
	}
	if len(whereSql) != 0 && len(vals) != 0 {
		mod.Where(whereSql, vals...)
	}
	result := mod.Find(&addressFirst)
	return addressFirst, result.RowsAffected, result.Error
}
