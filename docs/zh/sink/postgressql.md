PostgresSQL
====


## 概述

`PostgresSQLSinkConfig` 结构体用于配置 PostgresSQL Sink，包括与 PostgresSQL 数据库连接和数据传输相关的所有设置，以便将数据批量写入指定的数据库表中。
## PostgresSQLSinkConfig 结构体

```golang
type PostgresSQLSinkConfig struct {
	Address   string `json:"address"`    // PostgresSQL address
	Username  string `json:"username"`   // PostgresSQL username
	Password  string `json:"password"`   // PostgresSQL password
	Database  string `json:"database"`   // PostgresSQL database
	TableName string `json:"table_name"` // PostgresSQL table name
	BulkSize  int    `json:"bulk_size"`  // Bulk size
}
```

## 字段说明

| 字段名           | 类型       | 描述                   |
|---------------|----------|----------------------|
| **Address**   | `string` | PostgreSQL 数据库地址。    |
| **Username**  | `string` | PostgreSQL 用户名。      |
| **Password**  | `string` | PostgreSQL 密码。       |
| **Database**  | `string` | 目标 PostgreSQL 数据库名称。 |
| **TableName** | `string` | 目标 PostgreSQL 表名称。   |
| **BulkSize**  | `int`    | 每次批量处理的记录数大小。        |


## 示例

```json
{
  "sink": [
    {
      "type": "postgresql",
      "sink_name": "postgresql_data_sink",
      "postgresql": {
        "address": "localhost:5432",
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
