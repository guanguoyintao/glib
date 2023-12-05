package uhash

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/spaolacci/murmur3"
	"io"
	"os"
	"strings"
)

func HashMD532(s string) string {
	sum := md5.Sum([]byte(strings.TrimSpace(s)))
	return hex.EncodeToString(sum[:])
}

func HashMurmurHash340(s string) (string, error) {
	h := murmur3.New128()
	_, err := h.Write([]byte(s))
	if err != nil {
		return "", err
	}
	sum := h.Sum([]byte(s))

	return hex.EncodeToString(sum[:]), nil
}

func CalcFileSHA256(ctx context.Context, filePath string) (string, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 创建 SHA256 哈希对象
	hash := sha256.New()

	// 读取文件内容并计算哈希值
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	// 获取计算得到的哈希值
	hashValue := hash.Sum(nil)

	// 将哈希值转换为十六进制字符串
	hashString := fmt.Sprintf("%x", hashValue)

	return hashString, nil
}

func CalcContentSHA256(ctx context.Context, content io.Reader) (string, error) {
	// 创建 SHA256 哈希对象
	hash := sha256.New()

	// 读取文件内容并计算哈希值
	if _, err := io.Copy(hash, content); err != nil {
		return "", err
	}

	// 获取计算得到的哈希值
	hashValue := hash.Sum(nil)

	// 将哈希值转换为十六进制字符串
	hashString := fmt.Sprintf("%x", hashValue)

	return hashString, nil
}

func CalcFileSHA1(ctx context.Context, filePath string) (string, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 创建 SHA1 哈希对象
	hash := sha1.New()

	// 读取文件内容并计算哈希值
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	// 获取计算得到的哈希值
	hashValue := hash.Sum(nil)

	// 将哈希值转换为十六进制字符串
	hashString := fmt.Sprintf("%x", hashValue)

	return hashString, nil
}

func CalcContentSHA1(ctx context.Context, content io.Reader) (string, error) {
	// 创建 SHA1 哈希对象
	hash := sha1.New()

	// 读取文件内容并计算哈希值
	if _, err := io.Copy(hash, content); err != nil {
		return "", err
	}

	// 获取计算得到的哈希值
	hashValue := hash.Sum(nil)

	// 将哈希值转换为十六进制字符串
	hashString := fmt.Sprintf("%x", hashValue)

	return hashString, nil
}
