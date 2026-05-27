# Markdown To Word Convert

把 `input/` 文件夹中的 Markdown 批量转换为 `output/` 中的 Word `.docx` 文件。

## 环境要求

需要先安装 Pandoc。本工具会按以下顺序查找 Pandoc：

1. `convert.exe` 同目录下的 `pandoc.exe`。
2. 环境变量 `PANDOC_PATH` 指向的路径。
3. 系统 `PATH` 中的 `pandoc`。
4. Windows 常见安装路径，例如 `C:\Program Files\Pandoc\pandoc.exe`。

## 编译

在当前目录运行：

```powershell
go build .
```

会生成：

```text
convert.exe
```

## 使用

第一次运行会自动创建：

- `input/`: 放入 `.md` 或 `.markdown` 文件。
- `output/`: 输出同名 `.docx` 文件。

运行：

```powershell
.\convert.exe
```

支持子目录。例如：

```text
input/api/H5健康检查接口使用说明.md
```

会输出为：

```text
output/api/H5健康检查接口使用说明.docx
```

## 自定义 Word 样式

如果希望固定字体、标题样式、代码块样式、页边距等，可以把 `reference.docx` 放在 `convert.exe` 同目录。

生成 Pandoc 默认模板：

```powershell
& "C:\Program Files\Pandoc\pandoc.exe" -o reference.docx --print-default-data-file reference.docx
```

然后打开 `reference.docx` 修改 Word 样式，再重新运行转换程序。
