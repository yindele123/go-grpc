package model

import (
	"project/user_srv/global"
)

type User struct {
	ID        uint32
	Mobile    string `gorm:"type:varchar(11);unique;not null;comment:'手机号码'"`
	Password  string `gorm:"type:varchar(100);not null;comment:'密码'"`
	NickName  string `gorm:"type:varchar(20);default:'';comment:'昵称'"`
	HeadUrl   string `gorm:"type:varchar(200);default:'';comment:'头像'"`
	Birthday  uint64 `gorm:"comment:'生日';default:0"`
	Address   string `gorm:"type:varchar(200);comment:'地址';default:''"`
	Desc      string `gorm:"type:text;comment:'个人简历'"`
	Gender    uint32 `gorm:"type:tinyint(2) UNSIGNED;comment:'性别  1:女  2:男  3::保密';default:3"`
	Role      uint32 `gorm:"type:tinyint(2) UNSIGNED;comment:'用户角色,1普通用户 2管理员';default:1"`
	CreatedAt uint32 `gorm:"comment:'添加时间';default:0"`
	UpdatedAt uint32 `gorm:"comment:'更新时间';default:0"`
	DeletedAt uint32 `gorm:"comment:'删除时间';default:0"`
}

func GetUserList(where interface{}, fields string, Offset int, limit int) (User []User, rows int64, err error) {
	mod := global.MysqlDb.Limit(limit).Offset(Offset)
	if len(fields) != 0 {
		mod.Select(fields)
	}
	if where != nil {
		mod.Where(where)
	}
	result := mod.Find(&User)
	return User, result.RowsAffected, result.Error
}

func GetUserCount(where interface{}) (count int64, err error) {
	mod := global.MysqlDb.Model(&User{})
	if where != nil {
		mod.Where(where)
	}
	result := mod.Count(&count)
	return count, result.Error
}

func GetUserFirst(where interface{}, fields string) (UserFirst User, rows int64, err error) {
	mod := global.MysqlDb.Limit(1)
	if len(fields) != 0 {
		mod.Select(fields)
	}
	if where != nil {
		mod.Where(where)
	}
	result := mod.Find(&UserFirst)
	return UserFirst, result.RowsAffected, result.Error
}

func UpdateUser(data interface{}, where interface{}) (err error) {
	if data == nil || where == nil {
		return
	}
	result := global.MysqlDb.Model(&User{}).Where(where).Updates(data)
	return result.Error
}

func CreateUser(user User) (data User, err error) {
	result := global.MysqlDb.Create(&user)
	return user, result.Error
}
