/*
Copyright © 2021 zido <wuhongxu1208@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

//ResolveConfigPath 获取当前可用配置文件路径
func ResolveConfigPath(currentPath string) (p string) {
	for {
		p = filepath.Join(currentPath, ".markit.toml")
		if FileExists(p) {
			return currentPath
		}
		if currentPath == "/" {
			return ""
		}
		currentPath = filepath.Dir(currentPath)
	}
}

//IsDir 判断目录是否是一个文件夹
func IsDir(p string) bool {
	stat, err := os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		err = fmt.Errorf("文件状态获取错误:%s\n", err)
		if err != nil {
			fmt.Println(err)
		}
		return false
	}
	return stat.IsDir()
}

//FileExists 判断文件是否存在（若path是一个文件夹，也会返回false，应当使用DirExists方法）
func FileExists(path string) bool {
	stat, err := os.Stat(path)
	if err == nil {
		if stat.IsDir() {
			return false
		}
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

//DirExists 判断文件夹是否存在（若path是一个文件，也会返回false，应当使用FileExists方法）
func DirExists(path string) bool {
	stat, err := os.Stat(path)
	if err == nil {
		if stat.IsDir() {
			return true
		}
		return false
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

//Exists 判断所给路径文件/文件夹是否存在(返回true是存在)
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
