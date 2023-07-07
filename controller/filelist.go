package controller

import (
	"bcloud/netdisk/netfile"
	"strings"
)

func (a *App) GetTopList(path string) []string {
	return strings.Split(path, "/")[1:]
}

func (a *App) GetFileList(path string) []*netfile.FileListInfo {
	var start, limit uint64
	fileListInfos := make([]*netfile.FileListInfo, 0)
	start, limit = 0, 1000
	arg := netfile.NewFilelistArg(path, "time", start, limit, 1)
	for {
		b, f, l, err := netfile.FileList(arg, a.Config["PublicToken"])
		if err != nil || !b {
			fileListInfos = append(fileListInfos, f...)
			break
		}
		start = limit
		limit = uint64(l + 1000)
		arg = netfile.NewFilelistArg(path, "time", start, limit, 1)
		fileListInfos = append(fileListInfos, f...)
	}
	return fileListInfos
}

func (a *App) GetFileOtherList(path string, name string) string {
	p := strings.Split(path, "/")[1:]
	var index int
	for k, v := range p {
		if v == name {
			index = k
			break
		}
	}
	var oPath string
	for _, v := range p[:index+1] {
		oPath = oPath + "/" + v
	}
	return oPath
}
