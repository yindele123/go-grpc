package model

import "project/userop_srv/global"

type Leavingmessages struct {
	ID          uint32
	User        uint32 `gorm:"comment:'用户id';default:0"`
	MessageType int32  `gorm:"type:tinyint(2) UNSIGNED;comment:'留言类型: 1(留言),2(投诉),3(询问),4(售后),5(求购)';default:1"`
	Subject     string `gorm:"type:varchar(100);default:'';comment:'主题'"`
	Message     string `gorm:"type:text;comment:'留言内容'"`
	File        string `gorm:"type:varchar(100);default:'';comment:'上传的文件'"`
	CreatedAt   uint32 `gorm:"comment:'添加时间';default:0"`
	UpdatedAt   uint32 `gorm:"comment:'更新时间';default:0"`
	IsDeleted   bool   `gorm:"type:bool;comment:'是否删除,1:是 0:否';default:false"`
	DeletedAt   uint32 `gorm:"comment:'删除时间';default:0"`
}

func GetMessagesList(whereSql string, vals []interface{}, fields string, Offset int, limit int) (resMessages []Leavingmessages, rows int64, err error) {
	mod := global.MysqlDb.Limit(limit).Offset(Offset)
	if len(fields) != 0 {
		mod.Select(fields)
	}
	if len(whereSql) != 0 && len(vals) != 0 {
		mod.Where(whereSql, vals...)
	}
	result := mod.Debug().Find(&resMessages)
	return resMessages, result.RowsAffected, result.Error
}

func GetMessagesCount(whereSql string, vals []interface{}) (count int64, err error) {
	mod := global.MysqlDb.Model(&Leavingmessages{})
	if len(whereSql) != 0 && len(vals) != 0 {
		mod.Where(whereSql, vals...)
	}
	result := mod.Count(&count)
	return count, result.Error
}

func CreateMessages(messages Leavingmessages) (data Leavingmessages, err error) {
	result := global.MysqlDb.Create(&messages)
	return messages, result.Error
}
