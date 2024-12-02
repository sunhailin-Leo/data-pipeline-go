行拆分
========= 
> 在 row 模式下，可以设置特定分隔符拆分数据

| 字段            | 类型     | 必填 | 默认值  | 说明   |
|---------------|--------|----|------|------|
| row_separator | string | 否  | `""` | 行分隔符 |

在行模式下,根据 行分隔符 拆分数据，默认使用 空字符串 进行拆分。

**在拆分后字段数量 与 transform.schemas 长度必须一致**

## 示例：
kafka 数据中 data 存在日志格式如下：
```text
2024-11-27 14:40:05.471|-|INFO|-|"data-pipeline-go/pkg/stream/handler.go:135"|-|github.com/sunhailin-Leo/data-pipeline-go/pkg/stream.(*Handler).Start|-|data-pipeline-go|-|service start successful!
```
### 需求
* 使用 `|-|` 将信息按照 日期, 日志级别, 日志位置, 函数名, 项目名称, 日志内容 进行拆分
* 全部字段转为 string 写入 clickhouse 

```json
{
  "transform": {
    "mode": "row",
    "row_separator": "|-|",
    "schemas": [
      {
        "source_key": "date",
        "sink_key": "data",
        "converter": "toString",
        "source_name": "kafka-1",
        "sink_name": "clickhouse-1"
      },
      {
        "source_key": "date",
        "sink_key": "level",
        "converter": "toString",
        "source_name": "kafka-1",
        "sink_name": "clickhouse-1"
      },
      {
        "source_key": "date",
        "sink_key": "path",
        "converter": "toString",
        "source_name": "kafka-1",
        "sink_name": "clickhouse-1"
      },
      {
        "source_key": "date",
        "sink_key": "function",
        "converter": "toString",
        "source_name": "kafka-1",
        "sink_name": "clickhouse-1"
      },
      {
        "source_key": "date",
        "sink_key": "project",
        "converter": "toString",
        "source_name": "kafka-1",
        "sink_name": "clickhouse-1"
      },
      {
        "source_key": "date",
        "sink_key": "message",
        "converter": "toString",
        "source_name": "kafka-1",
        "sink_name": "clickhouse-1"
      }
    ]
  }
}
```
### 转换后数据表（写入到 `clickhouse-1`）：
| date       | level | path                                       | function                                                              | project          | message                   |
|------------|-------|--------------------------------------------|-----------------------------------------------------------------------|------------------|---------------------------|
| 2024-11-27 | INFO  | data-pipeline-go/pkg/stream/handler.go:135 | github.com/sunhailin-Leo/data-pipeline-go/pkg/stream.(*Handler).Start | data-pipeline-go | service start successful! |