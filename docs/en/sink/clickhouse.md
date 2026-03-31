```
ClickHouse
=====
## Overview
ClickHouse Sink is responsible for writing data to the ClickHouse database. This configuration allows users to customize connection settings, table structure, and write options for efficient data storage and analysis.

## ClickhouseSinkConfig Struct
```golang
type ClickhouseTableColumn struct {
  Name     string `json:"name"`     // Column name
  Type     string `json:"type"`     // Column data type
  Compress string `json:"compress"` // Column compression method (optional)
  Comment  string `json:"comment"`  // Column comment information (optional)
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

## Field Description
### `ClickhouseSinkConfig`

| Field Name         | Type                       | Description                                                                 |
|--------------------|----------------------------|-----------------------------------------------------------------------------|
| **Address**        | `string`                   | ClickHouse database address (IP or domain), used to establish database connection. |
| **Username**       | `string`                   | Username for connecting to ClickHouse.                                     |
| **Password**       | `string`                   | Password for connecting to ClickHouse.                                     |
| **Database**       | `string`                   | Name of the database to write to.                                          |
| **TableName**      | `string`                   | Target table name for data writing.                                        |
| **BulkSize**       | `int`                      | Bulk write size, determines the number of data records written to ClickHouse each time. |
| **IsAutoCreateTable** | `bool`                   | Whether to automatically create the table in the target database. The table will be created automatically if it does not exist. |
| **Columns**        | `[]ClickhouseTableColumn`  | Column configuration in the table, containing detailed information about column names and types. |
| **Engine**         | `string`                   | Table storage engine type, such as `MergeTree`, `SummingMergeTree`, etc.    |
| **Partition**      | `[]string`                 | List of partition fields, supports partitioning by specified fields to improve query performance. |
| **PrimaryKey**     | `[]string`                 | List of primary key fields, used to uniquely identify each row of data.    |
| **OrderBy**        | `[]string`                 | List of sort fields in the table to improve query efficiency.              |
| **TTL**            | `string`                   | Data expiration time setting, allows specifying the maximum time data is stored in the table. |
| **Comment**        | `string`                   | Table comment information, providing description or purpose of the table.  |
| **Settings**       | `[]string`                 | Other configuration options for ClickHouse table to optimize performance or meet specific requirements. |

### `ClickhouseTableColumn`

| Field Name  | Type     | Description                                  |
|-------------|----------|----------------------------------------------|
| **Name**    | `string` | Column name, defines the field name in the table. |
| **Type**    | `string` | Column data type, such as `UInt32`, `String`, etc. |
| **Compress** | `string` | Column compression method, specifies how to compress data in this column (optional). |
| **Comment** | `string` | Column comment information, providing description about the column (optional). |

## Example
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
            "comment": "User unique identifier"
          },
          {
            "name": "name",
            "type": "String",
            "compress": "",
            "comment": "User name"
          },
          {
            "name": "created_at",
            "type": "DateTime",
            "compress": "",
            "comment": "User creation time"
          }
        ],
        "engine": "MergeTree()",
        "partition": ["toYYYYMM(created_at)"],
        "primary_key": ["id"],
        "order_by": ["created_at"],
        "ttl": "created_at + INTERVAL 30 DAY",
        "comment": "Table for storing user data",
        "settings": [
          "max_insert_block_size=1000",
          "enable_http_compression=1"
        ]
      }
    }

  ]
}

```
```
