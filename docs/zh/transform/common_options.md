转换常见选项
======
> 源端连接器的常见参数
## 常见选项
| 字段名       | 类型                  | 必填 | 默认值 | 说明                                                           |
|-----------|---------------------|----|-----|--------------------------------------------------------------|
| `mode`    | `string`            | 是  | -   | 支持转换模式：`row`（按行分隔）; `json`（按 JSON 格式） ;`jsonPath` (数据访问路径模式) | 
| `schemas` | `[]TransformSchema` | 是  | -   | 过滤器配制列表                                                      |

* `row` 模式 以行读取来源数据，可自定义分隔符，支持处理多种不同的文件格式。
* `json` 模式 以 JSON 格式解析来源数据，支持提取整个 JSON 对象。
* `jsonPath` 通过 jsonPath 提取特定字段值。适用于嵌套结构复杂的 JSON 数据。

## 模式对比
| **对比维度**  | **row 模式**      | **json 模式**        | **jsonPath 模式**      |
|-----------|-----------------|--------------------|----------------------|
| **数据结构**  | 平坦结构，按行划分       | 嵌套结构，键值对格式         | 嵌套结构，指定字段路径          |
| **灵活性**   | 中等，需配置分隔符       | 高，支持动态字段           | 高，精准提取需要的字段          |
| **适合数据源** | CSV、TSV 文件，日志文本 | JSON 文件或 API 返回值   | JSON 文件或 API 返回值     |
| **配置复杂度** | 简单，仅需配置分隔符      | 中等，需匹配 JSON 结构     | 较高，需熟悉 JsonPath 查询语法 |
| **性能**    | 高效，适合大规模平坦数据    | 较高，适合解析中小型 JSON 数据 | 高效，仅解析所需部分字段         |

## 过滤器配制参数
| 字段名           | 类型       | 必填 | 默认值 | 说明                                  |
|---------------|----------|----|-----|-------------------------------------|
| `source_name` | `string` | 是  | -   | 源数据来源名，对应 source 配置中的 `source_name` |
| `source_key`  | `string` | 是  | -   | 在从 SourceName 读取数据时应该使用的键或字段名       |
| `sink_name`   | `string` | 是  | -   | 数据目标名，对应 sink 配置中的 `sink_name`      |
| `sink_key`    | `string` | 是  | -   | 用于指定在将数据写入到 SinkName 时应该使用的键或字段名    |

* 在数据处理管道中，这四个字段用于定义数据从源头到目标的映射关系，确保数据能够正确读取、转换和写入。它们共同构成了数据在管道中的映射逻辑，确保从源到目标的完整传输与转换。
* SourceName 与 SinkName 必须对应 source 与 sink 配置中的 `source_name` 和 `sink_name` 字段值。
* SourceName 和 SourceKey 负责定义数据来源及其字段；SinkName 和 SinkKey 负责定义数据的目标及其字段。

## 字段映射示例
### 源数据表,来源 mysql-1：
| id | name | age |
|----|------|-----|
| 1  | Jane | 20  |
| 2  | Max  | 20  |
| 3  | King | 20  |
| 4  | Tom  | 20  |

#### 转换需求
1. . 将字段 `name` 重命名为 `new_name`。
2. 将字段 `age` 重命名为 `user_age`。

#### 配置示例
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

#### 转换后数据表（写入到 `clickhouse-1`）：
| id | new_name | user_age |
|----|----------|----------|
| 1  | Jane     | 20       |
| 2  | Max      | 20       |
| 3  | King     | 20       |
| 4  | Tom      | 20       |


