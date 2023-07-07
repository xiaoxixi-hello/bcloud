package controller

import (
	"bcloud/common/snowflake"
	"bcloud/netdisk/floder"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"os"
)

var (
	ConfigTypes = [...]string{"PublicToken", "PanID", "DownPath", "MaxDownProcess"}
)

type ConfigItem struct {
	Id    int64  `json:"id"`    // ID
	Name  string `json:"name"`  // 配置项名称
	Value string `json:"value"` // 配置项值
}

type BaiduApiKey struct {
	Token string
}

func (a *App) GetConfig() map[string]string {
	config := make(map[string]string)
	for _, name := range ConfigTypes {
		var c ConfigItem
		a.DB.Take(&c, "name = ?", name)

		if c.Value == "" && name != "PanID" {
			c.Value = a.initConfigData(name, "")
		}
		if name == "PanID" && c.Value == "" {
			c.Value = a.initConfigData(name, config["PublicToken"])
		}
		config[name] = c.Value

	}
	zap.L().Info(fmt.Sprintf("系统配置初始化完成!,%v", config))
	return config
}

func (a *App) getToken() {
	var b BaiduApiKey
	var c ConfigItem
	if a.CB.Find(&b).Error != nil {
		panic("获取百度网盘token失败")
	}
	err := a.DB.Take(&c, "name = ?", "PublicToken").Error
	if err == gorm.ErrRecordNotFound {
		a.DB.Create(&ConfigItem{Name: "PublicToken", Value: b.Token})
		return
	}
	a.DB.Select("PublicToken").Save(&c)
	zap.L().Info("获取Token成功" + b.Token)
}

func (a *App) initConfigData(name, token string) string {
	switch name {
	case "PanID":
		panID := snowflake.InitPanID()
		floder.CreateDir(panID, token)
		if a.DB.Create(&ConfigItem{Name: "PanID", Value: "/" + panID}).Error != nil {
			panic("初始化网盘ID失败")
		}
		zap.L().Info(fmt.Sprintf("初始化网盘成功,PanID %s", panID))
		return "/" + panID
	case "DownPath":
		if a.DB.Create(&ConfigItem{Name: "DownPath", Value: floder.GetConfigDir()}).Error != nil {
			panic("初始化网盘下载路径失败")
		}
		zap.L().Info(fmt.Sprintf("初始化网盘下载路径成功,DownPath %s", floder.GetConfigDir()))
		return floder.GetConfigDir()
	case "MaxDownProcess":
		if a.DB.Create(&ConfigItem{Name: "MaxDownProcess", Value: "4"}).Error != nil {
			panic("初始化网盘并发下载线程数失败")
		}
		zap.L().Info(fmt.Sprintf("初始化网盘并发下载线程,MaxDownProcess %d", 2))
		return "2"
	default:
		zap.L().Error("没有找到该系统配置,无法初始化" + name)
		os.Exit(1)
	}
	return ""
}

func (a *App) UpdateConfigItem(config ConfigItem) {
	var oldConfig ConfigItem
	a.DB.Find(&oldConfig, "name = ?", config.Name)
	if oldConfig.Value != config.Value {
		oldConfig.Value = config.Value
		a.DB.Select("value").Save(&oldConfig)
		a.GetConfig()
	}
}
func (a *App) ChangeDownPath() string {
	dir, _ := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{})
	var s ConfigItem
	a.DB.Take(&s, "name = ?", "DownPath")
	if dir != "" {
		s.Value = dir
		a.DB.Select("value").Save(&s)
		a.GetConfig()
		return dir
	}
	return s.Value
}
