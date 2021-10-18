package model

import "project/goods_srv/global"

type Category struct {
	ID               uint32
	Name             string `gorm:"type:varchar(20);default:'';comment:'名称'"`
	ParentCategoryId uint32 `gorm:"index:category_parent_category_id;comment:'父类别';default:0"`
	Level            int32  `gorm:"type:tinyint(2) UNSIGNED;comment:'级别';default:0"`
	IsTab            bool   `gorm:"type:bool;comment:'是否显示在首页tab,1:是,0:否';default:false"`
	CreatedAt        uint32 `gorm:"comment:'添加时间';default:0"`
	UpdatedAt        uint32 `gorm:"comment:'更新时间';default:0"`
	IsDeleted        bool   `gorm:"type:bool;comment:'是否删除,1:是 0:否';default:false"`
	DeletedAt        uint32 `gorm:"comment:'删除时间';default:0"`
}

type Goodscategorybrand struct {
	ID         uint32
	CategoryId uint32 `gorm:"index:category_id;index:category_id_brand_id,unique;comment:'类别';default:0"`
	BrandId    uint32 `gorm:"index:brand_id;index:category_id_brand_id,unique;comment:'品牌';default:0"`
	CreatedAt  uint32 `gorm:"comment:'添加时间';default:0"`
	UpdatedAt  uint32 `gorm:"comment:'更新时间';default:0"`
	IsDeleted  bool   `gorm:"type:bool;comment:'是否删除,1:是 0:否';default:false"`
	DeletedAt  uint32 `gorm:"comment:'删除时间';default:0"`
}

func GetCategoryList(whereSql string,vals []interface{}, fields string, Offset int, limit int) (resCategory []Category, rows int64, err error) {
	mod := global.MysqlDb.Limit(limit).Offset(Offset)
	if len(fields) != 0 {
		mod.Select(fields)
	}
	if len(whereSql) != 0 && len(vals) != 0{
		mod.Where(whereSql,vals...)
	}
	result := mod.Find(&resCategory)
	return resCategory, result.RowsAffected, result.Error
}

func GetCategoryFirst(whereSql string,vals []interface{}, fields string) (categoryFirst Category, rows int64, err error) {
	mod := global.MysqlDb.Limit(1)
	if len(fields) != 0 {
		mod.Select(fields)
	}
	if len(whereSql) != 0 && len(vals) != 0{
		mod.Where(whereSql,vals...)
	}
	result := mod.Find(&categoryFirst)
	return categoryFirst, result.RowsAffected, result.Error
}


func CreateCategory(category Category) (data Category, err error) {
	result := global.MysqlDb.Create(&category)
	return category, result.Error
}

func UpdateCategory(data interface{}, whereSql string,vals []interface{}) (err error) {
	if data == nil || len(whereSql) == 0 || len(vals) == 0{
		return
	}
	result := global.MysqlDb.Model(&Category{}).Where(whereSql,vals...).Updates(data)
	return result.Error
}
