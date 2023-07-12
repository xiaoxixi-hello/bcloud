package netfile

import (
	"bcloud/common/http"
	"bcloud/netdisk/download"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/url"
	"strconv"
	"sync"
	"time"
)

type FileMetasArg struct {
	Fsids []uint64 `json:"fsids"`
	Path  string   `json:"path"` //查询共享目录或专属空间内文件时需要
}

func NewFileMetasArg(fsid []uint64, path string) *FileMetasArg {
	s := new(FileMetasArg)
	s.Fsids = fsid
	s.Path = path
	return s
}

type ListInfo struct {
	Size        uint64            `json:"size"`
	Path        string            `json:"path"`
	Isdir       int               `json:"isdir"`
	ServerCtime uint64            `json:"server_ctime"`
	ServerMtime uint64            `json:"server_mtime"`
	Fsid        uint64            `json:"fs_id"`
	OperId      int               `json:"oper_id"`
	Md5         string            `json:"md5"`
	Filename    string            `json:"filename"`
	Category    int               `json:"category"`
	Dlink       string            `json:"dlink"` // 文件才返回dlink
	Duration    int               `json:"duration"`
	Thumbs      map[string]string `json:"thumbs"`
	Height      int               `json:"height"`
	Width       int               `json:"width"`
	DateTaken   int               `json:"date_taken"`
}
type FileMetasReturn struct {
	Errno     int                    `json:"errno"`
	Errmsg    string                 `json:"errmsg"`
	RequestID string                 `json:"request_id"`
	Names     map[string]interface{} `json:"names"`
	List      []ListInfo             `json:"list"`
}

// FileMetas 获取文件详细信息
func fileMetas(token string, arg *FileMetasArg) (FileMetasReturn, error) {
	ret := FileMetasReturn{}

	protocal := "https"
	host := "pan.baidu.com"
	router := "/rest/2.0/xpan/multimedia?method=filemetas&"
	uri := protocal + "://" + host + router

	params := url.Values{}
	params.Set("access_token", token)
	fsidJs, err := json.Marshal(arg.Fsids)
	if err != nil {
		return ret, err
	}
	params.Set("fsids", string(fsidJs))
	params.Set("dlink", "1") // 对于下载，dlink为必选参数，才能拿到dlink下载地址
	params.Set("path", arg.Path)
	uri += params.Encode()

	headers := map[string]string{
		"Host":         host,
		"Content-Type": "application/x-www-form-urlencoded",
	}

	var postBody io.Reader
	body, _, err := http.DoHTTPRequest(uri, postBody, headers)
	if err != nil {
		return ret, err
	}
	if err = json.Unmarshal([]byte(body), &ret); err != nil {
		return ret, errors.New("unmarshal filemetas body failed,body")
	}
	if ret.Errno != 0 {
		return ret, errors.New("call filemetas failed")
	}
	return ret, nil
}

func FileMetasAll(token string, fileFsid []uint64) []*download.DownDetail {
	rf := make([]*download.DownDetail, 0)
	if len(fileFsid) < 100 {
		metasArg := NewFileMetasArg(fileFsid, "")
		metas, err := fileMetas(token, metasArg)
		if err != nil {
			zap.L().Error("获取文件详情失败", zap.Error(err))
			panic("获取文件详情失败")
		}
		for _, info := range metas.List {
			if info.Isdir != 1 {
				f := &download.DownDetail{
					CreatedAt: time.Now().Format("2006-01-02"),
					Name:      info.Filename,
					Path:      info.Path,
					Status:    "初始化完成",
					Dlink:     info.Dlink,
					Fsid:      info.Fsid,
				}
				if info.Size == 0 {
					f.Size = "--"
				} else if info.Size/1024/1024/1024 < 1 {
					f.Size = fmt.Sprintf("%dMB", info.Size/1024/1024)
				} else {
					f.Size = fmt.Sprintf("%vGB", strconv.FormatFloat(float64(info.Size)/1024/1024/1024, 'f', 2, 64))
				}
				rf = append(rf, f)
			}
		}
		return rf
	}
	count := len(fileFsid)/100 + 1
	var wg sync.WaitGroup
	var mutex sync.Mutex
	wg.Add(count)
	for i := 0; i < count; i++ {
		end := (i + 1) * 100
		if end > len(fileFsid) { // 最后一段
			end = len(fileFsid)
		}
		go func(i int, end int) {
			defer wg.Done()
			metasArg := NewFileMetasArg(fileFsid[i:end], "")
			metas, err := fileMetas(token, metasArg)
			if err != nil {
				zap.L().Info("获取文件详情失败", zap.Error(err))
				panic("获取文件详情失败")
			}
			for _, info := range metas.List {
				mutex.Lock()
				if info.Isdir != 1 {
					f := &download.DownDetail{
						CreatedAt: time.Now().Format("2006-01-02"),
						Name:      info.Filename,
						Path:      info.Path,
						Status:    "初始化完成",
						Dlink:     info.Dlink,
						Fsid:      info.Fsid,
					}
					if info.Size == 0 {
						f.Size = "--"
					} else if info.Size/1024/1024/1024 < 1 {
						f.Size = fmt.Sprintf("%dMB", info.Size/1024/1024)
					} else {
						f.Size = fmt.Sprintf("%dGB", info.Size/1024/1024/1024)
					}
					rf = append(rf, f)
				}
				mutex.Unlock()
			}
		}(i*100, end)
	}
	wg.Wait()
	return rf
}
