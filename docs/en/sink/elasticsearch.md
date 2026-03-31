```
Elasticsearch
====

## Overview

The `ElasticsearchSinkConfig` struct is used to configure Elasticsearch Sink, including all settings related to Elasticsearch cluster connection and data transmission, to write data to the specified Elasticsearch index.

## ElasticsearchSinkConfig Struct

```go
type ElasticsearchSinkConfig struct {
Address   string `json:"address"`     // Elasticsearch address
Username  string `json:"username"`    // Elasticsearch username
Password  string `json:"password"`    // Elasticsearch password
IndexName string `json:"index_name"`  // Elasticsearch index name
DocIdName string `json:"doc_id_name"` // Elasticsearch document id name (will take it from transform data)
Version   string `json:"version"`     // Only use 7.X or 8.X
}
```

## Field Description

| Field Name        | Type       | Description                                 |
|-------------------|------------|---------------------------------------------|
| **Address**       | `string`   | Elasticsearch cluster address.              |
| **Username**      | `string`   | Elasticsearch username.                     |
| **Password**      | `string`   | Elasticsearch password.                     |
| **IndexName**     | `string`   | Elasticsearch index name.                   |
| **DocIdName**     | `string`   | Elasticsearch document ID name (will be taken from transform data). |
| **Version**       | `string`   | Elasticsearch version, only supports 7.X or 8.X. |

## Example

```json
{
  "sink": [
    {
      "type": "elasticsearch",
      "sink_name": "elasticsearch_data_sink",
      "elasticsearch": {
        "address": "http://localhost:9200",
        "username": "elastic",
        "password": "password",
        "index_name": "example_index",
        "doc_id_name": "id",
        "version": "7.X"
      }
    }
  ]
}

```

```
