package flowcontrol

import (
	"context"
	"math"
)

// ExponentialCeilingPolicy implements UsageLimitPolicy with parameter-free
// exponential per-band dispatch ceilings.
//
// Formula: ceiling[i] = exp(-N * sat * i / (N-1))
//
// The exponential shape creates widening separation between bands as saturation
// increases. At low load all bands dispatch freely (ceilings near 1.0). Under
// pressure, low-priority bands are held back exponentially harder than high-priority.
//
// Transfer to llm-d: register via StaticPolicyFactory pattern. No CLI flags needed.
// The formula derives behavior entirely from saturation and active band count.
type ExponentialCeilingPolicy struct{}

func NewExponentialCeilingPolicy() *ExponentialCeilingPolicy {
	return &ExponentialCeilingPolicy{}
}

func (p *ExponentialCeilingPolicy) Name() string { return "exponential-ceiling" }

// ComputeLimit returns per-band ceilings using exponential decay.
//
// ceiling[i] = exp(-N * saturation * i / (N-1))
//
// Properties:
//   - Band 0 (highest priority): always 1.0 (exp(0) = 1)
//   - Band N-1 (lowest priority): exp(-N * sat), drops fast with saturation
//   - Single band: returns [1.0] (no gating)
//   - sat=0: all ceilings = 1.0 (dormant, identical to constant policy)
//   - Stateless, goroutine-safe
func (p *ExponentialCeilingPolicy) ComputeLimit(
	ctx context.Context, saturation float64, priorities []int,
) []float64 {
	n := len(priorities)
	ceilings := make([]float64, n)

	if n <= 1 {
		if n == 1 {
			ceilings[0] = 1.0
		}
		return ceilings
	}

	for i := range priorities {
		ceilings[i] = math.Exp(-float64(n) * saturation * float64(i) / float64(n-1))
	}
	return ceilings
}
