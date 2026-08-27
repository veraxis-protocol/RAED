package raed

import (
	"crypto/ed25519"
	"path/filepath"
	"testing"
	"time"
)

func id(t *testing.T, ns, ref, typ string) string {
	t.Helper()
	x, err := MintID(ns, ref, typ)
	if err != nil {
		t.Fatal(err)
	}
	return x.ID
}

func TestIdentityDeterministic(t *testing.T) {
	a, _ := MintID("concept", "x", "concept")
	b, _ := MintID("concept", "x", "concept")
	if a.ID != b.ID {
		t.Fatal("identity drift")
	}
}

func TestSemanticZeroAndConflict(t *testing.T) {
	s := NewStore()
	sub := id(t, "p", "P", "proposition")
	pred := id(t, "pred", "means", "predicate")
	obj1 := id(t, "o", "A", "concept")
	obj2 := id(t, "o", "B", "concept")
	ctx := id(t, "c", "C", "context")
	obs := id(t, "a", "observer", "actor")
	r1, _ := s.AppendRelation(Relation{SubjectID: sub, PredicateID: pred, ObjectID: obj1, ContextID: ctx, ObserverID: obs, EpistemicStatus: Verified, Dimensions: Dimensions{SemanticEpoch: "E1"}, PayloadDigest: "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"})
	if got := s.Resolve(Query{SubjectID: "missing"}); got.Status != SemanticZero {
		t.Fatalf("got %s", got.Status)
	}
	_, _ = r1, s
	_, err := s.AppendRelation(Relation{SubjectID: sub, PredicateID: pred, ObjectID: obj2, ContextID: ctx, ObserverID: obs, EpistemicStatus: Verified, Dimensions: Dimensions{SemanticEpoch: "E1"}, PayloadDigest: "sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"})
	if err != nil {
		t.Fatal(err)
	}
	if got := s.Resolve(Query{SubjectID: sub, ContextID: ctx}); got.Status != Conflict {
		t.Fatalf("expected conflict, got %s", got.Status)
	}
}

func TestRevocationPropagatesAndReplayStable(t *testing.T) {
	s := NewStore()
	ctx := id(t, "c", "C", "context")
	obs := id(t, "a", "O", "actor")
	pred := id(t, "p", "supports", "predicate")
	a := id(t, "x", "A", "proposition")
	b := id(t, "x", "B", "proposition")
	c := id(t, "x", "C", "proposition")
	r0, err := s.AppendRelation(Relation{SubjectID: a, PredicateID: pred, ObjectID: b, ContextID: ctx, ObserverID: obs, EpistemicStatus: Verified, Dimensions: Dimensions{SemanticEpoch: "E1"}, PayloadDigest: "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"})
	if err != nil {
		t.Fatal(err)
	}
	r1, err := s.AppendRelation(Relation{SubjectID: b, PredicateID: pred, ObjectID: c, ContextID: ctx, ObserverID: obs, EpistemicStatus: Derived, ParentRelationIDs: []string{r0.RelationID}, Dimensions: Dimensions{SemanticEpoch: "E1"}, PayloadDigest: "sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"})
	if err != nil {
		t.Fatal(err)
	}
	affected, err := s.AppendLifecycleEvent(LifecycleEvent{TargetID: r0.RelationID, Status: Revoked, ReasonCode: "TEST", EffectiveAt: time.Now().UTC()})
	if err != nil {
		t.Fatal(err)
	}
	if len(affected.AffectedIDs) != 1 || affected.AffectedIDs[0] != r1.RelationID {
		t.Fatalf("bad affected set: %+v", affected)
	}
	if rr, _ := s.Relation(r1.RelationID); rr.LifecycleStatus != ReevaluationRequired {
		t.Fatalf("child not invalidated: %s", rr.LifecycleStatus)
	}
	path := filepath.Join(t.TempDir(), "ledger.jsonl")
	if err := s.SaveLedger(path); err != nil {
		t.Fatal(err)
	}
	rp, err := ReplayLedger(path)
	if err != nil {
		t.Fatal(err)
	}
	if s.SnapshotDigest() != rp.SnapshotDigest() {
		t.Fatalf("snapshot mismatch %s != %s", s.SnapshotDigest(), rp.SnapshotDigest())
	}
}

func TestSignatureIntegrity(t *testing.T) {
	pub, priv, _ := ed25519.GenerateKey(nil)
	r := Relation{SubjectID: "a", PredicateID: "b", ObjectID: "c", ContextID: "d", ObserverID: "e", RecordedAt: time.Unix(0, 0).UTC(), EpistemicStatus: Unverified, LifecycleStatus: Active, PayloadDigest: "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}
	signed, err := SignRelation(r, priv)
	if err != nil {
		t.Fatal(err)
	}
	if !VerifyRelationSignature(signed, pub) {
		t.Fatal("signature should verify")
	}
	signed.ObjectID = "tampered"
	if VerifyRelationSignature(signed, pub) {
		t.Fatal("tamper not detected")
	}
}
