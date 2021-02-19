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
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"markit/engine"
	"markit/styles"
	"markit/utils"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	targetPath string
	//是否生成完整html
	body bool
	//是否添加css样式
	styled bool
	//是否附加独立的lib，这会将css样式放在css/文件夹下，js代码放置到js/文件夹下
	stand bool

	//css位置
	cssDir = "css"
)

// renderCmd format 命令
var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "渲染markdown文档",
	Long: `渲染markdown文档，可指定目录或文件路径，例如：

 markit render .

将渲染当前目录下所有markdown文本

若不指定路径，默认从标准输入流获取内容直到EOF，并写出到标准输出流，例如：

 cat README.me | markit render
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
				renderDir(p)
			} else {
				targetPath, err = filepath.Abs(targetPath)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				renderFile(p, targetPath)
			}
		} else {
			p, err := filepath.Abs(os.Args[0])
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			loadConfig(p)
			render(os.Stdin, os.Stdout)
		}
	},
}

func render(r io.Reader, w io.Writer) {
	content, err := ioutil.ReadAll(r)
	if err != nil {
		err = fmt.Errorf("读取文件内容出错:%s", err)
		if err != nil {
			fmt.Println(err)
		}
		os.Exit(1)
	}
	content = engine.NewRender(engine.NewOptions()).Render(content)
	w.Write(content)
}

func renderFile(p string, targetPath string) {
	resultPath := targetPath
	var filename string
	if utils.IsDir(resultPath) {
		arr := strings.Split(p, string(filepath.Separator))
		filename = arr[len(arr)-1]
		filename = strings.TrimSuffix(filename, filepath.Ext(filename))
		resultPath = filepath.Join(resultPath, filename) + ".html"
	}
	content, err := ioutil.ReadFile(p)
	if err != nil {
		fmt.Printf("读取文件内容出错:%s\n", err)
		os.Exit(1)
	}
	content = engine.NewRender(engine.NewOptions()).Render(content)
	if !body {
		ioutil.WriteFile(resultPath, content, os.ModePerm)
		return
	}
	var buffer bytes.Buffer
	buffer.WriteString("<html>\n<head>\n<title>")
	buffer.WriteString(filename)
	buffer.WriteString("</title>\n")
	if styled {
		buffer.WriteString(styles.GetExtraHead("github"))
		//非独立样式会将样式放在head中
		if !stand {
			buffer.WriteString("<style type=\"text/css\">\n")
			buffer.WriteString(styles.Get("github"))
			buffer.WriteString("</style>")
		} else {
			cssPath := filepath.Join(cssDir, "github.css")
			buffer.WriteString("<link rel=\"stylesheet\" href=\"" + cssPath + "\"/>")
			if !utils.FileExists(cssPath) {
				if !utils.DirExists(cssDir) {
					os.MkdirAll(cssDir, os.ModePerm)
				}
				ioutil.WriteFile(cssPath, utils.StrToBytes(styles.Get("github")), os.ModePerm)
			}
		}
	}
	buffer.WriteString("</head>\n")
	if styled {
		buffer.WriteString("<body class=\"markdown-body\">")
	} else {
		buffer.WriteString("<body>")
	}
	buffer.Write(content)
	buffer.WriteString("</body>\n</html>")
	ioutil.WriteFile(resultPath, buffer.Bytes(), os.ModePerm)
}

func renderDir(p string) {
	WalkMdFile(p, func(path string) {
		relativePath, err := filepath.Rel(p, filepath.Dir(path))
		if err != nil {
			fmt.Printf("渲染文件出错：%s", err)
			os.Exit(1)
		}
		resultPath := filepath.Join(targetPath, relativePath)
		if !utils.DirExists(resultPath) {
			os.MkdirAll(resultPath, os.ModePerm)
		}
		renderFile(path, resultPath)
	})
}

func init() {
	processCmd(renderCmd)
	renderCmd.Flags().StringVarP(&targetPath, "out", "o", "", "当输入为文件时，目标路径必须为文件路径。当输入为文件夹时，目标路径必须为文件夹路径")
	renderCmd.Flags().BoolVarP(&styled, "style", "s", false, "是否添加css样式")
	renderCmd.Flags().BoolVarP(&stand, "stand", "S", false, "是否使用独立样式")
	renderCmd.Flags().BoolVarP(&body, "body", "b", false, "生成完整的html，这会为渲染结果添加<html><title><body>等标签")
}
