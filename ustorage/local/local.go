package ulocal

import (
	"bufio"
	"fmt"
	"git.umu.work/AI/uglib/uerrors"
	"os"
	"path/filepath"
	"strings"
)

// IsPathExists 判断文件夹是否存在
func IsPathExists(dir string) (os.FileInfo, bool, error) {
	fileInfo, err := os.Stat(dir)
	if err == nil {
		return fileInfo, true, nil
	}
	if os.IsNotExist(err) {
		return nil, false, nil
	}

	return nil, false, err
}

// GetFd 获取文件描述符
func GetFd(filePath string, flag int) (fd *os.File, err error) {
	dir := filepath.Dir(filePath)
	_, isExist, err := IsPathExists(dir)
	if err != nil {
		return nil, err
	}
	if !isExist {
		// 递归创建文件夹
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}
	// 检查文件是否存在，多一次判断不影响
	_, isFileExist, err := IsPathExists(filePath)
	if err != nil {
		return nil, err
	}
	if !isFileExist {
		tmpFd, err := os.Create(filePath)
		if err != nil {
			return nil, err
		}
		err = tmpFd.Close()
		if err != nil {
			return nil, err
		}
	}
	fd, err = os.OpenFile(filePath, flag, 0666)
	if err != nil {
		return nil, err
	}

	return fd, nil
}

func WriteHead(head []byte, filePath, newFilePath string) error {
	file, err := os.OpenFile(filePath, os.O_RDWR, 0544)
	if err != nil {
		return err
	}
	reader := bufio.NewReader(file)
	var tempFilePath string
	if filePath != newFilePath {
		tempFilePath = newFilePath
	} else {
		tempFilePath = fmt.Sprintf("%s.tmp", filePath)
	}
	// 新建临时文件
	tempFile, err := os.OpenFile(tempFilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(tempFile)
	_ = writer.Flush()
	// head内容写入临时文件
	_, _ = writer.Write(head)
	_ = writer.Flush()
	// 把源文件的内容写入临时文件
	// 按块读取
	buf := make([]byte, 1<<(1*10)) // 1k
	for {
		n, err := reader.Read(buf) // 依次读一行
		if err != nil {
			if uerrors.IsCloseError(err) {
				break
			}
			return err
		}
		_, err = writer.Write(buf[:n])
		if err != nil {
			return err
		}
	}
	_ = writer.Flush()
	// 释放文件占用
	err = file.Close()
	if err != nil {
		return err
	}
	err = tempFile.Close()
	if err != nil {
		return err
	}
	// 根据是否生成新文件判断是否需要rename
	if filePath == newFilePath {
		err = os.Rename(tempFilePath, filePath)
		if err != nil {
			return err
		}
	}

	return nil
}

func RemoveFileExtension(filename string) string {
	// 使用 filepath 包获取文件名的基础部分
	base := filepath.Base(filename)

	// 使用 strings 包的 TrimSuffix 函数去除后缀
	name := strings.TrimSuffix(base, filepath.Ext(base))

	return name
}

// RemoveFile 安全删除文件，不报错
func RemoveFile(filePath string) {
	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return // 文件不存在，无需删除
	}
	// 删除文件，忽略错误
	_ = os.Remove(filePath)
}
