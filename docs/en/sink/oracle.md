```
Oracle
====

## Overview

The `OracleSinkConfig` struct is used to configure Oracle Sink, including all settings related to Oracle database connection and data transmission, to batch write data to the specified Oracle table.

## OracleSinkConfig Struct

```golang
type OracleSinkConfig struct {
	Address   string `json:"address"`    // Oracle address
	Username  string `json:"username"`   // Oracle username
	Password  string `json:"password"`   // Oracle password
	Database  string `json:"database"`   // Oracle database
	TableName string `json:"table_name"` // Oracle table name
	BulkSize  int    `json:"bulk_size"`  // Bulk size
}
```

## Field Description

| Field Name        | Type       | Description                        |
|-------------------|------------|------------------------------------|
| **Address**       | `string`   | Oracle database address.           |
| **Username**      | `string`   | Oracle username.                   |
| **Password**      | `string`   | Oracle password.                   |
| **Database**      | `string`   | Target Oracle database name.       |
| **TableName**     | `string`   | Target Oracle table name.          |
| **BulkSize**      | `int`      | Number of records processed per batch. |

## Example
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

```
