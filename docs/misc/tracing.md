# Distributed Tracing

`Janus` uses [OpenCensus](https://opencensus.io) as the standard way to trace requests. It can be used for monitoring microservices-based distributed systems:

- Distributed context propagation
- Distributed transaction monitoring
- Root cause analysis
- Service dependency analysis
- Performance / latency optimization

OpenCensus supports several tracing backend systems (i.e. [exporters](https://opencensus.io/exporters/supported-exporters/go/)) which are:
- Azure Monitor
- Honeycomb.io
- AWS X-Ray
- Datadog
- Jaeger
- Stackdriver
- Zipkin

Currently, only Jaeger exporter is available in `Janus`.

```toml
# Tracing Configuration

[tracing]
  # Backend system to export traces to
  #
  # Default: None
  #
  Exporter: "jaeger"
  
  # Service name used in the backend
  #
  # Default: "janus"
  #
  ServiceName: "janus"
  
  # SamplingStrategy specifies the sampling strategy
  #
  # Valid Values: "probabilistic", "always", "never"
  #
  # Default: "probabilistic"
  #
  SamplingStrategy: "probabilistic"
  
  # SamplingParam is an additional value passed to the sampler.
  #
  # Valid Values:
  #   - for "always" and "never" sampler, this value is unused
  #   - for "probabilistic" sampler, a probability between 0 and 1
  #
  # Default: "0.15"
  #
  SamplingParam: "0.15"
  
  [tracing.jaeger]
    # SamplingServerURL is the address to the sampling server
    #
    # Default: None
    #
    SamplingServerURL: "localhost:6832"
```
