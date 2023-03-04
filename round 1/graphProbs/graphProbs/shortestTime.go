package graphProbs

import (
	"container/heap"
	"math"
)

type pair struct {
	node   Node
	weight Weight
}

// An pairHeap is a min-heap of Weights.
type pairHeap []pair

func (h pairHeap) Len() int           { return len(h) }
func (h pairHeap) Less(i, j int) bool { return h[i].weight < h[j].weight }
func (h pairHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *pairHeap) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(pair))
}

func (h *pairHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (g *Graph) ShortestTime(start Node, end Node) *Weight {
	if !g.IsValidNode(start) || !g.IsValidNode(end) {
		return nil
	}
	dists := map[Node]Weight{}
	pq := make(pairHeap, len(g.nodes))
	pq[0] = pair{node: start, weight: 0}
	dists[start] = 0
	i := 1
	for node := range g.nodes {
		if node == start {
			continue
		}
		dists[node] = uint(math.Inf(0))
		pq[i] = pair{node, uint(math.Inf(0))}
		i += 1
	}
	visited := map[Node]bool{}
	for {
		if len(pq) == 0 {
			break
		}
		pr := heap.Pop(&pq).(pair)
		if pr.node == end {
			return &pr.weight
		}
		visited[pr.node] = true
		du := pr.weight
		for _, neighbor := range g.Neighbors(pr.node) {
			vnode := neighbor.To
			wt := neighbor.Wt
			_, ok := visited[vnode]
			if ok {
				continue
			}
			dv := dists[vnode]
			if du+wt <= dv {
				dists[vnode] = du + wt
				heap.Push(&pq, pair{node: vnode, weight: du + wt})
			}
		}
	}
	return nil
}
