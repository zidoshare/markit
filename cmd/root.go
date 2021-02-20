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
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
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

默认使用.markit.toml 文件配置整个仓库的处理方式风格
`,
}

//Execute 执行命令行解析
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
