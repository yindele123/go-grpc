package model

import "project/goods_srv/global"

type Banners struct {
	ID        uint32
	Image     string `gorm:"type:varchar(200);default:'';comment:'图片url'"`
	Url       string `gorm:"type:varchar(200);default:'';comment:'访问url'"`
	Index     uint32 `gorm:"comment:'轮播顺序';default:0"`
	CreatedAt uint32 `gorm:"comment:'添加时间';default:0"`
	UpdatedAt uint32 `gorm:"comment:'更新时间';default:0"`
	IsDeleted bool   `gorm:"type:bool;comment:'是否删除,1:是 0:否';default:false"`
	DeletedAt uint32 `gorm:"comment:'删除时间';default:0"`
}

func GetBannersList(whereSql string, vals []interface{}, fields string, Offset int, limit int, order string) (resBanners []Banners, rows uint32, err error) {
	mod := global.MysqlDb.Limit(limit).Offset(Offset)
	if len(fields) != 0 {
		mod.Select(fields)
	}
	if len(order) != 0 {
		mod.Order(order)
	}
	if len(whereSql) != 0 && len(vals) != 0 {
		mod.Where(whereSql, vals...)
	}
	result := mod.Find(&resBanners)
	return resBanners, uint32(result.RowsAffected), result.Error
}

func GetBannersCount(whereSql string, vals []interface{}) (resCount uint32, err error) {
	mod := global.MysqlDb.Model(&Banners{})
	if len(whereSql) != 0 && len(vals) != 0 {
		mod.Where(whereSql, vals...)
	}
	var count int64
	result := mod.Count(&count)
	return uint32(count), result.Error
}

func CreateBanners(banners Banners) (data Banners, err error) {
	result := global.MysqlDb.Create(&banners)
	return banners, result.Error
}

func UpdateBanners(data interface{}, whereSql string, vals []interface{}) (err error) {
	if data == nil || len(whereSql) == 0 || len(vals) == 0 {
		return
	}
	result := global.MysqlDb.Model(&Banners{}).Where(whereSql, vals...).Updates(data)
	return result.Error
}

func GetBannersFirst(whereSql string, vals []interface{}, fields string) (bannerFirst Banners, rows uint32, err error) {
	mod := global.MysqlDb.Limit(1)
	if len(fields) != 0 {
		mod.Select(fields)
	}
	if len(whereSql) != 0 && len(vals) != 0 {
		mod.Where(whereSql, vals...)
	}
	result := mod.Find(&bannerFirst)
	return bannerFirst, uint32(result.RowsAffected), result.Error
}
