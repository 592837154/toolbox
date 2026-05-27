# Markdown To Word (Go + Pandoc)

一个用 Go 封装 Pandoc 的 Markdown 批量转 Word 工具。当前只有两个核心目录：`convert/` 负责转换，`clean/` 负责清理。

## 目录结构

```text
Markdown-To-Word/
├─ convert/
│  ├─ go.mod
│  ├─ main.go
│  ├─ README.md
│  ├─ convert.exe     # 转换程序，编译后生成
│  ├─ clean.exe       # 清理程序，从 clean/ 编译输出到这里
│  ├─ input/          # 放入待转换的 Markdown 文件，运行后自动创建
│  ├─ output/         # 输出 Word 文件，运行后自动创建
│  └─ reference.docx  # 可选，Pandoc Word 样式模板
├─ clean/
│  ├─ go.mod
│  └─ main.go
└─ README.md
```

## 转换工具

进入 `convert/` 目录编译：

```powershell
cd D:\learn\toolbox\Markdown-To-Word\convert
go build .
```

会生成：

```text
convert.exe
```

使用方式：

```powershell
.\convert.exe
```

第一次运行会自动创建：

- `convert/input/`: 放入 `.md` 或 `.markdown` 文件。
- `convert/output/`: 输出同名 `.docx` 文件。

支持递归扫描子目录。例如：

```text
convert/input/api/H5健康检查接口使用说明.md
```

会输出为：

```text
convert/output/api/H5健康检查接口使用说明.docx
```

## 清理工具

进入 `clean/` 目录编译：

```powershell
cd D:\learn\toolbox\Markdown-To-Word\clean
go build -o ..\convert\clean.exe .
```

会生成到 `convert/` 目录：

```text
D:\learn\toolbox\Markdown-To-Word\convert\clean.exe
```

说明：Go 原生的 `go build .` 只能把 exe 生成到当前目录；如果要让清理程序直接生成到 `convert/`，需要使用上面的 `-o ..\convert\clean.exe`。

运行清理程序：

```powershell
cd D:\learn\toolbox\Markdown-To-Word\convert
.\clean.exe
```

`clean.exe` 会清空：

- `Markdown-To-Word/convert/input/`
- `Markdown-To-Word/convert/output/`

它只删除这两个目录里面的内容，不删除 `input/`、`output/` 文件夹本身，也不会删除 `reference.docx`、源码或 exe。

## Pandoc 要求

需要先安装 Pandoc。本工具会按以下顺序查找 Pandoc：

1. `convert.exe` 同目录下的 `pandoc.exe`。
2. 环境变量 `PANDOC_PATH` 指向的路径。
3. 系统 `PATH` 中的 `pandoc`。
4. Windows 常见安装路径，例如 `C:\Program Files\Pandoc\pandoc.exe`。

检查 Pandoc：

```powershell
& "C:\Program Files\Pandoc\pandoc.exe" --version
```

## 自定义 Word 样式

如果希望输出 Word 使用固定字体、标题样式、代码块样式、页边距等，可以准备一个 `reference.docx` 放到 `convert.exe` 同目录。

生成 Pandoc 默认模板：

```powershell
& "C:\Program Files\Pandoc\pandoc.exe" -o reference.docx --print-default-data-file reference.docx
```

然后打开 `reference.docx` 修改 Word 样式，再重新运行 `convert.exe`。

## 说明

转换参数使用：

```text
--from gfm+tex_math_dollars+emoji
--to docx
--standalone
--toc
--toc-depth 3
```

Mermaid 等需要额外渲染的扩展暂未实现，建议先导出为图片后在 Markdown 中引用。
