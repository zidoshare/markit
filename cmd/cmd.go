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
	"fmt"
	"markit/utils"
	"os"
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "默认配置文件(将从指定路径及上层路径递归查找.markit.toml,如果未能找到将尝试查询$HOME/.markit.toml)")
}

func processCmd(cmd *cobra.Command) {
	rootCmd.AddCommand(cmd)
}

func loadConfig(path string) {
	viper.SetConfigName(".markit")
	viper.SetConfigType("toml")
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		if !utils.IsDir(path) {
			path = filepath.Dir(path)
		}
		path = utils.ResolveConfigPath(path)
		if path == "" {
			var err error
			path, err = homedir.Dir()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		viper.AddConfigPath(path)
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return
		}
		fmt.Println(err)
		os.Exit(1)
	}
}

//WalkMdFile 递归遍历目录中所有的markdown文件
func WalkMdFile(dir string, cb func(path string)) {
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && Ignore(info) {
			return filepath.SkipDir
		} else if !info.IsDir() && Ignore(info) {
			return nil
		}
		if filepath.Ext(path) == ".md" {
			cb(path)
		}
		return nil
	})
}

//Ignore 按照规则是否忽略该文件
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
