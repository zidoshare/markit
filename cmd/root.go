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

package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zidoshare/markit/utils/paths"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:   "markit",
		Short: "专为中文语境优化的 Markdown 命令行工具",
		Long: `专为中文语境优化的 Markdown 命令行工具，包含：

* 格式化
* 渲染

针对中文语境，它参考 中文文案排版指北 [https://github.com/sparanoid/chinese-copywriting-guidelines](https://github.com/sparanoid/chinese-copywriting-guidelines)

格式化处理项

* 中英文之间增加空格(排除 "ing"， 例如工作中，可能写为 "工作ing"，这种情况不会自动加空格)
* 专有名词使用正确的大小写(例如 自动将 'github' 转换为 'Github')
* 自动处理表格对齐（务必使用等宽字符）

可以提供整个项目文档的统一处理，特别是在处理 Gitbook 等项目时，能够在多人合作中实现格式统一。

默认使用 .markit.toml 文件配置整个仓库的处理方式风格
`,
	}
)

// Execute 执行命令行解析
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "默认配置文件(将从指定路径及上层路径递归查找.markit.toml,如果未能找到将尝试查询$HOME/.markit.toml)")

	rootCmd.AddCommand(formatCmd)
	rootCmd.AddCommand(renderCmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName(".markit")
		viper.SetConfigType("toml")
		path, err := os.Getwd()
		if err != nil {
			cobra.CheckErr(err)
		}
		path = resolveConfigPath(path)
		if path == "" {
			var err error
			path, err = homedir.Dir()
			if err != nil {
				cobra.CheckErr(err)
			}
		}
		viper.AddConfigPath(path)
	}
	err := viper.ReadInConfig()
	notFound := &viper.ConfigFileNotFoundError{}
	switch {
	case err != nil && !errors.As(err, notFound):
		cobra.CheckErr(err)
	case err != nil && errors.As(err, notFound):
		break
	default:
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func loadConfig(path string) {
	viper.SetConfigName(".markit")
	viper.SetConfigType("toml")
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		if !paths.IsDir(path) {
			path = filepath.Dir(path)
		}
		path = resolveConfigPath(path)
		if path == "" {
			var err error
			path, err = homedir.Dir()
			if err != nil {
				cobra.CheckErr(err)
			}
		}
		viper.AddConfigPath(path)
	}

}

// resolveConfigPath 获取当前可用配置文件路径
func resolveConfigPath(currentPath string) (p string) {
	for {
		p = filepath.Join(currentPath, ".markit.toml")
		if paths.FileExists(p) {
			return currentPath
		}
		if paths.DirExists(path.Join(currentPath, ".git")) || currentPath == "/" {
			return ""
		}
		currentPath = filepath.Dir(currentPath)
	}
}

// WalkMdFile 递归遍历目录中所有的markdown文件
func WalkMdFile(dir string, cb func(path string) error) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && Ignore(info) {
			return filepath.SkipDir
		} else if !info.IsDir() && Ignore(info) {
			return nil
		}
		if filepath.Ext(path) == ".md" {
			return cb(path)
		}
		return nil
	})
}

// Ignore 按照规则是否忽略该文件
func Ignore(info os.FileInfo) bool {
	if info.IsDir() {
		//以.开头的文件夹均忽略
		if strings.HasPrefix(info.Name(), ".") {
			return true
		}
		if info.Name() == "node_modules" {
			return true
		}
	} else {
		if filepath.Ext(info.Name()) != ".md" {
			return true
		}
	}

	return false
}
