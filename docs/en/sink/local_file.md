```
LocalFile
====

## Overview

The `LocalFileSinkConfig` struct is used to configure local file Sink, including all settings related to file storage, to save data to local files.

## LocalFileSinkConfig Struct

```golang
type LocalFileSinkConfig struct {
	FileName       string `json:"file_name"`
	FileFormatType string `json:"file_format_type"` // file format type: text, csv
	RowDelimiter   string `json:"row_delimiter"`    // only file_format_type = text is affect
}
```

## Field Description

| Field Name        | Type       | Description                                       |
|-------------------|------------|--------------------------------------------------|
| **FileName**      | `string`   | Name of the local file.                           |
| **FileFormatType**| `string`   | File format type, supports `text` and `csv`.      |
| **RowDelimiter**  | `string`   | Row delimiter, only effective when `file_format_type` is `text`. |

## Example
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
```
