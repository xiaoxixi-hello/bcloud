package mysql

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func InitDB() *gorm.DB {
	host := "114.115.175.220"
	port := 45672
	username := "root"
	password := "2025907338"
	dbname := "db_bcloud"

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, dbname)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "tb_", //表名前缀
			SingularTable: true,  // 单数表名
		},
	})

	if err != nil {
		zap.L().Error("连接云数据库失败", zap.Error(err))
		panic("连接云数据库失败")
	}
	zap.L().Info("云数据库连接成功")
	//_ = db.AutoMigrate(controller.SharLinkInfo{}, controller.Authority{})
	return db
}

func InitLocalDB(path string) *gorm.DB {
	//host := "127.0.0.1"
	//port := 63306
	//username := "root"
	//password := "123456"
	//dbname := "db_bcloud"

	//	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, dbname)
	//	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{
	//	NamingStrategy: schema.NamingStrategy{
	//		TablePrefix:   "tb_", //表名前缀
	//		SingularTable: true,  // 单数表名
	//	},
	//})
	db, err := gorm.Open(sqlite.Open(path))
	if err != nil {
		panic("初始化本地数据库失败")
	}

	// 自动迁移数据库表结构
	zap.L().Info("本地数据库创建成功")
	return db
}
