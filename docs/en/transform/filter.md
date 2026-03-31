Field Filter
========

## TransformSchema Attribute Configuration
| Field      | Type   | Required | Default   | Description      |
|-----------|------|----|-------|---------|
| is_ignore | bool | No  | false | Whether to filter this field |

## Example
### Source Data Table, from mysql-1:
| id | name | age | card |
|----|------|-----|------|
| 1  | Jane | 20  | 123  |
| 2  | Max  | 20  | 123  |
| 3  | King | 20  | 123  |
| 4  | Tom  | 20  | 123  |

```json
{
  "transform": {
    "mode": "json",
    "schemas": [
      {
        "source_key": "card",
        "sink_key": " ",
        "is_ignore": true,
        "source_name": "mysql-1",
        "sink_name": "mysql-2"
      }
    ]
  }
}
```
### Transformed Data Table (written to `mysql-2`):
| id | name | age |
|----|------|-----|
| 1  | Jane | 20  |
| 2  | Max  | 20  |
| 3  | King | 20  |
| 4  | Tom  | 20  |
