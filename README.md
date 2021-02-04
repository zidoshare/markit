# Markit

专为中文语境优化的 Markdown 命令行工具，包含：

* 解析
* 渲染
* 导出
* 格式化

针对中文语境，它参考 中文文案排版指北 [https://github.com/sparanoid/chinese-copywriting-guidelines](https://github.com/sparanoid/chinese-copywriting-guidelines)

格式化处理项

* 中英文之间增加空格(排除 "ing"， 例如工作中，可能写为 "工作ing"，这种情况不会自动加空格)
* 专有名词使用正确的大小写(例如 自动将 `github` 转换为 `Github`)

可以提供整个项目文档的统一处理，特别是在处理 Gitbook 等项目时，能够在多人合作中实现格式统一。

默认使用.markit.toml 文件配置整个仓库的处理方式风格

```
Usage:
  markit [command]

Available Commands:
  format      格式化markdown文档
  help        Help about any command

Flags:
  -h, --help   help for markit

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
