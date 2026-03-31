Transform
===

### TransformConfig Configuration

| Field Name | Type | Required | Description |
|------------|------|----------|-------------|
| `Mode` | `string` | Yes | Transformation mode: `row` (line-separated), `json` (JSON format), or `jsonPath` (extract via JSONPath) |
| `Schemas` | `[]TransformSchema` | No | List of field mapping configurations, including field conversion, aliases, converters, etc. |
| `RowSeparator` | `string` | No | Line separator, only used when `Mode` is `row`, separated using `strings.Split` |
| `Paths` | `[]TransformJsonPath` | No | List of JSON path configurations, only used when `Mode` is `json` |

### TransformSchema Field Mapping Configuration

| Field Name | Type | Required | Description |
|------------|------|----------|-------------|
| `SourceKey` | `string` | Yes | Original field name |
| `SinkKey` | `string` | Yes | Target field name |
| `Converter` | `string` | No | Converter type, such as `toInt`, `toFloat32`, `toString`, etc. |
| `IsIgnore` | `bool` | No | Whether to ignore this field |
| `IsStrictMode` | `bool` | No | Whether to enable strict mode, fields will be validated in strict mode |
| `IsKeepKeys` | `bool` | No | Whether to keep original field names |
| `IsExpand` | `bool` | No | Whether to expand this field (for column expansion) |
| `ExpandValue` | `any` | No | Value when expanding the field |
| `SourceName` | `string` | Yes | Data source alias |
| `SinkName` | `string` | Yes | Target data alias |

#### Supported Converter Types

| Converter Type | Description |
|---------------|-------------|
| `toBool` | Convert to boolean type (`bool`) |
| `toFloat64` | Convert to 64-bit floating point (`float64`) |
| `toFloat32` | Convert to 32-bit floating point (`float32`) |
| `toInt64` | Convert to 64-bit integer (`int64`) |
| `toInt32` | Convert to 32-bit integer (`int32`) |
| `toInt16` | Convert to 16-bit integer (`int16`) |
| `toInt8` | Convert to 8-bit integer (`int8`) |
| `toInt` | Convert to integer (`int`), size determined by system architecture |
| `toUint` | Convert to unsigned integer (`uint`) |
| `toUint64` | Convert to 64-bit unsigned integer (`uint64`) |
| `toUint32` | Convert to 32-bit unsigned integer (`uint32`) |
| `toUint16` | Convert to 16-bit unsigned integer (`uint16`) |
| `toUint8` | Convert to 8-bit unsigned integer (`uint8`) |
| `toString` | Convert to string (`string`) |
| `toStringMap` | Convert to string key-value pairs (`map[string]interface{}`) |

**Notes**
1. **Converter name must be valid:**  
   The converter type name used in the `Converter` field must be a supported type (such as `toBool`, `toInt64`, etc.). If an unsupported name is used, it will return `nil` and log an "unknown convertor" error.

2. **Conversion failure handling:**  
   During data conversion, if conversion fails (such as data type mismatch or conversion exception), an error will be logged and `nil` will be returned. It is recommended to use it in strict mode (when `IsStrictMode` is `true`) and ensure input data matches the expected type.

3. **Default behavior:**  
   If the `Converter` field is an empty string, no data conversion will be performed, and the original data will be returned directly.

4. **Compatibility note:**  
   The `toInt` converter automatically selects 32-bit or 64-bit integer based on the system architecture, so it may vary in different system environments.

5. **Performance recommendation:**  
   For large batch data conversion, try to avoid using complex converter types (such as `toStringMap`) to reduce performance overhead.

6. **Logging:**  
   All conversion errors and unrecognized converter types will be logged for later troubleshooting and debugging.

### TransformJsonPath JSON Path Configuration

| Field Name | Type | Required | Description |
|------------|------|----------|-------------|
| `SrcField` | `string` | Yes | JSON source field name |
| `Path` | `string` | Yes | JSON path (using JSONPath) |
| `DestField` | `string` | Yes | Target field name |

### Usage Examples

#### Example 1: `row` mode configuration
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

#### Example 2: json mode configuration
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
