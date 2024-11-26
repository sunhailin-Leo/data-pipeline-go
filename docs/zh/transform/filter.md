字段过滤器
========

## TransformSchema 属性配置
| 字段        | 类型   | 必填 | 默认值   | 说明      |
|-----------|------|----|-------|---------|
| is_ignore | bool | 否  | false | 是否过滤该字段 |

## 示例
### 源数据表,来源 mysql-1：
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
### 转换后数据表（写入到 `mysql-2`）：
| id | name | age |
|----|------|-----|
| 1  | Jane | 20  |
| 2  | Max  | 20  |
| 3  | King | 20  |
| 4  | Tom  | 20  |
