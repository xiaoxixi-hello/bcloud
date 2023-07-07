package floder

import (
	"bcloud/common/http"
	"bcloud/netdisk"
	"encoding/json"
	"fmt"
	"golang.org/x/exp/slog"
	"log"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func GetConfigDir() string {
	// 获取配置文件夹路径路径
	userInfo, err := user.Current()
	if err != nil {
		panic("Bcloud配置文件夹路径获取失败" + err.Error())
	}
	var homeDir = userInfo.HomeDir
	// 判断 homeDir/Bcloud 文件夹是否存在
	var gtDir = homeDir + "/Bcloud"
	_, err = os.Stat(gtDir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(gtDir, os.ModePerm)
			if err != nil {
				panic("Bcloud配置文件夹创建失败")
			}
		} else {
			panic("Bcloud配置文件夹不存在--" + err.Error())
		}
	}
	return gtDir
}

// CreateDir 创建远程文件夹
func CreateDir(p string, token string) {
	protocal := "https"
	host := "pan.baidu.com"
	router := "/rest/2.0/xpan/file?method=create&"
	uri := protocal + "://" + host + router

	//  设置url参数
	params := url.Values{}
	params.Set("access_token", token)

	uri += params.Encode()

	headers := map[string]string{
		"Host":         host,
		"Content-Type": "application/x-www-form-urlencoded",
	}

	// 设置body参数
	postBody := url.Values{}
	postBody.Add("isdir", "1")
	postBody.Add("path", p)
	body, _, err := http.DoHTTPRequest(uri, strings.NewReader(postBody.Encode()), headers)
	if err != nil {
		log.Println("HTTP响应失败", err)
		panic(err)
	}
	r := netdisk.Resp{}
	if err = json.Unmarshal([]byte(body), &r); err != nil {
		fmt.Println(err)
		panic(err)
	}
	if r.Errno != 0 {
		log.Println("创建文件夹失败", r)
		panic(err)
	}
	log.Println("创建文件夹成功", r)
}

func DownDirPath(f string, base string) (string, string, error) {
	var localPath string
	localPath = base + f
	if !filepath.IsAbs(f) {
		dir, err := os.Getwd()
		if err != nil {
			slog.Error("err", err)
			return "", "", err
		}
		localPath = filepath.Join(dir, "/", localPath)
	}
	// 创建文件夹
	var path, fileName string
	s, _ := os.Stat(path)
	if s != nil && s.IsDir() {
		path = localPath
	} else {
		pathList := strings.Split(localPath, "/")
		fileName = pathList[len(pathList)-1]
		path = localPath[:len(localPath)-len(fileName)]
	}
	err := os.MkdirAll(path, os.ModePerm)
	return path, fileName, err
}
