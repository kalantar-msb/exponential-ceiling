package flowcontrol

import "context"

// ConstantCeilingControl implements UsageLimitPolicy with a constant ceiling of 1.0
// for all bands. This is functionally identical to llm-d's built-in NewConstPolicy.
//
// Purpose: deploy via the same plugin interface as ExponentialCeilingPolicy to confirm
// the delivery mechanism introduces no behavioral difference. If this control produces
// results identical to default llm-d (no custom plugin), the plugin wiring is correct.
// If results differ, there's a framework bug to investigate before trusting treatment results.
type ConstantCeilingControl struct{}

func NewConstantCeilingControl() *ConstantCeilingControl {
	return &ConstantCeilingControl{}
}

func (p *ConstantCeilingControl) Name() string { return "constant-ceiling-control" }

// ComputeLimit returns ceiling = 1.0 for all bands regardless of saturation.
// This means dispatch is only blocked when saturation >= 1.0 — the same behavior
// as llm-d's default NewConstPolicy.
func (p *ConstantCeilingControl) ComputeLimit(
	ctx context.Context, saturation float64, priorities []int,
) []float64 {
	ceilings := make([]float64, len(priorities))
	for i := range ceilings {
		ceilings[i] = 1.0
	}
	return ceilings
}
