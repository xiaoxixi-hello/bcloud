package netfile

import (
	"bcloud/common/http"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/url"
	"strconv"
	"time"
)

type FileListInfo struct {
	Fsid     string
	FileName string
	FileSize string
	FileTime string
	FileType string
	FilePath string
	Isdir    int
}

type FileInfo struct {
	Fsid           uint64            `json:"fs_id"`
	Path           string            `json:"path"`
	ServerFilename string            `json:"server_filename"`
	Isdir          int               `json:"isdir"`
	Size           uint64            `json:"size"`
	Category       int               `json:"category"`
	Md5            string            `json:"md5"`
	DirEmpty       int               `json:"dir_empty"`
	LocalCtime     uint64            `json:"local_ctime"`
	LocalMtime     uint64            `json:"local_mtime"`
	ServerCtime    uint64            `json:"server_ctime"`
	ServerMtime    uint64            `json:"server_mtime"`
	Thumbs         map[string]string `json:"thumbs"` // 当文件类型为图片时，且请求参数含有web=1时，返回thumbs
}

type FilelistReturn struct {
	Cursor  int `json:"cursor"` //下一页查询起点
	Errno   int `json:"errno"`
	HasMore int `json:"has_more"` // 1 下一页还有数据 0 没有下一页

	List []FileInfo `json:"list"` //文件列表
}

type FilelistArg struct {
	Dir   string `json:"dir"`
	Order string `json:"order"`
	Start uint64 `json:"start"`
	Limit uint64 `json:"limit"`
	Desc  int    `json:"desc"`
	Token string `json:"token"`
}

func NewFilelistArg(dir string, order string, start uint64, limit uint64, desc int) *FilelistArg {
	s := new(FilelistArg)
	s.Dir = dir
	s.Start = start
	s.Limit = limit
	s.Order = order
	s.Desc = desc
	return s
}

// FileList 获取单个目录的文件列表
func FileList(arg *FilelistArg, token string) (bool, []*FileListInfo, int, error) {
	ret := FilelistReturn{}
	fileListInfos := make([]*FileListInfo, 0)

	protocal := "https"
	host := "pan.baidu.com"
	router := "/rest/2.0/xpan/file?method=list&"
	uri := protocal + "://" + host + router

	params := url.Values{}
	params.Set("access_token", token)
	params.Set("dir", arg.Dir)
	params.Set("start", strconv.FormatUint(arg.Start, 10))
	params.Set("limit", strconv.FormatUint(arg.Limit, 10))
	params.Set("order", arg.Order)
	params.Set("desc", strconv.FormatInt(int64(arg.Desc), 10))
	params.Set("folder", "0")
	params.Set("web", "0")
	params.Set("show_empty", "0")
	uri += params.Encode()

	headers := map[string]string{
		"Host":         host,
		"Content-Type": "application/x-www-form-urlencoded",
	}

	var postBody io.Reader
	body, _, err := http.DoHTTPRequest(uri, postBody, headers)
	if err != nil {
		zap.L().Error("获取文件列表失败", zap.Error(err))
		return false, nil, 0, errors.New("获取文件列表失败")
	}
	if err = json.Unmarshal([]byte(body), &ret); err != nil {
		zap.L().Error("获取文件列表失败", zap.Error(errors.New("unmarshal fileListAll body failed,body")))
		return false, nil, 0, errors.New("unmarshal fileListAll body failed,body")
	}
	if ret.Errno != 0 {
		if ret.Errno == 31034 {
			return false, nil, 0, errors.New("接口请求限频")
		}
		zap.L().Info(fmt.Sprintf("灾难性错误，获取文件列表失败！ret:%v", ret))
	}
	for _, info := range ret.List {
		f := &FileListInfo{
			FileName: info.ServerFilename,
			Fsid:     strconv.Itoa(int(info.Fsid)),
			FileTime: time.Unix(int64(info.LocalMtime), 0).Format("2006-01-02 15:04:05"),
			FilePath: info.Path,
			Isdir:    info.Isdir,
		}
		if info.Isdir == 1 {
			f.FileType = "文件夹"
		} else {
			f.FileType = "文件"
		}
		if info.Size == 0 {
			f.FileSize = "--"
		} else if info.Size/1024/1024/1024 < 1 {
			f.FileSize = fmt.Sprintf("%dMB", info.Size/1024/1024)
		} else {
			f.FileSize = fmt.Sprintf("%dGB", info.Size/1024/1024/1024)
		}
		fileListInfos = append(fileListInfos, f)
	}
	if ret.HasMore == 1 {
		return true, fileListInfos, ret.Cursor, nil
	}
	return false, fileListInfos, 0, nil
}

