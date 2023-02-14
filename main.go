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
package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zidoshare/markit/engine"
	"github.com/zidoshare/markit/utils/paths"
)

var rootCmd = &cobra.Command{
	Use:   "markit [path]...",
	Short: "专为中文语境优化的 Markdown 命令行格式化工具",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			for _, p := range args {
				if paths.IsDir(p) {
					cobra.CheckErr(formatDir(p))
				} else {
					cobra.CheckErr(formatFile(p))
				}
			}
		} else {
			cobra.CheckErr(format(os.Stdin, os.Stdout))
		}
	},
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
		// 以.开头的文件夹均忽略
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

func format(r io.Reader, w io.Writer) error {
	content, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	content = engine.NewFormatter().Format(content)
	if _, err := w.Write(content); err != nil {
		return err
	}
	return nil
}

func formatFile(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	content = engine.NewFormatter().Format(content)
	if err := os.WriteFile(path, content, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func formatDir(dir string) error {
	return WalkMdFile(dir, formatFile)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
