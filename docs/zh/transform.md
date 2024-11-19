Transform
===

### TransformConfig 配置说明

| 字段名            | 类型                    | 必填 | 说明                                                 |
|----------------|-----------------------|----|----------------------------------------------------|
| `Mode`         | `string`              | 是  | 转换模式：`row`（按行分隔）或 `json`（按 JSON 格式）                |
| `Schemas`      | `[]TransformSchema`   | 否  | 字段映射的配置列表，包含字段转换、别名、转换器等配置                         |
| `RowSeparator` | `string`              | 否  | 行分隔符，仅在 `Mode` 为 `row` 时使用，使用 `strings.Split` 进行分隔 |
| `Paths`        | `[]TransformJsonPath` | 否  | JSON 路径配置列表，仅在 `Mode` 为 `json` 时使用                 |

### TransformSchema 字段映射配置说明

| 字段名            | 类型       | 必填 | 说明                                         |
|----------------|----------|----|--------------------------------------------|
| `SourceKey`    | `string` | 是  | 原始字段名                                      |
| `SinkKey`      | `string` | 是  | 目标字段名                                      |
| `Converter`    | `string` | 否  | 转换器类型，如 `toInt`, `toFloat32`, `toString` 等 |
| `IsIgnore`     | `bool`   | 否  | 是否忽略此字段                                    |
| `IsStrictMode` | `bool`   | 否  | 是否开启严格模式，严格模式下会对字段进行校验                     |
| `IsKeepKeys`   | `bool`   | 否  | 是否保留原始字段名                                  |
| `IsExpand`     | `bool`   | 否  | 是否展开此字段（用于列扩展）                             |
| `ExpandValue`  | `any`    | 否  | 展开字段时的值                                    |
| `SourceName`   | `string` | 是  | 数据来源别名                                     |
| `SinkName`     | `string` | 是  | 目标数据别名                                     |

#### Converter 支持的类型说明

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



### TransformJsonPath JSON 路径配置说明

| 字段名         | 类型       | 必填 | 说明                   |
|-------------|----------|----|----------------------|
| `SrcField`  | `string` | 是  | JSON 源字段名            |
| `Path`      | `string` | 是  | JSON 路径（使用 JSONPath） |
| `DestField` | `string` | 是  | 目标字段名                |

### 使用示例

#### 示例 1：`row` 模式配置
```json
{
  "mode": "row",
  "row_separator": ",",
  "schemas": [
    {
      "source_key": "col1",
      "sink_key": "column1",
      "converter": "toInt",
      "is_ignore": false,
      "is_strict_mode": true,
      "is_keep_keys": false,
      "source_name": "source1",
      "sink_name": "sink1"
    }
  ]
}
```

#### 示例 2：json 模式配置
```json
{
  "mode": "json",
  "schemas": [
    {
      "source_key": "data",
      "sink_key": "username",
      "converter": "toString",
      "is_ignore": false,
      "is_strict_mode": true,
      "is_keep_keys": true,
      "is_expand": false,
      "expand_value": null,
      "source_name": "User Data",
      "sink_name": "User Info"
    }
  ],
  "paths": [
    {
      "src_field": "data",
      "path": "{UserNamePath}",
      "dest_field": "username"
    },
    {
      "src_field": "data",
      "path": "{UserAgePath}",
      "dest_field": "user_age",
      "converter": "toInt",
      "is_strict_mode": true
    }
  ]
}


```