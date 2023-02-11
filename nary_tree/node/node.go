package node

import "treeOfSpace/taggedSet"

type Node[U comparable, N comparable] struct {
	nodeId                       string
	lockedBy                     *U
	descendantLockingUserMapping *taggedSet.TaggedSet[U, N]
}

func New[U comparable, N comparable](nodeId string) Node[U, N] {
	descendantLockingUserMapping := taggedSet.New[U, N]()
	return Node[U, N]{
		nodeId:                       nodeId,
		descendantLockingUserMapping: &descendantLockingUserMapping,
	}
}

func (node *Node[U, N]) IsLocked() bool {
	return node.lockedBy != nil
}

func (node *Node[U, N]) LockingUser() *U {
	return node.lockedBy
}

func (node *Node[U, N]) Lock(userId U) {
	node.lockedBy = &userId
}

func (node *Node[U, N]) Unlock() *U {
	oldUser := node.lockedBy
	node.lockedBy = nil
	return oldUser
}

func (node *Node[U, N]) LockedDescendants(userId U) []N {
	return node.descendantLockingUserMapping.Lookup(userId)
}

func (node *Node[U, N]) AnyLockedDescendants() bool {
	return node.descendantLockingUserMapping.Size() > 0
}

func (node *Node[U, N]) IsUpgradable(userId U) bool {
	return !node.IsLocked() &&
		node.descendantLockingUserMapping.NumTags() == 1 &&
		node.descendantLockingUserMapping.Contains(userId)
}

func (node *Node[U, N]) AddDescedantLockingUser(userId U, nodeIdx N) {
	node.descendantLockingUserMapping.AddEntry(userId, nodeIdx)
}

func (node *Node[U, N]) RemoveDescedantLockingUser(userId U, nodeIdx N) {
	node.descendantLockingUserMapping.RemoveEntry(userId, nodeIdx)
}
