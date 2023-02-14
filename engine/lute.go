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
package engine

import (
	"github.com/88250/lute"
)

// Formatter 格式化接口
type Formatter interface {
	Format(bytes []byte) []byte
}

// LuteFormatter 使用Lute实现的格式化
type LuteFormatter struct {
	engine *lute.Lute
}

// NewFormatter 新建格式化
func NewFormatter(options Options) *LuteFormatter {
	return &LuteFormatter{
		engine: lute.New(func(engine *lute.Lute) {
			engine.SetAutoSpace(options.AutoSpace)
			engine.SetFixTermTypo(options.FixTermTypo)
		}),
	}
}

// Format 格式化
func (formatter *LuteFormatter) Format(bytes []byte) []byte {
	return formatter.engine.Format("Markit", bytes)
}

// Renderer 渲染器
type Renderer interface {
	Render(bytes []byte) []byte
}

// LuteRenderer 使用Lute实现的渲染器
type LuteRenderer struct {
	engine *lute.Lute
}

// NewRender 新建渲染器
func NewRender(options Options) *LuteRenderer {
	return &LuteRenderer{
		engine: lute.New(func(engine *lute.Lute) {
		}),
	}
}

// Render 渲染html
func (renderer *LuteRenderer) Render(bytes []byte) []byte {
	return renderer.engine.Markdown("Markit", bytes)
}
