Elasticsearch
====

## 概述

`ElasticsearchSinkConfig` 结构体用于配置 Elasticsearch Sink，包括与 Elasticsearch 集群连接和数据传输相关的所有设置，以便将数据写入指定的 Elasticsearch 索引。
## ElasticsearchSinkConfig 结构体

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

## 字段说明

| 字段名           | 类型       | 描述                                 |
|---------------|----------|------------------------------------|
| **Address**   | `string` | Elasticsearch 集群地址。                |
| **Username**  | `string` | Elasticsearch 用户名。                 |
| **Password**  | `string` | Elasticsearch 密码。                  |
| **IndexName** | `string` | Elasticsearch 索引名称。                |
| **DocIdName** | `string` | Elasticsearch 文档 ID 名称（将从转换数据中获取）。 |
| **Version**   | `string` | Elasticsearch 版本，仅支持 7.X 或 8.X。    |

## 示例

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
