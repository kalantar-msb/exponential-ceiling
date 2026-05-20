package exponentialceiling

import (
	"context"
	"math"
	"testing"
)

func TestComputeLimit_SingleBand(t *testing.T) {
	ceilings := computeLimit(context.Background(), 0.5, []int{100})
	if len(ceilings) != 1 {
		t.Fatalf("expected 1 ceiling, got %d", len(ceilings))
	}
	if ceilings[0] != 1.0 {
		t.Errorf("single band ceiling should be 1.0, got %f", ceilings[0])
	}
}

func TestComputeLimit_EmptyPriorities(t *testing.T) {
	ceilings := computeLimit(context.Background(), 0.5, []int{})
	if len(ceilings) != 0 {
		t.Fatalf("expected 0 ceilings, got %d", len(ceilings))
	}
}

func TestComputeLimit_ZeroSaturation(t *testing.T) {
	priorities := []int{100, 50, 10}
	ceilings := computeLimit(context.Background(), 0.0, priorities)

	for i, c := range ceilings {
		if c != 1.0 {
			t.Errorf("at zero saturation, ceiling[%d] should be 1.0, got %f", i, c)
		}
	}
}

func TestComputeLimit_HighestPriorityAlwaysOne(t *testing.T) {
	priorities := []int{100, 50, 10}
	saturations := []float64{0.0, 0.25, 0.5, 0.75, 1.0}

	for _, sat := range saturations {
		ceilings := computeLimit(context.Background(), sat, priorities)
		if ceilings[0] != 1.0 {
			t.Errorf("sat=%f: highest priority ceiling should be 1.0, got %f", sat, ceilings[0])
		}
	}
}

func TestComputeLimit_MonotonicallyDecreasing(t *testing.T) {
	priorities := []int{100, 75, 50, 25, 10}
	ceilings := computeLimit(context.Background(), 0.6, priorities)

	for i := 1; i < len(ceilings); i++ {
		if ceilings[i] > ceilings[i-1] {
			t.Errorf("ceilings not monotonically decreasing: ceiling[%d]=%f > ceiling[%d]=%f",
				i, ceilings[i], i-1, ceilings[i-1])
		}
	}
}

func TestComputeLimit_Formula(t *testing.T) {
	// Verify the formula: ceiling[i] = exp(-N * sat * i / (N-1))
	priorities := []int{100, 50, 10}
	sat := 0.5
	n := float64(len(priorities))

	ceilings := computeLimit(context.Background(), sat, priorities)

	for i := range priorities {
		expected := math.Exp(-n * sat * float64(i) / (n - 1))
		if math.Abs(ceilings[i]-expected) > 1e-12 {
			t.Errorf("ceiling[%d]: expected %f, got %f", i, expected, ceilings[i])
		}
	}
}

func TestComputeLimit_FullSaturation(t *testing.T) {
	priorities := []int{100, 50, 10}
	ceilings := computeLimit(context.Background(), 1.0, priorities)

	// Band 0 should still be 1.0
	if ceilings[0] != 1.0 {
		t.Errorf("full saturation: ceiling[0] should be 1.0, got %f", ceilings[0])
	}
	// Last band should be exp(-N) = exp(-3) ≈ 0.0498
	expected := math.Exp(-float64(len(priorities)))
	if math.Abs(ceilings[len(ceilings)-1]-expected) > 1e-12 {
		t.Errorf("full saturation: ceiling[N-1] expected %f, got %f", expected, ceilings[len(ceilings)-1])
	}
}

func TestFactory(t *testing.T) {
	p, err := Factory("test-policy", nil, nil)
	if err != nil {
		t.Fatalf("Factory returned error: %v", err)
	}
	if p == nil {
		t.Fatal("Factory returned nil plugin")
	}

	tn := p.TypedName()
	if tn.Name != "test-policy" {
		t.Errorf("expected name 'test-policy', got %q", tn.Name)
	}
}

func TestPolicyType(t *testing.T) {
	if PolicyType != "exponential-ceiling-policy" {
		t.Errorf("PolicyType = %q, want %q", PolicyType, "exponential-ceiling-policy")
	}
}
