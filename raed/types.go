package raed

import "time"

type EpistemicStatus string
type LifecycleStatus string
type ResolutionStatus string

const (
	Unverified  EpistemicStatus = "UNVERIFIED"
	Supported   EpistemicStatus = "SUPPORTED"
	Verified    EpistemicStatus = "VERIFIED"
	Derived     EpistemicStatus = "DERIVED"
	Conflicting EpistemicStatus = "CONFLICTING"

	Active               LifecycleStatus = "ACTIVE"
	Expired              LifecycleStatus = "EXPIRED"
	Revoked              LifecycleStatus = "REVOKED"
	Superseded           LifecycleStatus = "SUPERSEDED"
	Quarantined          LifecycleStatus = "QUARANTINED"
	ReevaluationRequired LifecycleStatus = "REEVALUATION_REQUIRED"

	Resolved      ResolutionStatus = "RESOLVED"
	Conflict      ResolutionStatus = "CONFLICTING"
	SemanticZero  ResolutionStatus = "SEMANTIC_ZERO"
	LimitExceeded ResolutionStatus = "LIMIT_EXCEEDED"
)

type Identity struct {
	ID            string    `json:"id"`
	Namespace     string    `json:"namespace"`
	NormalizedRef string    `json:"normalized_ref"`
	ObjectType    string    `json:"object_type"`
	CreatedAt     time.Time `json:"created_at"`
}

type Dimensions struct {
	ActorID        string   `json:"actor_id,omitempty"`
	ActionID       string   `json:"action_id,omitempty"`
	ObjectID       string   `json:"object_id,omitempty"`
	Modality       string   `json:"modality,omitempty"`
	Polarity       string   `json:"polarity,omitempty"`
	Guards         []string `json:"guards,omitempty"`
	Exceptions     []string `json:"exceptions,omitempty"`
	ScopeIDs       []string `json:"scope_ids,omitempty"`
	TemporalClass  string   `json:"temporal_class,omitempty"`
	NumericBounds  []string `json:"numeric_bounds,omitempty"`
	Quantification string   `json:"quantification,omitempty"`
	EpistemicClass string   `json:"epistemic_class,omitempty"`
	AuthorityClass string   `json:"authority_class,omitempty"`
	SemanticEpoch  string   `json:"semantic_epoch,omitempty"`
}

type Relation struct {
	RelationID            string          `json:"relation_id"`
	SubjectID             string          `json:"subject_id"`
	PredicateID           string          `json:"predicate_id"`
	ObjectID              string          `json:"object_id"`
	ContextID             string          `json:"context_id"`
	ObserverID            string          `json:"observer_id"`
	OccurredAt            *time.Time      `json:"occurred_at,omitempty"`
	RecordedAt            time.Time       `json:"recorded_at"`
	EpistemicStatus       EpistemicStatus `json:"epistemic_status"`
	LifecycleStatus       LifecycleStatus `json:"lifecycle_status"`
	Dimensions            Dimensions      `json:"semantic_dimensions"`
	SourceDependencies    []string        `json:"source_dependencies,omitempty"`
	ParentRelationIDs     []string        `json:"parent_relation_ids,omitempty"`
	ConstraintDigest      string          `json:"constraint_digest,omitempty"`
	VerificationProfileID string          `json:"verification_profile_id,omitempty"`
	PayloadDigest         string          `json:"payload_digest"`
	SigningKeyID          string          `json:"signing_key_id,omitempty"`
	Signature             string          `json:"signature,omitempty"`
}

type Query struct {
	SubjectID       string
	PredicateID     string
	ContextID       string
	SemanticEpoch   string
	MaxRecords      int
	IncludeInactive bool
}

type ResolutionReceipt struct {
	ReceiptID            string           `json:"receipt_id"`
	QueryDigest          string           `json:"query_digest"`
	Status               ResolutionStatus `json:"status"`
	ResolvedRelationIDs  []string         `json:"resolved_relation_ids,omitempty"`
	ConsultedRelationIDs []string         `json:"consulted_relation_ids,omitempty"`
	SnapshotDigest       string           `json:"raed_snapshot_digest"`
	IssuedAt             time.Time        `json:"issued_at"`
}

type LifecycleEvent struct {
	TargetID    string          `json:"target_id"`
	Status      LifecycleStatus `json:"status"`
	ReasonCode  string          `json:"reason_code"`
	EffectiveAt time.Time       `json:"effective_at"`
}

type AffectedSet struct {
	RootID       string   `json:"root_id"`
	AffectedIDs  []string `json:"affected_ids"`
	VisitedCount int      `json:"visited_count"`
}

type MemoryEnvelope struct {
	MemoryID              string          `json:"memory_id"`
	RAEDObjectID          string          `json:"raed_object_id"`
	RelationIDs           []string        `json:"relation_ids"`
	SemanticPayloadDigest string          `json:"semantic_payload_digest"`
	ContextID             string          `json:"context_id"`
	Purpose               string          `json:"purpose,omitempty"`
	SemanticEpoch         string          `json:"semantic_epoch"`
	CreatedAt             time.Time       `json:"created_at"`
	ExpiresAt             *time.Time      `json:"expires_at,omitempty"`
	EpistemicStatus       EpistemicStatus `json:"epistemic_status"`
	LifecycleStatus       LifecycleStatus `json:"lifecycle_status"`
	SourceDependencyIDs   []string        `json:"source_dependency_ids,omitempty"`
	ParentRelationIDs     []string        `json:"parent_relation_ids,omitempty"`
	ConstraintDigest      string          `json:"constraint_digest,omitempty"`
	SnapshotDigest        string          `json:"raed_snapshot_digest"`
}
