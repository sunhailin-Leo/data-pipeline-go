Oracle
====

# OracleSinkConfig 结构体

## 概述

`OracleSinkConfig` 结构体用于配置 Oracle Sink，包括与 Oracle 数据库连接和数据传输相关的所有设置，以便将数据批量写入指定的 Oracle 表中。

## 结构体定义

```go
type OracleSinkConfig struct {
	Address   string `json:"address"`    // Oracle address
	Username  string `json:"username"`   // Oracle username
	Password  string `json:"password"`   // Oracle password
	Database  string `json:"database"`   // Oracle database
	TableName string `json:"table_name"` // Oracle table name
	BulkSize  int    `json:"bulk_size"`  // Bulk size
}
```

## 字段说明

| 字段名           | 类型       | 描述               |
|---------------|----------|------------------|
| **Address**   | `string` | Oracle 数据库地址。    |
| **Username**  | `string` | Oracle 用户名。      |
| **Password**  | `string` | Oracle 密码。       |
| **Database**  | `string` | 目标 Oracle 数据库名称。 |
| **TableName** | `string` | 目标 Oracle 表名称。   |
| **BulkSize**  | `int`    | 每次批量处理的记录数大小。    |


## 示例
```json
{
  "sink": [
    {
      "type": "oracle",
      "sink_name": "oracle_data_sink",
      "oracle": {
        "address": "localhost:1521",
        "username": "user",
        "password": "password",
        "database": "example_db",
        "table_name": "example_table",
        "bulk_size": 1000
      }
    }
  ]
}
```
