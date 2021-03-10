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
package html

import "bytes"

//Element Html节点
type Element struct {
	name       string
	attributes map[string]string
	text       []byte
	children   []*Element
	parent     *Element
}

//NewElement 新建节点
func NewElement(name string) *Element {
	return &Element{
		name:       name,
		attributes: make(map[string]string),
	}
}

//NewTextElement 文本节点
func NewTextElement(text []byte) *Element {
	return &Element{
		name: "",
		text: text,
	}
}

//Name 设置节点名
func (e *Element) Name(name string) *Element {
	e.name = name
	return e
}

//Attr 设置节点属性
func (e *Element) Attr(key, value string) *Element {
	e.attributes[key] = value
	return e
}

//Text 添加文本
func (e *Element) Text(text []byte) *Element {
	e.children = append(e.children, NewTextElement(text))
	return e
}

//WriteElement 获取字节数组
func WriteElement(e *Element) []byte {
	if e.name == "" {
		var buffer bytes.Buffer
		buffer.Write(e.text)
		buffer.WriteString("\n")
		return buffer.Bytes()
	}
	var buffer bytes.Buffer
	buffer.WriteString("<")
	buffer.WriteString(e.name)
	for key, val := range e.attributes {
		buffer.WriteString(" ")
		buffer.WriteString(key)
		buffer.WriteString("=")
		buffer.WriteString("\"")
		buffer.WriteString(val)
		buffer.WriteString("\"")
	}
	if len(e.children) == 0 {
		//闭合标签
		buffer.WriteString("/>\n")
		return buffer.Bytes()
	}
	buffer.WriteString(">\n")
	for _, child := range e.children {
		buffer.Write(WriteElement(child))
	}
	buffer.WriteString("</")
	buffer.WriteString(e.name)
	buffer.WriteString(">\n")
	return buffer.Bytes()
}

//In 将节点放置于目标子节点中
func (e *Element) In(parent *Element) *Element {
	parent.children = append(parent.children, e)
	e.parent = parent
	return e
}

//Append 将节点放置于目标相邻节点中
func (e *Element) Append(last *Element) *Element {
	return e.In(last.parent)
}
