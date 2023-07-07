package download

import (
	"bcloud/netdisk/floder"
	"crypto/tls"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type DownDetail struct {
	ID            int64  `gorm:"primarykey"`
	CreatedAt     string // 创建时间
	Name          string // 下载的名字
	Path          string
	Size          string // 文件的大小
	Status        string // 下载的状态 百分比、下载完成、下载失败
	Dlink         string // 文件的下载链接
	Fsid          uint64
	ProcessStatus int // 删除正在下载任务的标识
}

type FileDownloader struct {
	fileSize   int64
	dir        string
	fileName   string
	url        string
	md5        string
	blockSize  int64
	retryTimes int
}

func download(dlink, filename, base, id, token string, db *gorm.DB) error {
	path, name, err := floder.DownDirPath(filename, base)
	if err != nil {
		zap.L().Error("递归创建文件夹失败", zap.Error(err))
		return err
	}

	f := FileDownloader{
		url:        dlink + "&" + "access_token=" + token,
		fileName:   name,
		dir:        path,
		blockSize:  5 * 1024 * 1024,
		retryTimes: 5,
	}
	headers := map[string]string{
		"User-Agent": "pan.baidu.com",
	}
	timeout := 60 * time.Second
	tr := &http.Transport{
		MaxIdleConnsPerHost: -1,
		MaxIdleConns:        10,
		IdleConnTimeout:     30 * time.Second,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: tr}
	httpClient.Timeout = timeout

	var postBody io.Reader
	req, err := http.NewRequest("POST", f.url, postBody)
	if err != nil {
		zap.L().Error("http请求失败", zap.Error(err))
		return err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	var res *http.Response
	for i := 0; i < f.retryTimes; i++ {
		res, err = httpClient.Do(req)
		if err == nil {
			break
		}
		if i == f.retryTimes {
			zap.L().Error("http请求失败", zap.Error(err))
		}
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		upDataProcess(dlink, "下载失败", db)
		return errors.New("http返回码异常")
	}
	f.fileSize, err = strconv.ParseInt(res.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		zap.L().Error("获取HTTP大小失败", zap.Error(err))
		return err
	}

	// 创建本地临时文件
	tmpFile, err := os.CreateTemp(f.dir, "download_*.tmp")
	if err != nil {
		zap.L().Error("创建本地临时文件夹失败", zap.Error(err))
		return err
	}
	defer os.Remove(f.dir + tmpFile.Name())
	defer tmpFile.Close()

	var downloadedSize int64
	var buf = make([]byte, 1024)
	for {
		startTime := time.Now()
		start, end := downloadedSize, downloadedSize+f.blockSize-1
		if end >= f.fileSize {
			end = f.fileSize - 1
		}
		//fmt.Println("循环:", start, end)
		// 根据已下载的数据位置发送 Range 请求
		var fres *http.Response
		for i := 0; i < 5; i++ {
			req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))
			fres, err = http.DefaultClient.Do(req)
			if err == nil {
				break
			}
			zap.L().Info(fmt.Sprintf("在range下载时,获取响应失败,正在进行%d次重试", i+1), zap.String("name", filename), zap.String("进程任务ID", id))
			time.Sleep(3 * time.Second)
		}

		// 写入数据到临时文件
		n, err := io.CopyBuffer(tmpFile, fres.Body, buf)
		if err != nil {
			zap.L().Info(fmt.Sprintf("%s 需要重入队列", filename), zap.String("任务进程ID", id))
			return err
		}
		//if (n + 1) == (end - start) {
		//
		//}
		//fmt.Println("写入临时文件字节数:", n, "range请求字节数:", end-start)
		err = fres.Body.Close()
		if err != nil {
			zap.L().Error("关闭相应体失败", zap.Error(err))
			return err
		}
		if err != nil && err != io.EOF {
			zap.L().Error("临时文件保存异常", zap.Error(err))
			return err
		}

		downloadedSize += n
		//fmt.Printf("Downloaded %s: %d / %d bytes \n", f.fileName, downloadedSize, f.fileSize)
		elapsedTime := time.Since(startTime)
		downloadSpeed := float64(f.blockSize) / elapsedTime.Seconds() / 1024 / 1024
		// 更新下载进度的
		//upDataProcess(dlink, fmt.Sprintf("%.2f%%", (float32(downloadedSize)/float32(f.fileSize))*100), db)
		if upDataProcess(dlink, fmt.Sprintf("%.2fMB/s - %.0f%%", downloadSpeed, (float32(downloadedSize)/float32(f.fileSize))*100), db) == 1 {
			return errors.New("该任务需要删除")
		}
		// 下载完成后退出循环
		if downloadedSize == f.fileSize {
			upDataProcess(dlink, "下载完成", db)
			break
		}
	}

	// 将下载好的临时文件改名为目标文件名
	err = os.Rename(tmpFile.Name(), f.dir+f.fileName)
	if err != nil {
		zap.L().Error("重命名失败", zap.Error(err))
		return err
	}
	// 本地hash校验
	//hash := md5.New()
	//file, _ := os.Open(f.dir + f.fileName)
	//defer file.Close()
	//
	//if _, err := io.Copy(hash, file); err != nil {
	//	return errors.New("生成md5值失败")
	//}
	//
	//if hex.EncodeToString(hash.Sum(nil)) == f.md5 {
	//	slog.Info("hash验证通过", "进程任务ID", id, "msg", f.fileName)
	//	return nil
	//}
	//return errors.New("hash验证失败")
	return nil
}

func upDataProcess(link, status string, db *gorm.DB) int {
	var f DownDetail
	if db.Find(&f, "dlink = ?", link).Error != nil {
		panic("更新下载进度失败")
	}
	if f.ProcessStatus == 1 {
		db.Delete(&f, "dlink = ?", link)
		return 1
	}
	f.Status = status
	if db.Select("status").Save(&f).Error != nil {
		panic("更新下载进度失败")
	}
	return 0
}
