```
PostgresSQL
====

## Overview

The `PostgresSQLSinkConfig` struct is used to configure PostgresSQL Sink, including all settings related to PostgresSQL database connection and data transmission, to batch write data to the specified database table.

## PostgresSQLSinkConfig Struct

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

## Field Description

| Field Name        | Type       | Description                           |
|-------------------|------------|---------------------------------------|
| **Address**       | `string`   | PostgreSQL database address.          |
| **Username**      | `string`   | PostgreSQL username.                  |
| **Password**      | `string`   | PostgreSQL password.                  |
| **Database**      | `string`   | Target PostgreSQL database name.      |
| **TableName**     | `string`   | Target PostgreSQL table name.         |
| **BulkSize**      | `int`      | Number of records processed per batch. |

## Example

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

```
