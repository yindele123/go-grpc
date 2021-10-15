package model

import "project/goods_srv/global"

type Brands struct {
	ID        uint32
	Name      string `gorm:"uniqueIndex:brands_name;type:varchar(50);default:'';comment:'名称'"`
	Logo      string `gorm:"type:varchar(200);default:'';comment:'图标'"`
	CreatedAt uint32 `gorm:"comment:'添加时间';default:0"`
	UpdatedAt uint32 `gorm:"comment:'更新时间';default:0"`
	IsDeleted bool   `gorm:"type:bool;comment:'是否删除,1:是 0:否';default:false"`
	DeletedAt uint32 `gorm:"comment:'删除时间';default:0"`
}

func GetBrandsList(whereSql string,vals []interface{}, fields string, Offset int, limit int) (resBrands []Brands, rows int64, err error) {
	mod := global.MysqlDb.Limit(limit).Offset(Offset)
	if len(fields) != 0 {
		mod.Select(fields)
	}
	if len(whereSql) != 0 && len(vals) != 0{
		mod.Where(whereSql,vals...)
	}
	result := mod.Debug().Find(&resBrands)
	return resBrands, result.RowsAffected, result.Error
}

func GetBrandsFirst(whereSql string,vals []interface{}, fields string) (brandsFirst Brands, rows int64, err error) {
	mod := global.MysqlDb.Limit(1)
	if len(fields) != 0 {
		mod.Select(fields)
	}
	if len(whereSql) != 0 && len(vals) != 0{
		mod.Where(whereSql,vals...)
	}
	result := mod.Find(&brandsFirst)
	return brandsFirst, result.RowsAffected, result.Error
}

