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
	"io"
	"io/ioutil"
	"markit/engine"
	"markit/utils"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// formatCmd format 命令
var formatCmd = &cobra.Command{
	Use:   "format",
	Short: "格式化markdown文档",
	Long: `格式化markdown文档，可指定目录或文件路径，例如：

 markit format .

将格式化当前目录下所有markdown文本

若不指定路径，默认从标准输入流获取内容直到EOF，并写出到标准输出流，例如：

 cat README.me | markit format
	`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 0 {
			p, err := filepath.Abs(args[0])
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			loadConfig(p)
			if utils.IsDir(p) {
				formatDir(p)
			} else {
				formatFile(p)
			}
		} else {
			p, err := filepath.Abs(os.Args[0])
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			loadConfig(p)
			format(os.Stdin, os.Stdout)
		}
	},
}

func format(r io.Reader, w io.Writer) {
	content, err := ioutil.ReadAll(r)
	if err != nil {
		err = fmt.Errorf("读取文件内容出错:%s", err)
		if err != nil {
			fmt.Println(err)
		}
		os.Exit(1)
	}
	content = engine.NewFormatter(engine.NewOptions()).Format(content)
	w.Write(content)
}
func formatFile(path string) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("读取文件内容出错:%s\n", err)
		if err != nil {
			fmt.Println(err)
		}
		os.Exit(1)
	}
	content = engine.NewFormatter(engine.NewOptions()).Format(content)
	ioutil.WriteFile(path, content, os.ModePerm)
}
func formatDir(dir string) {
	WalkMdFile(dir, formatFile)
}

func init() {
	processCmd(formatCmd)
}
