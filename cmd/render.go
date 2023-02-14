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
	"os"
	"path/filepath"
	"strings"

	"github.com/zidoshare/markit/engine"
	"github.com/zidoshare/markit/styles"
	"github.com/zidoshare/markit/utils/html"
	"github.com/zidoshare/markit/utils/paths"
	"github.com/zidoshare/markit/utils/strbytesconv"

	"github.com/spf13/cobra"
)

var (
	// 输出路径
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

	// renderCmd format 命令
	renderCmd = &cobra.Command{
		Use:   "render [file]",
		Short: "渲染markdown文档",
		Long: `渲染markdown文档，可指定目录或文件路径，例如：

 markit render .

将渲染当前目录下所有markdown文本

若不指定路径，默认从标准输入流获取内容直到EOF，并写出到标准输出流，例如：

 cat README.me | markit render
	`,
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			in := ""
			if len(args) != 0 {
				in = args[0]
			} else {
				in = os.Args[0]
			}
			in, err = filepath.Abs(in)
			if err != nil {
				cobra.CheckErr(err)
			}
			out, err = filepath.Abs(out)
			if err != nil {
				cobra.CheckErr(err)
			}
			if stand {
				if paths.IsDir(in) || paths.IsDir(out) {
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
					cobra.CheckErr(fmt.Errorf("css目录必须在输出目录内部"))
				}
			}
			if len(args) != 0 {
				if paths.IsDir(in) {
					err = renderDir(in, out)
				} else {
					err = renderFile(in, out)
				}
			} else {
				err = render(os.Stdin, os.Stdout)
			}
			cobra.CheckErr(err)
		},
	}
)

func init() {
	renderCmd.Flags().StringVarP(&out, "out", "o", "", "当输入为文件时，目标路径必须为文件路径。当输入为文件夹时，目标路径必须为文件夹路径")
	renderCmd.Flags().BoolVarP(&styled, "style", "s", false, "是否添加css样式")
	renderCmd.Flags().BoolVarP(&stand, "stand", "S", false, "是否使用独立样式，css和js文件将单独进行打包")
	renderCmd.Flags().BoolVarP(&body, "body", "b", false, "生成完整的html，为渲染结果添加<html><title><body>等标签")
	renderCmd.Flags().BoolVarP(&single, "single", "", false, "是否生成渲染为单页，使用局部加载，页面不发生跳转")
	renderCmd.Flags().BoolVarP(&graph, "graph", "g", false, "是否绘制关系图")
	renderCmd.Flags().StringVarP(&cssDir, "css-dir", "", "css", "指定css路径")
}

func render(r io.Reader, w io.Writer) error {
	content, err := ioutil.ReadAll(r)
	if err != nil {
		return fmt.Errorf("读取文件内容出错：%w", err)
	}
	content = engine.NewRender().Render(content)
	if _, err := w.Write(content); err != nil {
		return fmt.Errorf("写入文件时出错：%w", err)
	}
	return nil
}

// 渲染文件
func renderFile(mdPath, outPath string) error {
	arr := strings.Split(mdPath, string(filepath.Separator))
	filename := arr[len(arr)-1]
	filename = strings.TrimSuffix(filename, filepath.Ext(filename))
	if paths.IsDir(outPath) {
		outPath = filepath.Join(outPath, filename+".html")
	}
	content, err := ioutil.ReadFile(mdPath)
	if err != nil {
		return fmt.Errorf("读取文件内容出错：%w", err)
	}
	content = engine.NewRender().Render(content)
	if !body {
		if err := ioutil.WriteFile(outPath, content, os.ModePerm); err != nil {
			return fmt.Errorf("写入文件 %s 时出错：%w", outPath, err)
		}
		return nil
	}
	document := html.NewElement("html")
	head := html.NewElement("head").In(document)
	html.NewElement("title").Text(strbytesconv.ToBytes(filename)).In(head)
	if styled {
		style, styleHead := styles.Get("github")
		head.Text(strbytesconv.ToBytes(styleHead))
		//非独立样式会将样式放在head中
		if !stand {
			html.NewElement("style").Attr("type", "text/css").Text(strbytesconv.ToBytes(style)).In(head)
		} else {
			//创建css文件
			var cssPath = filepath.Join(cssDir, "github.css")
			if !paths.FileExists(cssPath) {
				if !paths.DirExists(filepath.Dir(cssPath)) {
					os.MkdirAll(filepath.Dir(cssPath), os.ModePerm)
				}
				ioutil.WriteFile(cssPath, strbytesconv.ToBytes(style), os.ModePerm)
			}
			//计算相对路径
			relCSSPath, err := filepath.Rel(filepath.Dir(outPath), cssPath)
			if err != nil {
				return fmt.Errorf("无法计算 css 样式文件放置的相对路径: %w", err)
			}
			html.NewElement("link").Attr("rel", "stylesheet").Attr("href", relCSSPath).In(head)
		}
	}
	var body *html.Element
	if styled {
		body = html.NewElement("body").Attr("class", "markdown-body").Append(head)
	} else {
		body = html.NewElement("body").Append(head)
	}
	body.Text(content)
	if err := ioutil.WriteFile(outPath, html.WriteElement(document), os.ModePerm); err != nil {
		return fmt.Errorf("写入文件 %s 时出错：%w", outPath, err)
	}
	return nil
}

func renderDir(dirIn, dirOut string) error {
	return WalkMdFile(dirIn, func(mdPath string) error {
		relativePath, err := filepath.Rel(dirIn, filepath.Dir(mdPath))
		if err != nil {
			return fmt.Errorf("遍历文件时出错，文件路径：%s，错误原因：%w", mdPath, err)
		}
		resultPath := filepath.Join(dirOut, relativePath)
		if !paths.DirExists(resultPath) {
			return os.MkdirAll(resultPath, os.ModePerm)
		}
		arr := strings.Split(mdPath, string(filepath.Separator))
		filename := arr[len(arr)-1]
		filename = strings.TrimSuffix(filename, filepath.Ext(filename))
		resultPath = filepath.Join(resultPath, filename) + ".html"
		renderFile(mdPath, resultPath)
		return nil
	})
}
