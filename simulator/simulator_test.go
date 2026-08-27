package simulator

import "testing"

func TestDeterministicAffectedCount(t *testing.T) {
	a, err := Run(1000, 10000, 42)
	if err != nil {
		t.Fatal(err)
	}
	b, err := Run(1000, 10000, 42)
	if err != nil {
		t.Fatal(err)
	}
	if a.AffectedAgents != b.AffectedAgents {
		t.Fatalf("non-deterministic %d != %d", a.AffectedAgents, b.AffectedAgents)
	}
}
