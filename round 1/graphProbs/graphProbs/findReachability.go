package graphProbs

func (g *Graph) ReachableNodes(frm Node, blockedNodes map[Node]bool) <-chan Node {
	retCh := make(chan Node)
	if blockedNodes == nil {
		blockedNodes = map[Node]bool{}
	}
	go func() {
		frontier := []Node{frm}
		visited := map[Node]bool{frm: true}
		for {
			if len(frontier) == 0 {
				break
			}
			nextFrontier := []Node{}
			for _, node := range frontier {
				retCh <- node
				for _, neighbor := range g.Neighbors(node) {
					_, ok := blockedNodes[neighbor.To]
					if ok {
						continue
					}
					_, ok = visited[neighbor.To]
					if ok {
						continue
					}
					nextFrontier = append(nextFrontier, neighbor.To)
					visited[neighbor.To] = true
				}
			}
			frontier = nextFrontier
		}
		close(retCh)
	}()
	return retCh
}

func (g *Graph) CanReach(frm Node, to Node) bool {
	for node := range g.ReachableNodes(frm, nil) {
		if node == to {
			return true
		}
	}
	return false
}
