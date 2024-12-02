Prometheus
====
## Prometheus Source 配置
prom_metrics 是一个对象类型，用于配置从 Prometheus 获取指标的相关参数。下面是具体的配置说明：

```json
{
  "source": [
    {
      "type": "Prometheus",
      "source_name": "prometheus-1",
      "prom_metrics": {
        "address": "<Prometheus Address>",   // Prometheus 的地址，例如 "http://prometheus.example.com:9090"
        "interval": <Interval in seconds>    // 查询的时间间隔（秒），例如 30
      }
    }
  ]
}
```
## 配置参数说明:
**type：**

指定 source 的类型，这里是 "Prometheus"，表示使用 Prometheus 作为数据源。

**source_name：**

Prometheus source 的自定义名称，用于识别该 source。在本例中，它被设置为 "prometheus-1"。

**prom_metrics：**
包含 Prometheus source 的核心配置信息：

* **address：**

    Prometheus 实例的地址，通常为 HTTP 格式。例如，"http://prometheus.example.com:9090" 指定 Prometheus 的 API 访问地址。

* **interval：**

    查询 Prometheus 的时间间隔，单位为秒。定义了从 Prometheus 拉取最新指标数据的频率。例如，30 表示每 30 秒查询一次 Prometheus。

## 配置示例

```json
{
  "source": [
    {
      "type": "Prometheus",
      "source_name": "prometheus-1",
      "prom_metrics": {
        "address": "http://prometheus.example.com:9090",
        "interval": 30
      }
    }
  ]
}
```
在该示例中：

* "prometheus-1" 是定义的 Prometheus source 名称。

* "address" 指定了 Prometheus 实例的地址 "http://prometheus.example.com:9090"。

* "interval" 设置为 30 秒，表示每 30 秒从 Prometheus 拉取一次数据。