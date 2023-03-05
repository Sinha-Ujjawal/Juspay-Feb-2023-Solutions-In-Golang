package graphProbs

func (g *Graph) NeighborsToBlockToEnsureUnreachability(follower Node, following Node) <-chan Node {
	retCh := make(chan Node)
	go func() {
		immediateParents := g.ImmediateParents(following)
		for reachableNode := range g.ReachableNodes(follower, map[Node]bool{following: true}) {
			_, ok := immediateParents[reachableNode]
			if ok {
				retCh <- reachableNode
			}
		}
		close(retCh)
	}()
	return retCh
}
