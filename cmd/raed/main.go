package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/veraxis-protocol/RAED/raed"
	"github.com/veraxis-protocol/RAED/simulator"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	switch os.Args[1] {
	case "id":
		fs := flag.NewFlagSet("id", flag.ExitOnError)
		ns := fs.String("namespace", "concept", "namespace")
		ref := fs.String("ref", "", "normalized reference")
		typ := fs.String("type", "concept", "object type")
		_ = fs.Parse(os.Args[2:])
		x, err := raed.MintID(*ns, *ref, *typ)
		exitJSON(x, err)
	case "scale":
		fs := flag.NewFlagSet("scale", flag.ExitOnError)
		agents := fs.Int("agents", 100000, "agent count")
		edges := fs.Int("connections", 1000000, "directed connection count")
		seed := fs.Int64("seed", 42, "deterministic seed")
		_ = fs.Parse(os.Args[2:])
		r, err := simulator.Run(*agents, *edges, *seed)
		exitJSON(r, err)
	case "demo":
		demo()
	default:
		usage()
		os.Exit(2)
	}
}
func exitJSON(v any, err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}
func usage() { fmt.Println("raed <id|demo|scale>") }
func demo() {
	s := raed.NewStore()
	mint := func(ns, ref, typ string) string { x, _ := raed.MintID(ns, ref, typ); return x.ID }
	ctx := mint("context", "example-context", "context")
	obs := mint("actor", "observer", "actor")
	pred := mint("predicate", "permits", "predicate")
	sub := mint("proposition", "state-P", "proposition")
	obj := mint("action", "action-A", "action")
	guard := mint("proposition", "guard-G", "proposition")
	r, err := s.AppendRelation(raed.Relation{
		SubjectID: sub, PredicateID: pred, ObjectID: obj, ContextID: ctx, ObserverID: obs,
		RecordedAt: time.Now().UTC(), EpistemicStatus: raed.Verified, LifecycleStatus: raed.Active,
		Dimensions:         raed.Dimensions{Modality: "MAY", Guards: []string{guard}, SemanticEpoch: "E1"},
		SourceDependencies: []string{guard},
		PayloadDigest:      "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	})
	if err != nil {
		exitJSON(nil, err)
	}
	rec := s.Resolve(raed.Query{SubjectID: sub, ContextID: ctx})
	exitJSON(map[string]any{"relation": r, "resolution": rec}, nil)
}
