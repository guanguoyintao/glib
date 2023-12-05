package uhttp

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	uoss "git.umu.work/AI/uglib/ustorage/oss"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

// FileType 表示文件的MIME类型枚举
type FileType int

const (
	FileTypeUnknown FileType = iota
	FileTypeText
	FileTypeImage
	FileTypeAudio
	FileTypeVideo
)

// String 返回file类型的字符串表示
func (t FileType) String() string {
	switch t {
	case FileTypeText:
		return "text"
	case FileTypeImage:
		return "image"
	case FileTypeAudio:
		return "audio"
	case FileTypeVideo:
		return "video"
	default:
		return "unknown"
	}
}

// FileInfo 包含文件大小、文件类型和文件名的结构体
type FileInfo struct {
	Size         int64     // 文件大小
	Type         FileType  // 文件类型
	Name         string    // 文件名
	ModifiedTime time.Time // 文件修改时间
}

func QuickDownload(url string, tmpDir string) (string, error) {
	if tmpDir == "" {
		tmpDir = os.TempDir()
	}
	u, err := uoss.CovertUrl2Inner(url)
	if err != nil {
		return "", errors.Wrapf(err, "trans url to inner failed")
	}
	arr := md5.Sum([]byte(strings.TrimSpace(u)))
	uniqueId := hex.EncodeToString(arr[:8])
	filepath := fmt.Sprintf("%s/%s-%d", tmpDir, uniqueId, time.Now().Unix())
	file, err := os.Create(filepath)
	if err != nil {
		return "", errors.Wrapf(err, "create file failed, path=%s", filepath)
	}
	defer file.Close()
	resp, err := http.Get(u)
	if err != nil {
		return "", errors.Wrapf(err, "http get content failed, url=%s", u)
	}
	defer resp.Body.Close()
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", errors.Wrapf(err, "copy http resp failed, url=%s,path=%s", u, filepath)
	}
	return filepath, nil
}

// CheckFileExistence 检查在线资源是否存在
func CheckFileExistence(urlString string) bool {
	resp, err := http.Head(urlString)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	// 检查 status code
	exists := resp.StatusCode == http.StatusOK

	return exists
}

// GetFileInfo 获取指定 URL 的文件信息
func GetFileInfo(fileURL string) (*FileInfo, error) {
	// 发送 HEAD 请求以获取文件信息，不下载实际文件内容
	resp, err := http.Head(fileURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 从HTTP头中获取MIME类型
	mimeType := resp.Header.Get("Content-Type")

	// 根据MIME类型设置文件类型
	var fileType FileType
	switch {
	case mimeType == "text/plain":
		fileType = FileTypeText
	case mimeType[:5] == "image":
		fileType = FileTypeImage
	case mimeType[:5] == "audio":
		fileType = FileTypeAudio
	case mimeType[:5] == "video":
		fileType = FileTypeVideo
	}

	// 从HTTP头中获取修改时间
	modifiedTimeStr := resp.Header.Get("Last-Modified")
	modifiedTime, _ := time.Parse(time.RFC1123, modifiedTimeStr)

	// 创建 HTTPFileInfo 结构体以保存文件信息
	fileInfo := &FileInfo{
		Size:         resp.ContentLength,
		Type:         fileType,
		Name:         path.Base(fileURL), // 从URL中提取文件名
		ModifiedTime: modifiedTime,
	}

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return fileInfo, fmt.Errorf("http get failed：status code=%d", resp.StatusCode)
	}

	return fileInfo, nil
}
