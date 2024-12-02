ClickHouse
=====
## 概述
ClickHouse Sink 负责将数据写入 ClickHouse 数据库。此配置允许用户自定义连接设置、数据表结构以及写入选项，以便实现高效的数据存储与分析。

## ClickhouseSinkConfig 结构体
```golang
type ClickhouseTableColumn struct {
  Name     string `json:"name"`     // 列的名称
  Type     string `json:"type"`     // 列的数据类型
  Compress string `json:"compress"` // 列的压缩方式（可选）
  Comment  string `json:"comment"`  // 列的注释信息（可选）
}

type ClickhouseSinkConfig struct {
	Address           string                  `json:"address"`
	Username          string                  `json:"username"`
	Password          string                  `json:"password"`
	Database          string                  `json:"database"`
	TableName         string                  `json:"table_name"`
	BulkSize          int                     `json:"bulk_size"`
	IsAutoCreateTable bool                    `json:"is_auto_create_table"`
	Columns           []ClickhouseTableColumn `json:"columns"`
	Engine            string                  `json:"engine"`
	Partition         []string                `json:"partition"`
	PrimaryKey        []string                `json:"primary_key"`
	OrderBy           []string                `json:"order_by"`
	TTL               string                  `json:"ttl"`
	Comment           string                  `json:"comment"`
	Settings          []string                `json:"settings"`
}
```

## 字段说明
### `ClickhouseSinkConfig`

| 字段名               | 类型                       | 描述                                                                 |
|---------------------|----------------------------|----------------------------------------------------------------------|
| **Address**         | `string`                   | ClickHouse 数据库的地址（IP 或域名），用于建立数据库连接。            |
| **Username**        | `string`                   | 用于连接 ClickHouse 的用户名。                                      |
| **Password**        | `string`                   | 用于连接 ClickHouse 的密码。                                        |
| **Database**        | `string`                   | 要写入的数据库名称。                                                |
| **TableName**       | `string`                   | 数据写入的目标表名称。                                              |
| **BulkSize**        | `int`                      | 批量写入的大小，决定每次写入 ClickHouse 的数据条数。                |
| **IsAutoCreateTable** | `bool`                   | 是否在目标数据库中自动创建表，如果表不存在则自动创建。              |
| **Columns**         | `[]ClickhouseTableColumn`  | 表中列的配置，包含列名和类型的详细信息。                           |
| **Engine**          | `string`                   | 表的存储引擎类型，如 `MergeTree`、`SummingMergeTree` 等。           |
| **Partition**       | `[]string`                 | 分区字段的列表，支持根据指定字段进行分区以提高查询性能。            |
| **PrimaryKey**      | `[]string`                 | 主键字段的列表，用于唯一标识每一行数据。                            |
| **OrderBy**         | `[]string`                 | 数据在表中的排序字段列表，以提高查询效率。                         |
| **TTL**             | `string`                   | 数据的过期时间设置，允许指定数据在表中存储的最大时间。              |
| **Comment**         | `string`                   | 表的注释信息，提供表的描述或用途。                                  |
| **Settings**        | `[]string`                 | ClickHouse 表的其他配置选项，以优化性能或满足特定需求。            |


### `ClickhouseTableColumn`

| 字段名    | 类型     | 描述                                  |
|-----------|----------|---------------------------------------|
| **Name**  | `string` | 列的名称，定义在表中的字段名。       |
| **Type**  | `string` | 列的数据类型，例如 `UInt32`、`String` 等。 |
| **Compress** | `string` | 列的压缩方式，指定如何压缩该列的数据（可选）。 |
| **Comment** | `string` | 列的注释信息，提供关于列的描述（可选）。 |


## 示例
```json
{
  "sink": [
    {
      "type": "clickhouse",
      "sink_name": "user_data_sink",
      "clickhouse": {
        "address": "http://localhost:8123",
        "username": "default",
        "password": "password123",
        "database": "analytics",
        "table_name": "users",
        "bulk_size": 1000,
        "is_auto_create_table": true,
        "columns": [
          {
            "name": "id",
            "type": "UInt64",
            "compress": "LZ4",
            "comment": "用户唯一标识符"
          },
          {
            "name": "name",
            "type": "String",
            "compress": "",
            "comment": "用户姓名"
          },
          {
            "name": "created_at",
            "type": "DateTime",
            "compress": "",
            "comment": "用户创建时间"
          }
        ],
        "engine": "MergeTree()",
        "partition": ["toYYYYMM(created_at)"],
        "primary_key": ["id"],
        "order_by": ["created_at"],
        "ttl": "created_at + INTERVAL 30 DAY",
        "comment": "存储用户数据的表",
        "settings": [
          "max_insert_block_size=1000",
          "enable_http_compression=1"
        ]
      }
    }

  ]
}

```