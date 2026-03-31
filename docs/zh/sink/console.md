Console
======

## 概述

Console Sink 是最简单的数据输出方式，将数据直接输出到控制台（标准输出）。它主要用于开发调试阶段，可以快速查看数据流转情况，方便排查问题。

## 配置说明

Console Sink 不需要额外配置，只需在 sink 配置中指定 `type` 为 `"Console"` 即可。

## 示例

```json
{
  "sink": [
    {
      "type": "Console",
      "sink_name": "debug_sink"
    }
  ]
}
```

## 注意事项

- **仅建议在开发调试阶段使用**：Console Sink 会将所有数据输出到标准输出，不适合生产环境
- **性能影响**：在高吞吐量场景下，控制台输出可能成为性能瓶颈
- **生产环境建议**：生产环境请使用其他 Sink，如 Kafka、ClickHouse、HTTP 等，以保证数据持久化和性能
