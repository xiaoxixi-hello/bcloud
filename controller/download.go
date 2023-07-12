package controller

import (
	"bcloud/dao/mysql"
	"bcloud/netdisk/download"
	"bcloud/netdisk/netfile"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.uber.org/zap"
	"strconv"
)

func (a *App) DownFilePre(data []netfile.FileListInfo) int {
	var fileFsid []uint64
	for _, info := range data {
		if info.Isdir == 1 {
			fileListInfos, i := netfile.FileListAll(info.FilePath, a.Config["PublicToken"])
			if i == -1 {
				return -2
			}
			for _, listInfo := range fileListInfos {
				parseUint1, _ := strconv.ParseUint(listInfo.Fsid, 10, 64)
				fileFsid = append(fileFsid, parseUint1)
			}
		} else {
			parseUint, _ := strconv.ParseUint(info.Fsid, 10, 64)
			fileFsid = append(fileFsid, parseUint)
		}
	}
	if len(fileFsid) > 500 {
		return -1
	}
	fileMetasAll := netfile.FileMetasAll(a.Config["PublicToken"], fileFsid)
	if len(fileMetasAll) != 0 {
		mysql.InsertListDownDetail(a.DB, fileMetasAll)
	}
	go func() {
		a.DownFIle(fileMetasAll)
	}()
	result, err := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:          runtime.InfoDialog,
		Title:         "下载任务提交成功",
		Message:       "请前往下载列表查看进程!",
		DefaultButton: "No",
	})
	zap.L().Info(fmt.Sprintf("任务提交信息:%v,%v,%v", result, err, data))
	return 0
}

func (a *App) DownFIle(f []*download.DownDetail) {
	if err := download.MultiThreadDownRun(f, a.Config["DownPath"], "99", a.Config["PublicToken"], a.DB, a.Config["MaxDownProcess"]); err != nil {
		zap.L().Error("下载失败", zap.Error(err))
		return
	}
}

func (a *App) downFromDB() {
	go func() {
		a.DownFIle(mysql.FindListDownDetail(a.DB))
	}()
}

func (a *App) GetDownListDetail() []download.DownDetail {
	return mysql.FindListOrderId(a.DB)
}

func (a *App) DownRetry(id int) {
	d := mysql.FindDownDetailForId(a.DB, id)
	_, _ = a.DB.Exec("DELETE FROM tb_down_detail WHERE id = ?", id)
	downDetails := netfile.FileMetasAll(a.Config["PublicToken"], []uint64{d.Fsid})
	downDetails[0].Status = "重试任务"
	mysql.InsertListDownDetail(a.DB, downDetails)
	_ = download.MultiThreadDownRun(downDetails, a.Config["DownPath"], "99", a.Config["PublicToken"], a.DB, a.Config["MaxDownProcess"])
}
func (a *App) DeleteRetry(id int) {
	_, _ = a.DB.Exec("DELETE FROM tb_down_detail WHERE id = ?", id)
}
