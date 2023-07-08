package controller

import (
	"bcloud/common/logger"
	"bcloud/dao/mysql"
	"bcloud/netdisk/download"
	"bcloud/netdisk/floder"
	"context"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type App struct {
	ctx    context.Context
	DB, CB *gorm.DB          // 数据库信息
	Config map[string]string // 系统配置
}

func NewApp() *App {
	return &App{}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	/*
		1. 日志
		2. 数据库
		3. 系统配置
		4. 获取key、获取授权
	*/
	logger.InitLogger("debug", "dev")
	a.DB = mysql.InitLocalDB(fmt.Sprintf("%s/bcloud.db", floder.GetConfigDir()))
	_ = a.DB.AutoMigrate(&ConfigItem{}, &download.DownDetail{})
	a.CB = mysql.InitDB()

	a.getToken()
	a.Config = a.GetConfig()
	a.GetAuthorityCount()

	a.downFromDB()
}

// OnBeforeClose action
func (a *App) OnBeforeClose(ctx context.Context) bool {
	return true
}

// OnDOMReady action
func (a *App) OnDOMReady(ctx context.Context) {
	// 启动一个监听事件
	//runtime.EventsOn(a.ctx, "test", func(optionalData ...interface{}) {
	//	a.Log.Info(optionalData...)
	//})
}

func (a *App) Shutdown(ctx context.Context) {
	zap.L().Info("Shutting down. Goodbye!")
}
