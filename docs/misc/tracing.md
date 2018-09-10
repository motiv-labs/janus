# Distributed Tracing

Janus use [opentracing](http://opentracing.io/) as the standard way to trace requests. It can be used for monitoring microservices-based distributed systems:

- Distributed context propagation
- Distributed transaction monitoring
- Root cause analysis
- Service dependency analysis
- Performance / latency optimization

We support some [tracers](https://github.com/opentracing/specification/blob/master/specification.md#tracer) which are:

### Jaeger

[Jaeger](https://github.com/jaegertracing/jaeger) \ˈyā-gər\ is a distributed tracing system, originally open sourced by Uber Technologies. 
It provides distributed context propagation, distributed transaction monitoring, root cause analysis, service dependency analysis, 
and performance / latency optimization. Built with OpenTracing support from inception, Jaeger includes OpenTracing client libraries in 
several languages, including Java, Go, Python, Node.js, and C++. It is a Cloud Native Computing Foundation member project.

You can define that you want to use Jaeger on your configuration file under the `tracing` section:

```toml
# Jeager Distributed Tracing
[tracing]
  # Backend name used to send tracing data
  #
  # Default: "noop"
  #
  provider = "jaeger"

  # Service name used in the backend
  #
  # Default: "janus"
  #
  serviceName = "janus"

# Jeager Distributed Tracing
  [tracing.jaeger]
    # SamplingServerURL is the address of jaeger-agent's HTTP sampling server
    #
    # Default: "http://localhost:5778/sampling"
    #
    SamplingServerURL = "localhost:6832"

    # Sampling Type specifies the type of the sampler: const, probabilistic, rateLimiting
    #
    # Default: "const"
    #
    # SamplingType = "const"

    # SamplingParam Param is a value passed to the sampler.
    # Valid values for Param field are:
    #   - for "const" sampler, 0 or 1 for always false/true respectively
    #   - for "probabilistic" sampler, a probability between 0 and 1
    #   - for "rateLimiting" sampler, the number of spans per second
    #
    # Default: 1.0
    #
    # SamplingParam = 1.0

    # BufferFlushInterval controls how often the buffer is force-flushed, even if it's not full.
    # It is generally not useful, as it only matters for very low traffic services.
    # 
    # Default: 1s
    #
    # BufferFlushInterval = "1s"

    # LogSpans, when true, enables LoggingReporter that runs in parallel with the main reporter
    # and logs all submitted spans. Main Configuration.Logger must be initialized in the code
    # for this option to have any effect.
    # 
    # Default: false
    #
    # LogSpans = false

    # QueueSize controls how many spans the reporter can keep in memory before it starts dropping
	  # new spans. The queue is continuously drained by a background go-routine, as fast as spans
	  # can be sent out of process.
    # 
    # Default: 0
    #
    # QueueSize = 0

    # PropagationFormat is the propagation format jaeger will use.
    # Leave it blank to use jaeger default propagation mechanism
    # 
    # PropagationFormat = "zipkin"

```

### Google Cloud Platform - Tracing

[Stackdriver Trace](https://github.com/hellofresh/gcloud-opentracing) is a distributed tracing system that collects latency data from your applications and displays it in 
the Google Cloud Platform Console. You can track how requests propagate through your application and receive detailed near real-time performance insights. Stackdriver Trace 
automatically analyzes all of your application's traces to generate in-depth latency reports to surface performance degradations, and can capture traces from all of your VMs, 
containers, or Google App Engine projects.

You can define that you want to use Stackdriver Trace on your configuration file under the `tracing` section:

```toml
# Google Clooud Platform - Stackdriver Trace
[tracing]
  # Backend name used to send tracing data
  #
  # Default: "noop"
  #
  provider = "googleCloud"

  # Service name used in the backend
  #
  # Default: "janus"
  #
  serviceName = "janus"

  [tracing.googleCloud]
    projectID = ""
    email = ""
    privateKey = ""
    privateKeyID = ""
```
