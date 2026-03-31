Field Type Converter
================
> Field type conversion

## TransformSchema Attribute Configuration
| Field     | Type     | Required | Default | Description            |
|-----------|--------|----|-----|---------------|
| converter | string | No  | -   | Convert source field to target field type |

## Converter Supported Types

| Converter Type  | Description                                   |
|---------------|--------------------------------------|
| `toBool`      | Convert to boolean type (`bool`)                     |
| `toFloat64`   | Convert to 64-bit floating point (`float64`)              |
| `toFloat32`   | Convert to 32-bit floating point (`float32`)              |
| `toInt64`     | Convert to 64-bit integer (`int64`)                 |
| `toInt32`     | Convert to 32-bit integer (`int32`)                 |
| `toInt16`     | Convert to 16-bit integer (`int16`)                 |
| `toInt8`      | Convert to 8-bit integer (`int8`)                   |
| `toInt`       | Convert to integer (`int`), size determined by system bits             |
| `toUint`      | Convert to unsigned integer (`uint`)                    |
| `toUint64`    | Convert to 64-bit unsigned integer (`uint64`)             |
| `toUint32`    | Convert to 32-bit unsigned integer (`uint32`)             |
| `toUint16`    | Convert to 16-bit unsigned integer (`uint16`)             |
| `toUint8`     | Convert to 8-bit unsigned integer (`uint8`)               |
| `toString`    | Convert to string (`string`)                    |
| `toStringMap` | Convert to string key-value pair (`map[string]interface{}`) |

**Notes**
1. **Converter name must be valid:**
   The converter type name used in the `Converter` field must be a supported type (such as `toBool`, `toInt64`, etc.). If an unsupported name is used, `nil` will be returned and an "unknown convertor" error log will be recorded.

2. **Conversion failure handling:**
   During data conversion, if conversion fails (e.g., data type mismatch or conversion exception), an error log will be recorded and `nil` will be returned. It is recommended to ensure input data matches the expected type when using strict mode (`IsStrictMode` is `true`).

3. **Default behavior:**
   If the `Converter` field is an empty string, no data conversion will be performed and the original data will be returned directly.

4. **Compatibility note:**
   The `toInt` converter automatically selects 32-bit or 64-bit integer based on the system architecture's bit size, so it may vary in different system environments.

5. **Performance recommendation:**
   For large batch data conversion, try to avoid using complex converter types (such as `toStringMap`) to reduce performance overhead.

6. **Logging:**
   All conversion errors and unrecognized converter types will be logged for subsequent troubleshooting and debugging.

## Example
### Source Data Table, from mysql-1:
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
### Transformed Data Table (written to `mysql-2`):
| id | name | age（string） |
|----|------|-------------|
| 1  | Jane | 20          |
| 2  | Max  | 20          |
| 3  | King | 20          |
| 4  | Tom  | 20          |

This configuration indicates reading data from Kafka data source kafka-1, converting the key field to string, and writing it to HTTP target http-1
In the schemas array, there is a transform mode configuration object containing the following attributes:
* source_key: The key in the source data, i.e., the field to be converted.
* sink_key: The key of the transformed data in the target.
* converter: Converter type, here it is toString, indicating conversion of data to string type.
* source_name: The name of the source data, here it is kafka-1.
* sink_name: The name of the target data, here it is http-1.
