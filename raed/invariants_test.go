package raed

import (
	"testing"
	"time"
)

func TestContextChangesRelationIdentity(t *testing.T) {
	s := NewStore()
	sub := id(t, "p", "P", "proposition")
	pred := id(t, "p", "rel", "predicate")
	obj := id(t, "o", "O", "concept")
	obs := id(t, "a", "observer", "actor")
	c1 := id(t, "c", "C1", "context")
	c2 := id(t, "c", "C2", "context")
	base := Relation{SubjectID: sub, PredicateID: pred, ObjectID: obj, ObserverID: obs, RecordedAt: time.Unix(1, 0).UTC(), EpistemicStatus: Verified, LifecycleStatus: Active, PayloadDigest: "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}
	r1, err := s.AppendRelation(func() Relation { r := base; r.ContextID = c1; return r }())
	if err != nil {
		t.Fatal(err)
	}
	r2, err := s.AppendRelation(func() Relation { r := base; r.ContextID = c2; return r }())
	if err != nil {
		t.Fatal(err)
	}
	if r1.RelationID == r2.RelationID {
		t.Fatal("context mutation must change relation identity")
	}
}

func TestGuardLossChangesConstraintDigest(t *testing.T) {
	guard := id(t, "g", "G", "proposition")
	d1 := Dimensions{Modality: "MAY", Guards: []string{guard}, SemanticEpoch: "E1"}
	d2 := Dimensions{Modality: "MAY", SemanticEpoch: "E1"}
	b1, _ := canonicalJSON(d1)
	b2, _ := canonicalJSON(d2)
	if digestBytes(b1) == digestBytes(b2) {
		t.Fatal("guard loss became invisible")
	}
}

func TestInactiveRelationCannotExportOrdinaryEnvelope(t *testing.T) {
	s := NewStore()
	sub := id(t, "p", "P", "proposition")
	pred := id(t, "p", "rel", "predicate")
	obj := id(t, "o", "O", "concept")
	ctx := id(t, "c", "C", "context")
	obs := id(t, "a", "observer", "actor")
	r, err := s.AppendRelation(Relation{SubjectID: sub, PredicateID: pred, ObjectID: obj, ContextID: ctx, ObserverID: obs, RecordedAt: time.Unix(1, 0).UTC(), EpistemicStatus: Verified, LifecycleStatus: Active, Dimensions: Dimensions{SemanticEpoch: "E1"}, PayloadDigest: "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.AppendLifecycleEvent(LifecycleEvent{TargetID: r.RelationID, Status: Revoked, ReasonCode: "TEST", EffectiveAt: time.Unix(2, 0).UTC()}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.ToMemoryEnvelope([]string{r.RelationID}, "test"); err == nil {
		t.Fatal("revoked relation exported as ordinary active memory")
	}
}

func TestResolverLimitIsVisible(t *testing.T) {
	s := NewStore()
	sub := id(t, "p", "P", "proposition")
	pred := id(t, "p", "rel", "predicate")
	ctx := id(t, "c", "C", "context")
	obs := id(t, "a", "observer", "actor")
	for i, ref := range []string{"O1", "O2"} {
		obj := id(t, "o", ref, "concept")
		_, err := s.AppendRelation(Relation{SubjectID: sub, PredicateID: pred, ObjectID: obj, ContextID: ctx, ObserverID: obs, RecordedAt: time.Unix(int64(i+1), 0).UTC(), EpistemicStatus: Verified, LifecycleStatus: Active, PayloadDigest: "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"})
		if err != nil {
			t.Fatal(err)
		}
	}
	got := s.Resolve(Query{SubjectID: sub, MaxRecords: 1})
	if got.Status != LimitExceeded {
		t.Fatalf("silent truncation: %s", got.Status)
	}
}
