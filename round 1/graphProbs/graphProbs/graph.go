package graphProbs

type Node = string
type Weight = uint
type associations = map[Node]Weight
type adjacencyMatrix = map[Node]associations

type Graph struct {
	adjacencyMatrix adjacencyMatrix
	nodes           map[Node]bool
	edges           []Edge
}

type Edge struct {
	Frm Node
	To  Node
	Wt  Weight
}

func MkGraph(edges []Edge, nodes []Node) Graph {
	g := Graph{
		adjacencyMatrix: adjacencyMatrix{},
		nodes:           map[Node]bool{},
		edges:           append([]Edge{}, edges...),
	}
	for _, edge := range edges {
		g.AddEdge(edge)
	}
	for _, node := range nodes {
		g.AddNode(node)
	}
	return g
}

func (g *Graph) AddEdge(e Edge) {
	g.nodes[e.Frm] = true
	g.nodes[e.To] = true
	assocs, ok := g.adjacencyMatrix[e.Frm]
	if !ok {
		g.adjacencyMatrix[e.Frm] = associations{e.To: e.Wt}
		return
	}
	assocs[e.To] = e.Wt
}

func (g *Graph) AddNode(n Node) {
	g.nodes[n] = true
	_, ok := g.adjacencyMatrix[n]
	if !ok {
		g.adjacencyMatrix[n] = nil
	}
}

func (g *Graph) Neighbors(n Node) []Edge {
	assocs, ok := g.adjacencyMatrix[n]
	if !ok {
		return nil
	}
	neighbors := make([]Edge, len(assocs))
	i := 0
	for node := range assocs {
		neighbors[i].Frm = n
		neighbors[i].To = node
		neighbors[i].Wt = assocs[node]
		i += 1
	}
	return neighbors
}

func (g *Graph) IsValidNode(n Node) bool {
	_, ok := g.adjacencyMatrix[n]
	return ok
}

func (g *Graph) Nodes() []Node {
	nodes := make([]Node, len(g.nodes))
	i := 0
	for node := range g.nodes {
		nodes[i] = node
		i += 1
	}
	return nodes
}

func (g *Graph) Edges() []Edge {
	return append([]Edge{}, g.edges...)
}
