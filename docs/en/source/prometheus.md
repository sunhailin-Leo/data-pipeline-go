```
Prometheus
====
## Prometheus Source Configuration
prom_metrics is an object type used to configure parameters for fetching metrics from Prometheus. Here is the specific configuration description:

```json
{
  "source": [
    {
      "type": "Prometheus",
      "source_name": "prometheus-1",
      "prom_metrics": {
        "address": "<Prometheus Address>",   // Prometheus address, for example "http://prometheus.example.com:9090"
        "interval": <Interval in seconds>    // Query time interval (seconds), for example 30
      }
    }
  ]
}
```
## Configuration Parameter Description:
**type：**

Specifies the type of the source, here it is "Prometheus", indicating that Prometheus is used as the data source.

**source_name：**

A custom name for the Prometheus source, used to identify this source. In this example, it is set to "prometheus-1".

**prom_metrics：**
Contains the core configuration information for the Prometheus source:

* **address：**

    The address of the Prometheus instance, typically in HTTP format. For example, "http://prometheus.example.com:9090" specifies the API access address for Prometheus.

* **interval：**

    The time interval for querying Prometheus, in seconds. Defines the frequency of fetching the latest metric data from Prometheus. For example, 30 means querying Prometheus every 30 seconds.

## Configuration Example

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
In this example:

* "prometheus-1" is the defined Prometheus source name.

* "address" specifies the Prometheus instance address "http://prometheus.example.com:9090".

* "interval" is set to 30 seconds, indicating that data is fetched from Prometheus every 30 seconds.
```