func FileListAll(path, token string) ([]*FileListInfo, int) {
	var start, limit uint64
	fileDownInfoList := make([]*FileListInfo, 0)
	start, limit = 0, 1000
	arg := NewFilelistArg(path, "time", start, limit, 1)
	for {
		b, f, l, err := fileListAll(arg, token)
		if err != nil {
			if err.Error() == "接口请求限频" {
				zap.L().Error("FileListAll", zap.Error(err))
				return fileDownInfoList, -1
			}
		}
		if err != nil || !b {
			fileDownInfoList = append(fileDownInfoList, f...)
			break
		}
		start = limit
		limit = uint64(l + 1000)
		arg = NewFilelistArg(path, "time", start, limit, 1)
		fileDownInfoList = append(fileDownInfoList, f...)
	}
	return fileDownInfoList, 0
}

// FileListAll 递归获取文件列表
func fileListAll(arg *FilelistArg, token string) (bool, []*FileListInfo, int, error) {
	ret := FilelistReturn{}
	fileDownInfoList := make([]*FileListInfo, 0)

	protocal := "https"
	host := "pan.baidu.com"
	router := "/rest/2.0/xpan/multimedia?method=listall&"
	uri := protocal + "://" + host + router

	params := url.Values{}
	params.Set("access_token", token)
	params.Set("path", arg.Dir)
	params.Set("start", strconv.FormatUint(arg.Start, 10))
	params.Set("limit", strconv.FormatUint(arg.Limit, 10))
	params.Set("order", arg.Order)
	params.Set("desc", strconv.FormatInt(int64(arg.Desc), 10))
	params.Set("recursion", "1")
	uri += params.Encode()

	headers := map[string]string{
		"Host":         host,
		"Content-Type": "application/x-www-form-urlencoded",
	}

	var postBody io.Reader
	body, _, err := http.DoHTTPRequest(uri, postBody, headers)
	if err != nil {
		zap.L().Error("递归获取文件失败", zap.Error(err))
		return false, nil, 0, errors.New("递归获取文件失败")
	}
	if err = json.Unmarshal([]byte(body), &ret); err != nil {
		zap.L().Error("递归获取文件失败", zap.Error(errors.New("unmarshal fileListAll body failed,body")))
		return false, nil, 0, errors.New("unmarshal fileListAll body failed,body")
	}
	if ret.Errno != 0 {
		if ret.Errno == 31034 {
			return false, nil, 0, errors.New("接口请求限频")
		}
		zap.L().Info(fmt.Sprintf("灾难性错误，获取文件列表失败！ret:%v", ret))
	}
	for _, info := range ret.List {
		f := &FileListInfo{
			FileName: info.ServerFilename,
			Fsid:     strconv.Itoa(int(info.Fsid)),
			FileTime: time.Unix(int64(info.LocalMtime), 0).Format("2006-01-02 15:04:05"),
			FilePath: info.Path,
		}
		if info.Isdir == 1 {
			f.FileType = "文件夹"
		} else {
			f.FileType = "文件"
		}
		if info.Size == 0 {
			f.FileSize = "--"
		} else if info.Size/1024/1024/1024 < 1 {
			f.FileSize = fmt.Sprintf("%dMB", info.Size/1024/1024)
		} else {
			f.FileSize = fmt.Sprintf("%dGB", info.Size/1024/1024/1024)
		}
		fileDownInfoList = append(fileDownInfoList, f)
	}
	if ret.HasMore == 1 {
		return true, fileDownInfoList, ret.Cursor, nil
	}
	return false, fileDownInfoList, 0, nil
}
