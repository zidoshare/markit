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
package utils

import "unsafe"

//string 和 []byte 互转都涉及底层数据复制；通过 unsafe 强制转换绕过复制提高性能。
//string 类型的底层数组可能放在常量区（字面量初始化）也可能动态分配；无论哪一种，当以 string 类型出现时底层数据都是不可修改的（避免影响其他引用），string 的修改实际上是指向重新生成的底层数组。
//绕过类型检查强行转换，绕过了底层数据复制，提高性能同时也失去了检查和复制的保护，需要调用方自行确认不会出错（不能直接修改，必须确定是否）。

// BytesToStr 快速转换 []byte 为 string。
func BytesToStr(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}

// StrToBytes 快速转换 string 为 []byte。
func StrToBytes(str string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&str))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}
