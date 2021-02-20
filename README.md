# Markit

> 专为中文语境优化的 Markdown 命令行工具

# 关于

该项目为从命令行使用markdown提供了一个简单的配套解决方案。 

包含：

* 格式化
* 渲染

针对中文语境，它参考 中文文案排版指北 [https://github.com/sparanoid/chinese-copywriting-guidelines](https://github.com/sparanoid/chinese-copywriting-guidelines)

格式化处理项

* 中英文之间增加空格(排除 "ing"， 例如工作中，可能写为 "工作ing"，这种情况不会自动加空格)
* 专有名词使用正确的大小写(例如 自动将 `github` 转换为 `Github`)
* 自动处理表格对齐（务必使用等宽字符）

可以提供整个项目文档的统一处理，特别是在处理 Gitbook 等项目时，能够在多人合作中实现格式统一。

默认使用.markit.toml 文件配置整个仓库的处理方式风格

```
Usage:
  markit [command]

Available Commands:
  format      格式化markdown文档
  help        Help about any command
  render      渲染markdown文档

Flags:
  -c, --config string   默认配置文件(将从指定路径及上层路径递归查找.markit.toml,如果未能找到将尝试查询$HOME/.markit.toml)
  -h, --help            help for markit

Use "markit [command] --help" for more information about a command.
```

# 安装

## Go

```
go get -u github.com/zidoshare/markit
```

# 🙏 鸣谢

站在巨人们的肩膀上:

* [中文文案排版指北](https://github.com/sparanoid/chinese-copywriting-guidelines)：统一中文文案、排版的相关用法，降低团队成员之间的沟通成本，增强网站气质
* [lute](https://github.com/8825/lute)：一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript。
