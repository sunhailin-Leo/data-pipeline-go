JsonPath
========

## Description
Supports using JSONPath to select data. JsonPath provides a simple and powerful way to access and extract data from JSON documents.

## TransformJsonPath Attributes
| Field Name   | Type       | Required | Description                   |
|-------------|----------|----|----------------------|
| `SrcField`  | `string` | No  | JSON source field name            |
| `Path`      | `string` | No  | JSON path (using JSONPath) |
| `DestField` | `string` | No  | Target field name                |

> Use (GJSON Playground to test syntax online)[https://gjson.dev/]

## Example
Assume we have the following JSON data from kafka:
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
### Requirement Analysis
* Need to read all book authors, field name is `author`
* Extract all books with price below 10, field name is price

### Configuration Example
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

### After Data Transformation, Request Data Received by http-1
```json
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
```
