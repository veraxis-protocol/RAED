package raed

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"
)

type ledgerRecord struct {
	Kind      string          `json:"kind"`
	Relation  *Relation       `json:"relation,omitempty"`
	Lifecycle *LifecycleEvent `json:"lifecycle,omitempty"`
}

type Store struct {
	mu         sync.RWMutex
	relations  map[string]Relation
	bySubject  map[string][]string
	dependents map[string][]string
	events     []ledgerRecord
}

func NewStore() *Store {
	return &Store{
		relations:  make(map[string]Relation),
		bySubject:  make(map[string][]string),
		dependents: make(map[string][]string),
	}
}

func (s *Store) AppendRelation(r Relation) (Relation, error) {
	if r.SubjectID == "" || r.PredicateID == "" || r.ObjectID == "" || r.ContextID == "" || r.ObserverID == "" {
		return Relation{}, errors.New("subject, predicate, object, context and observer are required")
	}
	if r.RecordedAt.IsZero() {
		r.RecordedAt = time.Now().UTC()
	}
	if r.EpistemicStatus == "" {
		r.EpistemicStatus = Unverified
	}
	if r.LifecycleStatus == "" {
		r.LifecycleStatus = Active
	}
	id, err := relationDigest(r)
	if err != nil {
		return Relation{}, err
	}
	r.RelationID = id

	s.mu.Lock()
	defer s.mu.Unlock()
	if existing, ok := s.relations[id]; ok {
		return existing, nil
	}
	s.relations[id] = r
	s.bySubject[r.SubjectID] = append(s.bySubject[r.SubjectID], id)
	for _, dep := range append(append([]string{}, r.SourceDependencies...), r.ParentRelationIDs...) {
		s.dependents[dep] = append(s.dependents[dep], id)
	}
	cp := r
	s.events = append(s.events, ledgerRecord{Kind: "relation", Relation: &cp})
	return r, nil
}

func (s *Store) Relation(id string) (Relation, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.relations[id]
	return r, ok
}

func (s *Store) Resolve(q Query) ResolutionReceipt {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if q.MaxRecords <= 0 {
		q.MaxRecords = 256
	}
	qb, _ := json.Marshal(q)
	ids := append([]string(nil), s.bySubject[q.SubjectID]...)
	sort.Strings(ids)
	consulted := make([]string, 0, min(len(ids), q.MaxRecords))
	matches := make([]Relation, 0)
	for _, id := range ids {
		if len(consulted) >= q.MaxRecords {
			return s.receiptLocked(qb, LimitExceeded, nil, consulted)
		}
		r := s.relations[id]
		consulted = append(consulted, id)
		if !q.IncludeInactive && r.LifecycleStatus != Active {
			continue
		}
		if q.PredicateID != "" && r.PredicateID != q.PredicateID {
			continue
		}
		if q.ContextID != "" && r.ContextID != q.ContextID {
			continue
		}
		if q.SemanticEpoch != "" && r.Dimensions.SemanticEpoch != q.SemanticEpoch {
			continue
		}
		matches = append(matches, r)
	}
	if len(matches) == 0 {
		return s.receiptLocked(qb, SemanticZero, nil, consulted)
	}
	conflict := false
	base := matches[0]
	for _, r := range matches[1:] {
		if r.ObjectID != base.ObjectID || r.Dimensions.Polarity != base.Dimensions.Polarity || r.Dimensions.Modality != base.Dimensions.Modality {
			conflict = true
			break
		}
	}
	resolved := make([]string, len(matches))
	for i, r := range matches {
		resolved[i] = r.RelationID
	}
	status := Resolved
	if conflict {
		status = Conflict
	}
	return s.receiptLocked(qb, status, resolved, consulted)
}

func (s *Store) receiptLocked(qb []byte, status ResolutionStatus, resolved, consulted []string) ResolutionReceipt {
	snap := s.snapshotLocked()
	body := struct {
		Query     string           `json:"query"`
		Status    ResolutionStatus `json:"status"`
		Resolved  []string         `json:"resolved"`
		Consulted []string         `json:"consulted"`
		Snapshot  string           `json:"snapshot"`
	}{digestBytes(qb), status, sortedCopy(resolved), sortedCopy(consulted), snap}
	b, _ := json.Marshal(body)
	return ResolutionReceipt{
		ReceiptID: digestBytes(b), QueryDigest: digestBytes(qb), Status: status,
		ResolvedRelationIDs: sortedCopy(resolved), ConsultedRelationIDs: sortedCopy(consulted),
		SnapshotDigest: snap, IssuedAt: time.Now().UTC(),
	}
}

