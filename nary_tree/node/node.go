package node

import (
	"sync"
	"treeOfSpace/taggedSet"
)

type Node[U comparable, N comparable] struct {
	nodeId                           string
	lockedBy                         *U
	operationLock                    sync.RWMutex
	descendantLockingUserMapping     *taggedSet.TaggedSet[U, N]
	descendantLockingUserMappingLock sync.RWMutex
}

func New[U comparable, N comparable](nodeId string) Node[U, N] {
	descendantLockingUserMapping := taggedSet.New[U, N]()
	return Node[U, N]{
		nodeId:                           nodeId,
		descendantLockingUserMapping:     &descendantLockingUserMapping,
		operationLock:                    sync.RWMutex{},
		descendantLockingUserMappingLock: sync.RWMutex{},
	}
}

type OperationReadLockNode[U comparable, N comparable] struct {
	node *Node[U, N]
}

func (node *Node[U, N]) AcquireOperationReadLock() OperationReadLockNode[U, N] {
	node.operationLock.RLock()
	return OperationReadLockNode[U, N]{node}
}

type OperationWriteLockNode[U comparable, N comparable] struct {
	node *Node[U, N]
}

func (node *Node[U, N]) AcquireOperationWriteLock() OperationWriteLockNode[U, N] {
	node.operationLock.Lock()
	return OperationWriteLockNode[U, N]{node}
}

func (nodeH *OperationWriteLockNode[U, N]) Lock(userId U) {
	nodeH.node.lockedBy = &userId
}

func (nodeH *OperationWriteLockNode[U, N]) Unlock() *U {
	oldUser := nodeH.node.lockedBy
	nodeH.node.lockedBy = nil
	return oldUser
}

func (node *Node[U, N]) Unlock() *U {
	oldUser := node.lockedBy
	node.lockedBy = nil
	return oldUser
}

func (nodeH *OperationWriteLockNode[U, N]) ReleaseOperationLock() {
	nodeH.node.operationLock.Unlock()
}

func (nodeH *OperationReadLockNode[U, N]) ReleaseOperationLock() {
	nodeH.node.operationLock.RUnlock()
}

func (nodeH *OperationWriteLockNode[U, N]) IsLocked() bool {
	return nodeH.node.lockedBy != nil
}

func (nodeH *OperationReadLockNode[U, N]) IsLocked() bool {
	return nodeH.node.lockedBy != nil
}

func (nodeH *OperationWriteLockNode[U, N]) LockingUser() *U {
	return nodeH.node.lockedBy
}

func (nodeH *OperationReadLockNode[U, N]) LockingUser() *U {
	return nodeH.node.lockedBy
}

func (nodeH *OperationWriteLockNode[U, N]) LockedDescendants(userId U) []N {
	nodeH.node.descendantLockingUserMappingLock.RLock()
	defer nodeH.node.descendantLockingUserMappingLock.RUnlock()
	return nodeH.node.descendantLockingUserMapping.Lookup(userId)
}

func (nodeH *OperationReadLockNode[U, N]) LockedDescendants(userId U) []N {
	nodeH.node.descendantLockingUserMappingLock.RLock()
	defer nodeH.node.descendantLockingUserMappingLock.RUnlock()
	return nodeH.node.descendantLockingUserMapping.Lookup(userId)
}

func (nodeH *OperationWriteLockNode[U, N]) AnyLockedDescendants() bool {
	nodeH.node.descendantLockingUserMappingLock.RLock()
	defer nodeH.node.descendantLockingUserMappingLock.RUnlock()
	return nodeH.node.descendantLockingUserMapping.Size() > 0
}

func (nodeH *OperationReadLockNode[U, N]) AnyLockedDescendants() bool {
	nodeH.node.descendantLockingUserMappingLock.RLock()
	defer nodeH.node.descendantLockingUserMappingLock.RUnlock()
	return nodeH.node.descendantLockingUserMapping.Size() > 0
}

func (nodeH *OperationWriteLockNode[U, N]) IsUpgradable(userId U) bool {
	nodeH.node.descendantLockingUserMappingLock.RLock()
	defer nodeH.node.descendantLockingUserMappingLock.RUnlock()
	return !nodeH.IsLocked() &&
		nodeH.node.descendantLockingUserMapping.NumTags() == 1 &&
		nodeH.node.descendantLockingUserMapping.Contains(userId)
}

func (nodeH *OperationReadLockNode[U, N]) IsUpgradable(userId U) bool {
	nodeH.node.descendantLockingUserMappingLock.RLock()
	defer nodeH.node.descendantLockingUserMappingLock.RUnlock()
	return !nodeH.IsLocked() &&
		nodeH.node.descendantLockingUserMapping.NumTags() == 1 &&
		nodeH.node.descendantLockingUserMapping.Contains(userId)
}

func (node *Node[U, N]) AddDescedantLockingUser(userId U, nodeIdx N) {
	node.descendantLockingUserMappingLock.Lock()
	node.descendantLockingUserMapping.AddEntry(userId, nodeIdx)
	node.descendantLockingUserMappingLock.Unlock()
}

func (nodeH *OperationWriteLockNode[U, N]) AddDescedantLockingUser(userId U, nodeIdx N) {
	nodeH.node.descendantLockingUserMappingLock.Lock()
	nodeH.node.descendantLockingUserMapping.AddEntry(userId, nodeIdx)
	nodeH.node.descendantLockingUserMappingLock.Unlock()
}

func (nodeH *OperationReadLockNode[U, N]) AddDescedantLockingUser(userId U, nodeIdx N) {
	nodeH.node.descendantLockingUserMappingLock.Lock()
	nodeH.node.descendantLockingUserMapping.AddEntry(userId, nodeIdx)
	nodeH.node.descendantLockingUserMappingLock.Unlock()
}

func (node *Node[U, N]) RemoveDescedantLockingUser(userId U, nodeIdx N) {
	node.descendantLockingUserMappingLock.Lock()
	node.descendantLockingUserMapping.RemoveEntry(userId, nodeIdx)
	node.descendantLockingUserMappingLock.Unlock()
}

func (nodeH *OperationWriteLockNode[U, N]) RemoveDescedantLockingUser(userId U, nodeIdx N) {
	nodeH.node.descendantLockingUserMappingLock.Lock()
	nodeH.node.descendantLockingUserMapping.RemoveEntry(userId, nodeIdx)
	nodeH.node.descendantLockingUserMappingLock.Unlock()
}

func (nodeH *OperationReadLockNode[U, N]) RemoveDescedantLockingUser(userId U, nodeIdx N) {
	nodeH.node.descendantLockingUserMappingLock.Lock()
	nodeH.node.descendantLockingUserMapping.RemoveEntry(userId, nodeIdx)
	nodeH.node.descendantLockingUserMappingLock.Unlock()
}
