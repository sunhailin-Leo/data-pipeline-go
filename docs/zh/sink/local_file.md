LocalFile
====

## 概述

`LocalFileSinkConfig` 结构体用于配置本地文件 Sink，包括与文件存储相关的所有设置，以便将数据保存到本地文件中。
## LocalFileSinkConfig 结构体

```golang
type LocalFileSinkConfig struct {
	FileName       string `json:"file_name"`
	FileFormatType string `json:"file_format_type"` // file format type: text, csv
	RowDelimiter   string `json:"row_delimiter"`    // only file_format_type = text is affect
}
```
## 字段说明

| 字段名                | 类型       | 描述                                       |
|--------------------|----------|------------------------------------------|
| **FileName**       | `string` | 本地文件的名称。                                 |
| **FileFormatType** | `string` | 文件格式类型，支持 `text` 和 `csv`。                |
| **RowDelimiter**   | `string` | 行分隔符，仅在 `file_format_type` 为 `text` 时生效。 |


## 示例
```json
{
  "sink": [
    {
      "type": "local_file",
      "sink_name": "local_file_data_sink",
      "local_file": {
        "file_name": "output.txt",
        "file_format_type": "text",
        "row_delimiter": "\n"
      }
    }
  ]
}
```