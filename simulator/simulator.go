package simulator

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

type Edge struct{ From, To uint32 }

type Result struct {
	Agents             int    `json:"agents"`
	Connections        int    `json:"connections"`
	Seed               int64  `json:"seed"`
	GenerationTimeNS   int64  `json:"generation_time_ns"`
	InvalidationTimeNS int64  `json:"invalidation_time_ns"`
	AffectedAgents     int    `json:"affected_agents"`
	HeapAllocBytes     uint64 `json:"heap_alloc_bytes"`
	TotalAllocBytes    uint64 `json:"total_alloc_bytes"`
}

func Run(agents, connections int, seed int64) (Result, error) {
	if agents <= 0 || connections < 0 {
		return Result{}, fmt.Errorf("invalid scale")
	}
	start := time.Now()
	rng := rand.New(rand.NewSource(seed))
	edges := make([]Edge, connections)
	children := make([][]uint32, agents)
	for i := 0; i < connections; i++ {
		from := uint32(rng.Intn(agents))
		to := uint32(rng.Intn(agents))
		edges[i] = Edge{from, to}
		children[from] = append(children[from], to)
	}
	gen := time.Since(start)
	_ = edges // retain through invalidation measurement
	invStart := time.Now()
	seen := make([]bool, agents)
	queue := make([]uint32, 0, agents)
	seen[0] = true
	queue = append(queue, 0)
	affected := 0
	for head := 0; head < len(queue); head++ {
		cur := queue[head]
		for _, n := range children[cur] {
			if seen[n] {
				continue
			}
			seen[n] = true
			affected++
			queue = append(queue, n)
		}
	}
	inv := time.Since(invStart)
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return Result{Agents: agents, Connections: connections, Seed: seed, GenerationTimeNS: gen.Nanoseconds(), InvalidationTimeNS: inv.Nanoseconds(), AffectedAgents: affected, HeapAllocBytes: m.HeapAlloc, TotalAllocBytes: m.TotalAlloc}, nil
}
