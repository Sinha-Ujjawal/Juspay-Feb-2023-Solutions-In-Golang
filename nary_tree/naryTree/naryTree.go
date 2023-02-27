package naryTree

import nodeM "treeOfSpace/node"

type UserId = int64
type NodeId = string
type NodeIdx = uint

type NaryTree struct {
	nodes              []*nodeM.Node[UserId, NodeIdx]
	brachingFactor     uint
	nodeToIndexMapping map[NodeId]NodeIdx
}

func New(nodeIds []NodeId, brachingFactor uint) NaryTree {
	nodeToIndexMapping := make(map[string]uint, len(nodeIds))
	var nodes []*nodeM.Node[UserId, NodeIdx]
	for idx, nodeId := range nodeIds {
		nodeToIndexMapping[nodeId] = uint(idx)
		node := nodeM.New[UserId, NodeIdx](nodeId)
		nodes = append(nodes, &node)
	}
	return NaryTree{
		nodes,
		brachingFactor,
		nodeToIndexMapping,
	}
}

func (tree *NaryTree) parentOf(idx NodeIdx) NodeIdx {
	return (idx - 1) / tree.brachingFactor
}

func (tree *NaryTree) firstChildOf(idx NodeIdx) NodeIdx {
	return idx*tree.brachingFactor + 1
}

func (tree *NaryTree) ancestorsOf(nodeIdx NodeIdx) []NodeIdx {
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

func (tree *NaryTree) getNodeIdx(nodeId NodeId) (*NodeIdx, bool) {
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
	nodeH := node.AcquireOperationWriteLock()
	defer nodeH.ReleaseOperationLock()
	var ancestorHs []nodeM.OperationReadLockNode[UserId, NodeIdx]
	for _, ancestorIdx := range tree.ancestorsOf(*nodeIdx) {
		ancestor := tree.nodes[ancestorIdx]
		ancestorH := ancestor.AcquireOperationReadLock()
		ancestorHs = append(ancestorHs, ancestorH)
		if ancestorH.IsLocked() {
			clb(false)
			for _, ancestorH = range ancestorHs {
				ancestorH.ReleaseOperationLock()
			}
			return
		}
	}
	if nodeH.IsLocked() || nodeH.AnyLockedDescendants() {
		clb(false)
		for _, ancestorH := range ancestorHs {
			ancestorH.ReleaseOperationLock()
		}
		return
	}

	clb(true)
	nodeH.Lock(userId)
	for _, ancestorH := range ancestorHs {
		ancestorH.ReleaseOperationLock()
		ancestorH.AddDescedantLockingUser(userId, *nodeIdx)
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
	nodeH := node.AcquireOperationWriteLock()
	defer nodeH.ReleaseOperationLock()
	var ancestorHs []nodeM.OperationReadLockNode[UserId, NodeIdx]
	for _, ancestorIdx := range tree.ancestorsOf(*nodeIdx) {
		ancestor := tree.nodes[ancestorIdx]
		ancestorH := ancestor.AcquireOperationReadLock()
		ancestorHs = append(ancestorHs, ancestorH)
	}
	if !nodeH.IsLocked() {
		clb(false)
		for _, ancestorH := range ancestorHs {
			ancestorH.ReleaseOperationLock()
		}
		return
	}
	lockingUser := nodeH.LockingUser()
	if lockingUser == nil || *lockingUser != userId {
		clb(false)
		for _, ancestorH := range ancestorHs {
			ancestorH.ReleaseOperationLock()
		}
		return
	}

	clb(true)
	nodeH.Unlock()
	for _, ancestorH := range ancestorHs {
		ancestorH.ReleaseOperationLock()
		ancestorH.RemoveDescedantLockingUser(userId, *nodeIdx)
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
	nodeH := node.AcquireOperationWriteLock()
	defer nodeH.ReleaseOperationLock()
	var ancestorHs []nodeM.OperationReadLockNode[UserId, NodeIdx]
	for _, ancestorIdx := range tree.ancestorsOf(*nodeIdx) {
		ancestor := tree.nodes[ancestorIdx]
		ancestorH := ancestor.AcquireOperationReadLock()
		ancestorHs = append(ancestorHs, ancestorH)
	}
	if !nodeH.IsUpgradable(userId) {
		clb(false)
		for _, ancestorH := range ancestorHs {
			ancestorH.ReleaseOperationLock()
		}
		return
	}
	clb(true)
	nodeH.Lock(userId)
	for _, ancestorH := range ancestorHs {
		ancestorH.AddDescedantLockingUser(userId, *nodeIdx)
	}
	for _, lockedDescendantIdx := range nodeH.LockedDescendants(userId) {
		lockedDescendant := tree.nodes[lockedDescendantIdx]
		for _, ancestorIdx := range tree.ancestorsOf(lockedDescendantIdx) {
			ancestor := tree.nodes[ancestorIdx]
			ancestor.RemoveDescedantLockingUser(userId, lockedDescendantIdx)
		}
		lockedDescendant.Unlock()
	}
	for _, ancestorH := range ancestorHs {
		ancestorH.ReleaseOperationLock()
	}
	return
}
