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
	"bytes"
)

type Token string

const (
	HTML Token = "html"
	HEAD Token = "head"
	BODY Token = "body"
)

//Element Html节点
type Element struct {
	ElementBuilder
	Parent *Element
	Writer *bytes.Buffer
}

type ElementBuilder struct {
	Name       string
	Attributes map[string]string
	Text       string
	Children   []*Element
}

func NewElement(name string) *ElementBuilder {
	return &ElementBuilder{
		Name: name,
	}
}

func (e *Element) open() *Element {
	buffer := e.Writer
	buffer.WriteString("<")
	buffer.WriteString(e.Name)
	for key, val := range e.Attributes {
		buffer.WriteString(" ")
		buffer.WriteString(key)
		buffer.WriteString("=")
		buffer.WriteString("\"")
		buffer.WriteString(val)
		buffer.WriteString("\"")
	}
	buffer.WriteString(">")
	return e
}

func (e *Element) body() *Element {
	buffer := e.Writer
	buffer.WriteString(e.Text)
	return e
}

func (e *Element) cloze() *Element {
	buffer := e.Writer
	buffer.WriteString("</")
	buffer.WriteString(e.Name)
	buffer.WriteString(">")
	return e
}

func (e *Element) Bytes() []byte {
	offset := e.Writer.Len()
	e.open()
	e.body()
	e.cloze()
	end := e.Writer.Len()
	return e.Writer.Bytes()[offset:end]
}

func (e *Element) In(parent *Element) *Element {
	parent.Children = append(parent.Children, e)
	e.Parent = parent
	return e
}

func (e *Element) Append(last *Element) *Element {
	return e.In(last.Parent)
}

type Render interface {
	Write(token Token)
}

type StyleRender struct {
	//主题
	Theme string
	//是否独立
	Stand bool
}

func (ret *StyleRender) Write(token Token) {
}
