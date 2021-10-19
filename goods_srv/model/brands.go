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

func GetBrandsList(whereSql string,vals []interface{}, fields string, Offset int, limit int) (resBrands []Brands, rows uint32, err error) {
	mod := global.MysqlDb.Limit(limit).Offset(Offset)
	if len(fields) != 0 {
		mod.Select(fields)
	}
	if len(whereSql) != 0 && len(vals) != 0{
		mod.Where(whereSql,vals...)
	}
	result := mod.Find(&resBrands)
	return resBrands, uint32(result.RowsAffected), result.Error
}

func GetBrandsFirst(whereSql string,vals []interface{}, fields string) (brandsFirst Brands, rows uint32, err error) {
	mod := global.MysqlDb.Limit(1)
	if len(fields) != 0 {
		mod.Select(fields)
	}
	if len(whereSql) != 0 && len(vals) != 0{
		mod.Where(whereSql,vals...)
	}
	result := mod.Find(&brandsFirst)
	return brandsFirst, uint32(result.RowsAffected), result.Error
}

func GetBrandsCount(whereSql string, vals []interface{}) (resCount uint32, err error) {
	mod := global.MysqlDb.Model(&Brands{})
	if len(whereSql) != 0 && len(vals) != 0 {
		mod.Where(whereSql, vals...)
	}
	var count int64
	result := mod.Count(&count)
	return uint32(count), result.Error
}

func CreateBrands(brands Brands) (data Brands, err error) {
	result := global.MysqlDb.Create(&brands)
	return brands, result.Error
}

func UpdateBrands(data interface{}, whereSql string, vals []interface{}) (err error) {
	if data == nil || len(whereSql) == 0 || len(vals) == 0 {
		return
	}
	result := global.MysqlDb.Model(&Brands{}).Where(whereSql, vals...).Updates(data)
	return result.Error
}

