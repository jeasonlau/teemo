package main

import (
	"os"
	"path/filepath"
)

// 获取up和down的路径
func GetImgPath() (upPath, downPath string) {
	p, _ := os.Executable()
	path, _ := filepath.Abs(p)
	dir := filepath.Dir(path) + "/ipgwTool/teemo"
	return dir + "/img/up.png", dir + "/img/down.png"
}
