package naryTree

type UserId = int64
type NodeId = string
type nodeIdx = uint

type NaryTree struct {
	nodes              []*node
	brachingFactor     uint
	nodeToIndexMapping map[NodeId]nodeIdx
}

func New(nodeIds []NodeId, brachingFactor uint) NaryTree {
	nodeToIndexMapping := make(map[string]uint, len(nodeIds))
	var nodes []*node
	for idx, nodeId := range nodeIds {
		nodeToIndexMapping[nodeId] = uint(idx)
		node := newNode(nodeId)
		nodes = append(nodes, &node)
	}
	return NaryTree{
		nodes,
		brachingFactor,
		nodeToIndexMapping,
	}
}

func (tree *NaryTree) parentOf(idx nodeIdx) nodeIdx {
	return (idx - 1) / tree.brachingFactor
}

func (tree *NaryTree) firstChildOf(idx nodeIdx) nodeIdx {
	return idx*tree.brachingFactor + 1
}

func (tree *NaryTree) ancestorsOf(nodeIdx nodeIdx) []nodeIdx {
	if nodeIdx <= 0 {
		return nil
	}
	idx := tree.parentOf(nodeIdx)
	var ancestors []uint
	for {
		ancestors = append(ancestors, idx)
		if idx == 0 {
			break
		}
		idx = tree.parentOf(idx)
	}
	return ancestors
}

func (tree *NaryTree) getNodeIdx(nodeId NodeId) (*nodeIdx, bool) {
	nodeIdx, ok := tree.nodeToIndexMapping[nodeId]
	if !ok {
		return nil, false
	}
	return &nodeIdx, true
}

func (tree *NaryTree) Lock(
	nodeId NodeId,
	userId UserId,
	clb func(bool),
) {
	nodeIdx, ok := tree.getNodeIdx(nodeId)
	if !ok {
		clb(false)
		return
	}
	node := tree.nodes[*nodeIdx]
	nodeH := node.acquireWriteLock()
	defer nodeH.releaseLock()
	var ancestorHs []readLock
	for _, ancestorIdx := range tree.ancestorsOf(*nodeIdx) {
		ancestor := tree.nodes[ancestorIdx]
		ancestorH := ancestor.acquireReadLock()
		ancestorHs = append(ancestorHs, ancestorH)
		if ancestorH.isLocked() {
			clb(false)
			for _, ancestorH = range ancestorHs {
				ancestorH.releaseLock()
			}
			return
		}
	}
	if nodeH.isLocked() || nodeH.anyLockedDescendants() {
		clb(false)
		for _, ancestorH := range ancestorHs {
			ancestorH.releaseLock()
		}
		return
	}

	clb(true)
	nodeH.lock(userId)
	for _, ancestorH := range ancestorHs {
		ancestorH.releaseLock()
		ancestorH.addDescedantLockingUser(userId, *nodeIdx)
	}
	return
}

func (tree *NaryTree) Unlock(
	nodeId NodeId,
	userId UserId,
	clb func(bool),
) {
	nodeIdx, ok := tree.getNodeIdx(nodeId)
	if !ok {
		clb(false)
		return
	}
	node := tree.nodes[*nodeIdx]
	nodeH := node.acquireWriteLock()
	defer nodeH.releaseLock()
	var ancestorHs []readLock
	for _, ancestorIdx := range tree.ancestorsOf(*nodeIdx) {
		ancestor := tree.nodes[ancestorIdx]
		ancestorH := ancestor.acquireReadLock()
		ancestorHs = append(ancestorHs, ancestorH)
	}
	if !nodeH.isLocked() {
		clb(false)
		for _, ancestorH := range ancestorHs {
			ancestorH.releaseLock()
		}
		return
	}
	lockingUser := nodeH.lockingUser()
	if lockingUser != userId {
		clb(false)
		for _, ancestorH := range ancestorHs {
			ancestorH.releaseLock()
		}
		return
	}

	clb(true)
	nodeH.unlock()
	for _, ancestorH := range ancestorHs {
		ancestorH.releaseLock()
		ancestorH.removeDescedantLockingUser(userId, *nodeIdx)
	}
	return
}

func (tree *NaryTree) Upgrade(
	nodeId NodeId,
	userId UserId,
	clb func(bool),
) {
	nodeIdx, ok := tree.getNodeIdx(nodeId)
	if !ok {
		clb(false)
		return
	}
	node := tree.nodes[*nodeIdx]
	nodeH := node.acquireWriteLock()
	defer nodeH.releaseLock()
	var ancestorHs []readLock
	for _, ancestorIdx := range tree.ancestorsOf(*nodeIdx) {
		ancestor := tree.nodes[ancestorIdx]
		ancestorH := ancestor.acquireReadLock()
		ancestorHs = append(ancestorHs, ancestorH)
	}
	if !nodeH.isUpgradable(userId) {
		clb(false)
		for _, ancestorH := range ancestorHs {
			ancestorH.releaseLock()
		}
		return
	}
	clb(true)
	nodeH.lock(userId)
	for _, ancestorH := range ancestorHs {
		ancestorH.addDescedantLockingUser(userId, *nodeIdx)
	}
	for _, lockedDescendantIdx := range nodeH.lockedDescendants(userId) {
		lockedDescendant := tree.nodes[lockedDescendantIdx]
		for _, ancestorIdx := range tree.ancestorsOf(lockedDescendantIdx) {
			ancestor := tree.nodes[ancestorIdx]
			ancestor.removeDescedantLockingUser(userId, lockedDescendantIdx)
		}
		lockedDescendant.unlock()
	}
	for _, ancestorH := range ancestorHs {
		ancestorH.releaseLock()
	}
	return
}
