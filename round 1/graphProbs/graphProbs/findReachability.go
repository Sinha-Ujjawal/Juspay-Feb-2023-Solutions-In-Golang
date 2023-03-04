package graphProbs

func (g *Graph) CanReach(frm Node, to Node) bool {
	stack := []Node{frm}
	var n Node
	visited := map[Node]bool{frm: true}
	for {
		if len(stack) == 0 {
			break
		}
		stack, n = stack[:len(stack)-1], stack[len(stack)-1]
		if n == to {
			return true
		}
		visited[n] = true
		for _, neighbor := range g.Neighbors(n) {
			_, ok := visited[neighbor.To]
			if ok {
				continue
			}
			stack = append(stack, neighbor.To)
		}
	}
	return false
}
