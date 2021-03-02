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
	out string
	//是否生成完整html
	body bool
	//是否添加css样式
	styled bool
	//是否附加独立的lib，这会将css样式放在css/文件夹下，js代码放置到js/文件夹下
	stand bool
	//是否使用单页
	single bool
	//是否生成关系图
	graph bool

	//css位置
	cssDir string
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
		in := ""
		if len(args) != 0 {
			in = args[0]
		} else {
			in = os.Args[0]
		}
		in, err := filepath.Abs(in)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		loadConfig(in)
		out, err = filepath.Abs(out)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if stand {
			if utils.IsDir(in) || utils.IsDir(out) {
				cssDir = filepath.Join(out, cssDir)
			} else {
				cssDir = filepath.Join(filepath.Dir(out), cssDir)
			}
			cssDir, err := filepath.Abs(filepath.Clean(cssDir))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			//cssDir必须在out内
			if !strings.HasPrefix(cssDir, out) {
				fmt.Println("css目录必须在输出目录内部")
				os.Exit(1)
			}
		}
		if len(args) != 0 {
			if utils.IsDir(in) {
				renderDir(in, out)
			} else {
				renderFile(in, out)
			}
		} else {
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

//渲染文件
func renderFile(mdPath, outPath string) {
	arr := strings.Split(mdPath, string(filepath.Separator))
	filename := arr[len(arr)-1]
	filename = strings.TrimSuffix(filename, filepath.Ext(filename))
	if utils.IsDir(outPath) {
		outPath = filepath.Join(outPath, filename+".html")
	}
	content, err := ioutil.ReadFile(mdPath)
	if err != nil {
		fmt.Printf("读取文件内容出错:%s\n", err)
		os.Exit(1)
	}
	content = engine.NewRender(engine.NewOptions()).Render(content)
	if !body {
		ioutil.WriteFile(outPath, content, os.ModePerm)
		return
	}
	var buffer bytes.Buffer
	buffer.WriteString("<html>\n<head>\n<title>")

	buffer.WriteString(filename)
	buffer.WriteString("</title>\n")
	if styled {
		style, head := styles.Get("github")
		buffer.WriteString(head)
		//非独立样式会将样式放在head中
		if !stand {
			buffer.WriteString("<style type=\"text/css\">\n")
			buffer.WriteString(style)
			buffer.WriteString("</style>")
		} else {
			//创建css文件
			var cssPath = filepath.Join(cssDir, "github.css")
			if !utils.FileExists(cssPath) {
				if !utils.DirExists(filepath.Dir(cssPath)) {
					os.MkdirAll(filepath.Dir(cssPath), os.ModePerm)
				}
				ioutil.WriteFile(cssPath, utils.StrToBytes(style), os.ModePerm)
			}
			//计算相对路径
			relCssPath, err := filepath.Rel(filepath.Dir(outPath), cssPath)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			buffer.WriteString("<link rel=\"stylesheet\" href=\"" + relCssPath + "\"/>")
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
	ioutil.WriteFile(outPath, buffer.Bytes(), os.ModePerm)
}

func renderDir(dirIn, dirOut string) {
	WalkMdFile(dirIn, func(mdPath string) {
		relativePath, err := filepath.Rel(dirIn, filepath.Dir(mdPath))
		if err != nil {
			fmt.Printf("渲染文件出错：%s", err)
			os.Exit(1)
		}
		resultPath := filepath.Join(dirOut, relativePath)
		if !utils.DirExists(resultPath) {
			os.MkdirAll(resultPath, os.ModePerm)
		}
		arr := strings.Split(mdPath, string(filepath.Separator))
		filename := arr[len(arr)-1]
		filename = strings.TrimSuffix(filename, filepath.Ext(filename))
		resultPath = filepath.Join(resultPath, filename) + ".html"
		renderFile(mdPath, resultPath)
	})
}

func init() {
	processCmd(renderCmd)
	renderCmd.Flags().StringVarP(&out, "out", "o", "", "当输入为文件时，目标路径必须为文件路径。当输入为文件夹时，目标路径必须为文件夹路径")
	renderCmd.Flags().BoolVarP(&styled, "style", "s", false, "是否添加css样式")
	renderCmd.Flags().BoolVarP(&stand, "stand", "S", false, "是否使用独立样式，css和js文件将单独进行打包")
	renderCmd.Flags().BoolVarP(&body, "body", "b", false, "生成完整的html，为渲染结果添加<html><title><body>等标签")
	renderCmd.Flags().BoolVarP(&single, "single", "", false, "是否生成渲染为单页，使用局部加载，页面不发生跳转")
	renderCmd.Flags().BoolVarP(&graph, "graph", "g", false, "是否绘制关系图")
	renderCmd.Flags().StringVarP(&cssDir, "css-dir", "", "css", "指定css路径")
}
