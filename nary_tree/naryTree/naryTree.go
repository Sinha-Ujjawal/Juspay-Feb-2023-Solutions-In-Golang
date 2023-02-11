package naryTree

import "treeOfSpace/node"

type UserId = int64
type NodeId = string
type nodeIdx = uint

type NaryTree struct {
	nodes              []*node.Node[UserId, nodeIdx]
	brachingFactor     uint
	nodeToIndexMapping map[NodeId]nodeIdx
}

func New(nodeIds []NodeId, brachingFactor uint) NaryTree {
	nodeToIndexMapping := make(map[string]uint, len(nodeIds))
	var nodes []*node.Node[UserId, nodeIdx]
	for idx, nodeId := range nodeIds {
		nodeToIndexMapping[nodeId] = uint(idx)
		node := node.New[UserId, nodeIdx](nodeId)
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

func (tree *NaryTree) lock(node *node.Node[UserId, nodeIdx], nodeIdx *nodeIdx, userId UserId) {
	for _, ancestorIdx := range tree.ancestorsOf(*nodeIdx) {
		ancestor := tree.nodes[ancestorIdx]
		ancestor.AddDescedantLockingUser(userId, *nodeIdx)
	}
	node.Lock(userId)
}

func (tree *NaryTree) Lock(nodeId NodeId, userId UserId) bool {
	nodeIdx, ok := tree.getNodeIdx(nodeId)
	if !ok {
		return false
	}
	node := tree.nodes[*nodeIdx]
	if node.IsLocked() {
		return false
	}
	for _, ancestorIdx := range tree.ancestorsOf(*nodeIdx) {
		ancestor := tree.nodes[ancestorIdx]
		if ancestor.IsLocked() {
			return false
		}
	}
	if node.AnyLockedDescendants() {
		return false
	}
	tree.lock(node, nodeIdx, userId)
	return true
}

func (tree *NaryTree) unlock(node *node.Node[UserId, nodeIdx], nodeIdx *nodeIdx, userId UserId) {
	for _, ancestorIdx := range tree.ancestorsOf(*nodeIdx) {
		ancestor := tree.nodes[ancestorIdx]
		ancestor.RemoveDescedantLockingUser(userId, *nodeIdx)
	}
	node.Unlock()
}

func (tree *NaryTree) Unlock(nodeId NodeId, userId UserId) bool {
	nodeIdx, ok := tree.getNodeIdx(nodeId)
	if !ok {
		return false
	}
	node := tree.nodes[*nodeIdx]
	if !node.IsLocked() {
		return false
	}
	lockingUser := node.LockingUser()
	if lockingUser == nil || *lockingUser != userId {
		return false
	}
	tree.unlock(node, nodeIdx, userId)
	return true
}

func (tree *NaryTree) Upgrade(nodeId NodeId, userId UserId) bool {
	nodeIdx, ok := tree.getNodeIdx(nodeId)
	if !ok {
		return false
	}
	node := tree.nodes[*nodeIdx]
	if !node.IsUpgradable(userId) {
		return false
	}
	tree.lock(node, nodeIdx, userId)
	for _, lockedDescendantIdx := range node.LockedDescendants(userId) {
		lockedDescendant := tree.nodes[lockedDescendantIdx]
		tree.unlock(lockedDescendant, &lockedDescendantIdx, userId)
	}
	return true
}
