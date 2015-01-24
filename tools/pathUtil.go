// pathUtil
package tools

import (
	"os"
	"os/exec"
	"path"
)

// 获取当前执行程序的绝对路径，不包含文件名
func GetCurrentPath() string {
	file, _ := exec.LookPath(os.Args[0])
	return path.Dir(file)
}
