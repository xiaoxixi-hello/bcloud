package controller

import (
	"bcloud/dao/mail"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.uber.org/zap"
	"strings"
)

// 链接相关

type SharLinkInfo struct {
	ID   int `gorm:"primarykey"`
	Path string
	Link string
}

func (a *App) CreateShareLinkInfo(link string) int {
	if a.GetAuthorityCount() < 1 {
		//_, _ = runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		//	Type:          runtime.WarningDialog,
		//	Title:         "资源链接提交失败",
		//	Message:       "当前可用提交次数为0",
		//	DefaultButton: "No",
		//})
		return -1
	}
	_, _ = runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:          runtime.InfoDialog,
		Title:         "资源链接提交成功",
		Message:       "资源转存需要花费一点时间,请勿频繁提交",
		DefaultButton: "No",
	})

	if a.CB.Create(&SharLinkInfo{Path: a.Config["PanID"], Link: link}).Error != nil {
		panic("插入分享链接失败")
	}
	zap.L().Info("提交链接信息:" + link)
	if err := mail.SendToMe(fmt.Sprintf("目标路径:%s\n链接:%s", link, a.Config["PanID"])); err != nil {
		panic(err)
	}

	var auth Authority
	a.CB.Take(&auth, "path = ?", strings.TrimLeft(a.Config["PanID"], "/"))
	auth.Count = auth.Count - 1
	a.CB.Select("count").Save(&auth)
	return auth.Count
}
