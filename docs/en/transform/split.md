Row Split
========= 
> In row mode, you can set a specific delimiter to split data

| Field         | Type     | Required | Default  | Description   |
|---------------|--------|----|------|------|
| row_separator | string | No  | `""` | Row separator |

In row mode, split data according to the row separator. By default, empty string is used for splitting.

**After splitting, the number of fields must match the length of transform.schemas**

## Example:
In kafka data, the data field exists in log format as follows:
```text
2024-11-27 14:40:05.471|-|INFO|-|"data-pipeline-go/pkg/stream/handler.go:135"|-|github.com/sunhailin-Leo/data-pipeline-go/pkg/stream.(*Handler).Start|-|data-pipeline-go|-|service start successful!
```
### Requirements
* Use `|-|` to split information into date, log level, log location, function name, project name, log content
* Convert all fields to string and write to clickhouse

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
### Transformed Data Table (written to `clickhouse-1`):
| date       | level | path                                       | function                                                              | project          | message                   |
|------------|-------|--------------------------------------------|-----------------------------------------------------------------------|------------------|---------------------------|
| 2024-11-27 | INFO  | data-pipeline-go/pkg/stream/handler.go:135 | github.com/sunhailin-Leo/data-pipeline-go/pkg/stream.(*Handler).Start | data-pipeline-go | service start successful! |
