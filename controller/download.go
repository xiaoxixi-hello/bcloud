package controller

import (
	"bcloud/netdisk/download"
	"bcloud/netdisk/netfile"
	"errors"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.uber.org/zap"
	"gorm.io/gorm"
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
		if a.DB.Create(&fileMetasAll).Error != nil {
			panic("任务信息写入数据库失败")
		}
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
	var list []*download.DownDetail
	if a.DB.Find(&list, "status != ?", "下载完成").Error != nil {
		zap.L().Error("从获取下载信息失败")
		return
	}
	go func() {
		a.DownFIle(list)
	}()
}

func (a *App) GetDownListDetail() []download.DownDetail {
	var list []download.DownDetail
	if a.DB.Order("id desc").Find(&list).Error != nil {
		zap.L().Error("从获取下载信息失败")
		return nil
	}
	return list
}

func (a *App) DownRetry(id int) {
	var d download.DownDetail
	a.DB.Find(&d, "id = ?", id)
	downDetails := netfile.FileMetasAll(a.Config["PublicToken"], []uint64{d.Fsid})
	downDetails[0].Status = "重试任务"
	if a.DB.Create(&downDetails).Error != nil {
		panic("任务信息写入数据库失败")
	}
	_ = download.MultiThreadDownRun(downDetails, a.Config["DownPath"], "99", a.Config["PublicToken"], a.DB, a.Config["MaxDownProcess"])

	if a.DB.Delete(&download.DownDetail{}, "id = ?", id).Error != nil {
		panic("删除失败")
	}
}
func (a *App) DeleteRetry(id int) {
	var d download.DownDetail
	err := a.DB.Debug().Take(&d, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}
	if d.Status != "下载完成" {
		d.ProcessStatus = 1
		if a.DB.Debug().Select("process_status").Save(&d).Error != nil {
			panic("删除状态更新下载进程失败")
		}
		return
	}
	if a.DB.Debug().Delete(&d, "id = ?", id).Error != nil {
		panic("删除失败")
	}
}
