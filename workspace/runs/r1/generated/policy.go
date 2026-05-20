// Package exponentialceiling implements a UsageLimitPolicy that applies
// exponential per-band dispatch ceilings based on pool-wide saturation.
//
// Formula: ceiling[i] = exp(-N * saturation * i / (N-1))
//
// Properties:
//   - Band 0 (highest priority): always 1.0 (never gated)
//   - Band N-1 (lowest priority): exp(-N * sat), drops fast with saturation
//   - Single band: returns [1.0] (no gating needed)
//   - sat=0: all ceilings = 1.0 (dormant, identical to constant policy)
//   - Stateless and goroutine-safe
package exponentialceiling

import (
	"context"
	"encoding/json"
	"math"

	"github.com/llm-d/llm-d-router/pkg/epp/framework/interface/flowcontrol"
	"github.com/llm-d/llm-d-router/pkg/epp/framework/interface/plugin"
	"github.com/llm-d/llm-d-router/pkg/epp/framework/plugins/flowcontrol/usagelimits"
)

// PolicyType is the registered plugin type for this policy.
const PolicyType = "exponential-ceiling-policy"

// Factory creates an instance of the exponential ceiling usage limit policy.
// The policy is parameter-free; rawConfig is ignored.
func Factory(name string, _ json.RawMessage, _ plugin.Handle) (plugin.Plugin, error) {
	return newPolicy(name), nil
}

// newPolicy returns a UsageLimitPolicy using the exponential ceiling formula.
func newPolicy(name string) flowcontrol.UsageLimitPolicy {
	return usagelimits.NewPolicyFunc(name, computeLimit)
}

// computeLimit calculates per-band ceilings using exponential decay.
//
// ceiling[i] = exp(-N * saturation * i / (N-1))
func computeLimit(_ context.Context, saturation float64, priorities []int) []float64 {
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