func (s *Store) AppendLifecycleEvent(e LifecycleEvent) (AffectedSet, error) {
	if e.TargetID == "" || e.Status == "" {
		return AffectedSet{}, errors.New("target and status required")
	}
	if e.EffectiveAt.IsZero() {
		e.EffectiveAt = time.Now().UTC()
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	r, ok := s.relations[e.TargetID]
	if !ok {
		return AffectedSet{}, fmt.Errorf("target not found: %s", e.TargetID)
	}
	r.LifecycleStatus = e.Status
	s.relations[e.TargetID] = r
	cp := e
	s.events = append(s.events, ledgerRecord{Kind: "lifecycle", Lifecycle: &cp})
	affected := s.affectedLocked(e.TargetID)
	for _, id := range affected.AffectedIDs {
		if rr, ok := s.relations[id]; ok && rr.LifecycleStatus == Active {
			rr.LifecycleStatus = ReevaluationRequired
			s.relations[id] = rr
		}
	}
	return affected, nil
}

func (s *Store) AffectedBy(id string) AffectedSet {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.affectedLocked(id)
}

func (s *Store) affectedLocked(root string) AffectedSet {
	seen := map[string]bool{root: true}
	queue := []string{root}
	out := []string{}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		children := append([]string(nil), s.dependents[cur]...)
		sort.Strings(children)
		for _, c := range children {
			if seen[c] {
				continue
			}
			seen[c] = true
			out = append(out, c)
			queue = append(queue, c)
		}
	}
	sort.Strings(out)
	return AffectedSet{RootID: root, AffectedIDs: out, VisitedCount: len(seen)}
}

func (s *Store) SnapshotDigest() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.snapshotLocked()
}

func (s *Store) snapshotLocked() string {
	ids := make([]string, 0, len(s.relations))
	for id := range s.relations {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	rows := make([]struct {
		ID   string
		Life LifecycleStatus
		Epi  EpistemicStatus
	}, 0, len(ids))
	for _, id := range ids {
		r := s.relations[id]
		rows = append(rows, struct {
			ID   string
			Life LifecycleStatus
			Epi  EpistemicStatus
		}{id, r.LifecycleStatus, r.EpistemicStatus})
	}
	b, _ := json.Marshal(rows)
	return digestBytes(b)
}

func (s *Store) SaveLedger(path string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	for _, e := range s.events {
		if err := enc.Encode(e); err != nil {
			return err
		}
	}
	return f.Sync()
}

func ReplayLedger(path string) (*Store, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	s := NewStore()
	for dec.More() {
		var rec ledgerRecord
		if err := dec.Decode(&rec); err != nil {
			return nil, err
		}
		switch rec.Kind {
		case "relation":
			if rec.Relation == nil {
				return nil, errors.New("malformed relation record")
			}
			r := *rec.Relation
			origID := r.RelationID
			r.RelationID = ""
			added, err := s.AppendRelation(r)
			if err != nil {
				return nil, err
			}
			if added.RelationID != origID {
				return nil, errors.New("replay digest mismatch")
			}
		case "lifecycle":
			if rec.Lifecycle == nil {
				return nil, errors.New("malformed lifecycle record")
			}
			if _, err := s.AppendLifecycleEvent(*rec.Lifecycle); err != nil {
				return nil, err
			}
		default:
			return nil, errors.New("unknown ledger record")
		}
	}
	return s, nil
}

func SignRelation(r Relation, private ed25519.PrivateKey) (Relation, error) {
	cp := r
	cp.Signature = ""
	b, err := canonicalJSON(cp)
	if err != nil {
		return Relation{}, err
	}
	cp.Signature = base64.StdEncoding.EncodeToString(ed25519.Sign(private, b))
	return cp, nil
}

func VerifyRelationSignature(r Relation, public ed25519.PublicKey) bool {
	sig, err := base64.StdEncoding.DecodeString(r.Signature)
	if err != nil {
		return false
	}
	cp := r
	cp.Signature = ""
	b, err := canonicalJSON(cp)
	if err != nil {
		return false
	}
	return ed25519.Verify(public, b, sig)
}

func (s *Store) ToMemoryEnvelope(relationIDs []string, purpose string) (MemoryEnvelope, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(relationIDs) == 0 {
		return MemoryEnvelope{}, errors.New("at least one relation required")
	}
	sort.Strings(relationIDs)
	var first Relation
	deps := []string{}
	parents := []string{}
	for i, id := range relationIDs {
		r, ok := s.relations[id]
		if !ok {
			return MemoryEnvelope{}, fmt.Errorf("relation not found: %s", id)
		}
		if r.LifecycleStatus != Active {
			return MemoryEnvelope{}, fmt.Errorf("relation not active: %s (%s)", id, r.LifecycleStatus)
		}
		if i == 0 {
			first = r
		} else if r.ContextID != first.ContextID || r.Dimensions.SemanticEpoch != first.Dimensions.SemanticEpoch {
			return MemoryEnvelope{}, errors.New("mixed context or semantic epoch")
		}
		deps = append(deps, r.SourceDependencies...)
		parents = append(parents, r.ParentRelationIDs...)
	}
	payload, _ := json.Marshal(relationIDs)
	memID := digestBytes(append([]byte("raed:memory:v0.1\x00"), payload...))
	return MemoryEnvelope{
		MemoryID: memID, RAEDObjectID: first.SubjectID, RelationIDs: relationIDs,
		SemanticPayloadDigest: digestBytes(payload), ContextID: first.ContextID, Purpose: purpose,
		SemanticEpoch: first.Dimensions.SemanticEpoch, CreatedAt: time.Now().UTC(),
		EpistemicStatus: first.EpistemicStatus, LifecycleStatus: first.LifecycleStatus,
		SourceDependencyIDs: sortedCopy(deps), ParentRelationIDs: sortedCopy(parents),
		ConstraintDigest: first.ConstraintDigest, SnapshotDigest: s.snapshotLocked(),
	}, nil
}
