package trace

import (
	. "go.opencensus.io/trace"
)

// RespectParentSampler always uses decision made by parent span.
// If there is no parent span, decide with the a specified fallbackSampler.
func RespectParentSampler(fallbackSampler Sampler) Sampler {
	return func(p SamplingParameters) SamplingDecision {
		if (p.ParentContext != SpanContext{}) {
			return SamplingDecision{Sample: p.ParentContext.IsSampled()}
		}

		return fallbackSampler(p)
	}
}
