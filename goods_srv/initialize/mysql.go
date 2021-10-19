package initialize

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"project/goods_srv/global"
	"time"
)

//MysqlConfig
func InitMysql() {

	/*dsn := "root:root@tcp(192.168.226.108:3306)/goods_srv?charset=utf8mb4&parseTime=True&loc=Local"*/
	user:=global.ServerConfig.MysqlInfo.User
	password:=global.ServerConfig.MysqlInfo.Password
	host:=global.ServerConfig.MysqlInfo.Host
	port:=global.ServerConfig.MysqlInfo.Port
	dbname:=global.ServerConfig.MysqlInfo.DbName
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password,host,port, dbname)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   global.ServerConfig.MysqlInfo.TablePrefix, // 表名前缀，`User`表为`t_users`
			SingularTable: true, // 使用单数表名，启用该选项后，`User` 表将是`user`
		},
	})
	if err != nil {
		panic(err)
	}
	sqlDB, _ := db.DB()

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
	//db.AutoMigrate(&model.Goods{},&model.Banners{},&model.Brands{},&model.Category{},&model.Goodscategorybrand{})
	//MysqlDb
	global.MysqlDb=db
}
