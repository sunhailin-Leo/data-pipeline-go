JsonPath
========

## 描述
支持使用 JSONPath 选择数据，JsonPath 提供了一种简单而强大的方式来访问和提取 JSON 文档中的数据。


## TransformJsonPath 属性
| 字段名         | 类型       | 必填 | 说明                   |
|-------------|----------|----|----------------------|
| `SrcField`  | `string` | 否  | JSON 源字段名            |
| `Path`      | `string` | 否  | JSON 路径（使用 JSONPath） |
| `DestField` | `string` | 否  | 目标字段名                |

## 示例
假设我们有以下 来自kafka 的 JSON 数据：
```json
{
  "store": {
    "book": [
      {
        "category": "reference",
        "author": "Nigel Rees",
        "title": "Sayings of the Century",
        "price": 8.95
      },
      {
        "category": "fiction",
        "author": "Evelyn Waugh",
        "title": "Sword of Honour",
        "price": 12.99
      },
      {
        "category": "fiction",
        "author": "Herman Melville",
        "title": "Moby Dick",
        "isbn": "0-553-21311-3",
        "price": 8.99
      },
      {
        "category": "fiction",
        "author": "J. R. R. Tolkien",
        "title": "The Lord of the Rings",
        "isbn": "0-395-19395-8",
        "price": 22.99
      }
    ],
    "bicycle": {
      "color": "red",
      "price": 19.95
    }
  }
}
```
### 需求解析
* 需要读取所有书的作者,字段名为 `author`
* 提取所有价格低于 10 的书，字段名为price

### 配置示例
```json
{
  "transform": {
    "mode": "jsonPath",
    "paths": [
      {
        "src_field": "store.book[*].author",
        "path": "$.store.book[*].author",
        "dest_field": "author"
      },
    ],
    "schemas": [
      {
        "source_key": "store.book[*].author",
        "sink_key": "author",
        "source_name": "kafka-1",
        "sink_name": "http-1"
      },
      {
        "source_key": "store.book[?(@.price < 10)].price",
        "sink_key": "price",
        "source_name": "kafka-1",
        "sink_name": "http-1"
      }
    ]
  }
}
```

### 数据转换后,http-1接收到的请求数据
{
  "author": [
    "Nigel Rees",
    "Evelyn Waugh",
    "Herman Melville",
    "J. R. R. Tolkien"
  ],
  "price": [
    8.95,
    8.99
  ]
}