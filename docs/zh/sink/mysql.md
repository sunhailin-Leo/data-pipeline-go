MySQL
====

## 概述

`MySQLSinkConfig` 结构体用于配置 MySQL Sink，包括与 MySQL 数据库连接和数据传输相关的所有设置，以便将数据批量写入指定的 MySQL 表中。
## MySQLSinkConfig 结构体


```go
type MySQLSinkConfig struct {
Address   string `json:"address"`    // MySQL address
Username  string `json:"username"`   // MySQL username
Password  string `json:"password"`   // MySQL password
Database  string `json:"database"`   // MySQL database
TableName string `json:"table_name"` // MySQL table name
BulkSize  int    `json:"bulk_size"`  // Bulk size
}
```


## 字段说明

| 字段名           | 类型       | 描述              |
|---------------|----------|-----------------|
| **Address**   | `string` | MySQL 数据库地址。    |
| **Username**  | `string` | MySQL 用户名。      |
| **Password**  | `string` | MySQL 密码。       |
| **Database**  | `string` | 目标 MySQL 数据库名称。 |
| **TableName** | `string` | 目标 MySQL 表名称。   |
| **BulkSize**  | `int`    | 每次批量处理的记录数大小。   |


## 示例
```json
{
  "sink": [
    {
      "type": "mysql",
      "sink_name": "mysql_data_sink",
      "mysql": {
        "address": "localhost:3306",
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