// fileUtil
package tools

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

var _ = fmt.Println

// 返回文件大小，单位字节
// 支持目录大小
func GetFileSize(filePaht string) (int64, error) {
	f, e := os.Stat(filePaht)
	if e != nil {
		return 0, e
	}
	return f.Size(), nil
}

// 返回文件最后修改时间
func GetModifyTime(filePaht string) (time.Time, error) {
	f, e := os.Stat(filePaht)
	if e != nil {
		return time.Now(), e
	}
	return f.ModTime(), nil
}

// 删除文件
func DeleteFile(filePath string) error {
	return os.Remove(filePath)
}

// 重命名文件
func Rename(srcPath string, toPath string) error {
	return os.Rename(srcPath, toPath)
}

// 检查是否是一个文件，如果是目录或者不存在返回false
func IsFile(filePath string) bool {
	f, e := os.Stat(filePath)
	if e != nil {
		return false
	}
	return !f.IsDir()
}

// 文件或目录是否存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// 创建文件，如果目录不存在则尝试创建目录，如果存在直接返回
func CreateFile(path string, fileName string) (string, error) {
	var fullPath string
	if strings.LastIndex(path, "/") == len(path)-1 || strings.LastIndex(path, "\\") == len(path)-1 {
		fullPath = path + fileName
	} else {
		if strings.Index(path, "/") == -1 {
			fullPath = path + "\\" + fileName
		} else {
			fullPath = path + "/" + fileName
		}
	}
	if IsExist(fullPath) {
		return fullPath, nil
	}

	if err := os.MkdirAll(path, 0777); err != nil {
		return "", err
	}
	if _, err := os.Create(path + "/" + fileName); err != nil {
		return "", err
	}
	return fullPath, nil
}

// 读取配置文件，一行一条配置，注释行用#开头，键值之间用=分割
func ReadConfig(filePath string) (map[string]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	buff := bufio.NewReader(file)
	m := make(map[string]string)
	for {
		line, err := buff.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		if strings.Count(line, "#") > 0 {
			// 被注释的行，忽略掉
		} else {
			if (strings.Count(line, "=")) == 1 {
				q := strings.LastIndex(line, "=")
				if q != -1 {
					m[strings.TrimSpace(line[0:q-1])] = strings.TrimSpace(line[q+1:])
				}
			}
		}
	}
	return m, nil
}
