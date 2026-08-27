package raed

import "testing"

func FuzzMintIDNeverPanics(f *testing.F) {
	f.Add("concept", "state-P", "proposition")
	f.Fuzz(func(t *testing.T, ns, ref, typ string) {
		_, _ = MintID(ns, ref, typ)
	})
}

func FuzzCanonicalDimensionsNeverPanics(f *testing.F) {
	f.Add("MAY", "positive", "E1")
	f.Fuzz(func(t *testing.T, modality, polarity, epoch string) {
		_, _ = canonicalJSON(Dimensions{Modality: modality, Polarity: polarity, SemanticEpoch: epoch})
	})
}
