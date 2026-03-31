```
MySQL
====

## Overview

The `MySQLSinkConfig` struct is used to configure MySQL Sink, including all settings related to MySQL database connection and data transmission, to batch write data to the specified MySQL table.

## MySQLSinkConfig Struct

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

## Field Description

| Field Name        | Type       | Description                     |
|-------------------|------------|---------------------------------|
| **Address**       | `string`   | MySQL database address.         |
| **Username**      | `string`   | MySQL username.                 |
| **Password**      | `string`   | MySQL password.                 |
| **Database**      | `string`   | Target MySQL database name.     |
| **TableName**     | `string`   | Target MySQL table name.        |
| **BulkSize**      | `int`      | Number of records processed per batch. |

## Example
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
```
