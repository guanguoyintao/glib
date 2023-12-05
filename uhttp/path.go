package uhttp

import (
	"net/url"
	"path"
	"strings"
)

func ExtractFileNameFromURL(urlString string) (string, string, error) {
	// 解析URL
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return "", "", err
	}
	// 获取文件名部分
	fileName := path.Base(parsedURL.Path)
	fileDir := path.Dir(parsedURL.Path)

	// 去除可能的查询参数
	fileName = strings.Split(fileName, "?")[0]

	return fileName, fileDir, nil
}

func IsURL(p string) bool {
	// 解析 URL
	u, err := url.Parse(p)
	if err != nil {
		return false // 无法解析为有效的 URL，可能是本地路径
	}
	// 检查 URL 的 Scheme 是否为 http 或 https
	scheme := strings.ToLower(u.Scheme)
	if scheme == "http" || scheme == "https" {
		return true // 是在线的 HTTP 地址
	}

	return false // 不是在线的 HTTP 地址，可能是本地路径
}
