```
Console
======

## Overview

Console Sink is the simplest data output method, which directly outputs data to the console (standard output). It is mainly used in the development and debugging phase to quickly view data flow and facilitate troubleshooting.

## Configuration Description

Console Sink does not require additional configuration, just specify `type` as `"Console"` in the sink configuration.

## Example

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

## Notes

- **Recommended for development and debugging phase only**: Console Sink will output all data to standard output, which is not suitable for production environments.
- **Performance impact**: In high throughput scenarios, console output may become a performance bottleneck.
- **Production environment recommendation**: For production environments, please use other Sinks such as Kafka, ClickHouse, HTTP, etc. to ensure data persistence and performance.

```
