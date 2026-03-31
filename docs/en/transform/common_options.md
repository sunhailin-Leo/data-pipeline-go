Common Options
======
> Common parameters for source connectors
## Common Options
| Field Name   | Type                  | Required | Default | Description                                                           |
|-------------|---------------------|----|-----|--------------------------------------------------------------|
| `mode`      | `string`            | Yes  | -   | Supported transform modes: `row` (line-by-line); `json` (JSON format); `jsonPath` (data access path mode) | 
| `schemas`   | `[]TransformSchema` | Yes  | -   | Filter configuration list                                                      |

* `row` mode reads source data line by line, supports custom delimiters, and handles various file formats.
* `json` mode parses source data in JSON format, supports extracting entire JSON objects.
* `jsonPath` extracts specific field values through jsonPath. Suitable for JSON data with complex nested structures.

## Mode Comparison
| **Comparison Dimension** | **row mode**      | **json mode**        | **jsonPath mode**      |
|------------------------|-----------------|--------------------|----------------------|
| **Data Structure**     | Flat structure, line-by-line       | Nested structure, key-value format         | Nested structure, specified field paths          |
| **Flexibility**        | Medium, requires delimiter configuration       | High, supports dynamic fields           | High, precisely extracts required fields          |
| **Suitable Data Sources** | CSV, TSV files, log text | JSON files or API return values   | JSON files or API return values     |
| **Configuration Complexity** | Simple, only delimiter configuration needed      | Medium, needs to match JSON structure     | Higher, requires familiarity with JsonPath query syntax |
| **Performance**        | Efficient, suitable for large-scale flat data    | Relatively high, suitable for parsing small to medium JSON data | Efficient, only parses required fields         |

## Filter Configuration Parameters
| Field Name        | Type       | Required | Default | Description                                  |
|------------------|----------|----|-----|-------------------------------------|
| `source_name`    | `string` | Yes  | -   | Source data name, corresponding to `source_name` in source configuration |
| `source_key`     | `string` | Yes  | -   | Key or field name to use when reading data from SourceName       |
| `sink_name`      | `string` | Yes  | -   | Data target name, corresponding to `sink_name` in sink configuration      |
| `sink_key`       | `string` | Yes  | -   | Key or field name to use when writing data to SinkName    |

* In the data processing pipeline, these four fields define the mapping relationship from source to target, ensuring data can be correctly read, transformed, and written. Together they form the mapping logic of data in the pipeline, ensuring complete transmission and transformation from source to target.
* SourceName and SinkName must correspond to the `source_name` and `sink_name` field values in source and sink configurations.
* SourceName and SourceKey define the data source and its fields; SinkName and SinkKey define the data target and its fields.

## Field Mapping Example
### Source Data Table, from mysql-1:
| id | name | age |
|----|------|-----|
| 1  | Jane | 20  |
| 2  | Max  | 20  |
| 3  | King | 20  |
| 4  | Tom  | 20  |

#### Transformation Requirements
1. Rename field `name` to `new_name`.
2. Rename field `age` to `user_age`.

#### Configuration Example
```json
{
  "transform": {
    "mode": "json",
    "schemas": [
      {
        "source_key": "name",
        "sink_key": "new_name",
        "source_name": "mysql-1",
        "sink_name": "kafka-1"
      },
      {
        "source_key": "age",
        "sink_key": "user_age",
        "source_name": "mysql-1",
        "sink_name": "clickhouse-1"
      }
    ]
  }
}
```

#### Transformed Data Table (written to `clickhouse-1`):
| id | new_name | user_age |
|----|----------|----------|
| 1  | Jane     | 20       |
| 2  | Max      | 20       |
| 3  | King     | 20       |
| 4  | Tom      | 20       |
