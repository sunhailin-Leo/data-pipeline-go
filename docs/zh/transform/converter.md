类型转换
================
> 字段类型转换

## TransformSchema 属性配置
| 字段        | 类型     | 必填 | 默认值 | 说明            |
|-----------|--------|----|-----|---------------|
| converter | string | 否  | -   | 将源字段转换为目标字段类型 |

## converter 支持的类型说明

| 转换器类型         | 说明                                   |
|---------------|--------------------------------------|
| `toBool`      | 转换为布尔类型 (`bool`)                     |
| `toFloat64`   | 转换为 64 位浮点数 (`float64`)              |
| `toFloat32`   | 转换为 32 位浮点数 (`float32`)              |
| `toInt64`     | 转换为 64 位整数 (`int64`)                 |
| `toInt32`     | 转换为 32 位整数 (`int32`)                 |
| `toInt16`     | 转换为 16 位整数 (`int16`)                 |
| `toInt8`      | 转换为 8 位整数 (`int8`)                   |
| `toInt`       | 转换为整数 (`int`)，根据系统位数决定大小             |
| `toUint`      | 转换为无符号整数 (`uint`)                    |
| `toUint64`    | 转换为 64 位无符号整数 (`uint64`)             |
| `toUint32`    | 转换为 32 位无符号整数 (`uint32`)             |
| `toUint16`    | 转换为 16 位无符号整数 (`uint16`)             |
| `toUint8`     | 转换为 8 位无符号整数 (`uint8`)               |
| `toString`    | 转换为字符串 (`string`)                    |
| `toStringMap` | 转换为字符串键值对 (`map[string]interface{}`) |
**注意事项**
1. **转换器名称必须合法：**
   `Converter` 字段中使用的转换器类型名称必须是支持的类型（如 `toBool`, `toInt64` 等）。如果使用不支持的名称，将返回 `nil` 并记录 "unknown convertor" 错误日志。

2. **转换失败处理：**
   在数据转换过程中，如果转换失败（如数据类型不匹配或转换异常），会记录错误日志并返回 `nil`。建议在严格模式（`IsStrictMode` 为 `true`）下使用时，确保输入数据符合预期类型。

3. **默认行为：**
   如果 `Converter` 字段为空字符串，将不会进行任何数据转换，直接返回原始数据。

4. **兼容性注意：**
   `toInt` 转换器根据系统架构的位数自动选择为 32 位或 64 位整数，因此在不同系统环境下可能会有所不同。

5. **性能建议：**
   对于大批量数据转换，尽量避免使用复杂的转换器类型（如 `toStringMap`），以减少性能开销。

6. **日志记录：**
   所有转换错误和未识别的转换器类型都会记录到日志中，便于后续排查和调试。

## 示例
### 源数据表,来源 mysql-1：
| id | name | age（int） |
|----|------|----------|
| 1  | Jane | 20       |
| 2  | Max  | 20       |
| 3  | King | 20       |
| 4  | Tom  | 20       |

```json
{
  "transform": {
    "mode": "json",
    "schemas": [
      {
        "source_key": "age",
        "sink_key": "age",
        "converter": "toString",
        "source_name": "mysql-1",
        "sink_name": "mysql-2"
      }
    ]
  }
}
```
### 转换后数据表（写入到 `mysql-2`）：
| id | name | age（string） |
|----|------|-------------|
| 1  | Jane | 20          |
| 2  | Max  | 20          |
| 3  | King | 20          |
| 4  | Tom  | 20          |


此配置表示从 Kafka 数据源 kafka-1 读取数据，将 key 字段转为字符串后写入 HTTP 目标 http-1
在 schemas 数组中，有一个转换模式的配置对象，包含以下属性：
* source_key: 源数据中的键，即需要转换的字段。
* sink_key: 转换后的数据在目标中的键。
* converter: 转换器类型，这里是 toString，表示将数据转换为字符串类型。
* source_name: 源数据的名称，这里是 kafka-1。
* sink_name: 目标数据的名称，这里是 http-1。

