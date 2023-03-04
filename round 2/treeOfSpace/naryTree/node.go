package naryTree

import (
	"sync"
	"treeOfSpace/taggedSet"
)

type node struct {
	nodeId                           string
	lockedBy                         UserId
	isLocked                         bool
	operationLock                    sync.RWMutex
	descendantLockingUserMapping     taggedSet.TaggedSet[UserId, nodeIdx]
	descendantLockingUserMappingLock sync.RWMutex
}

func newNode(nodeId string) node {
	descendantLockingUserMapping := taggedSet.New[UserId, nodeIdx]()
	return node{
		nodeId:                           nodeId,
		isLocked:                         false,
		descendantLockingUserMapping:     descendantLockingUserMapping,
		operationLock:                    sync.RWMutex{},
		descendantLockingUserMappingLock: sync.RWMutex{},
	}
}

type operationReadLockNode struct {
	node *node
}

func (node *node) acquireOperationReadLock() operationReadLockNode {
	node.operationLock.RLock()
	return operationReadLockNode{node}
}

type operationWriteLockNode struct {
	node *node
}

func (node *node) acquireOperationWriteLock() operationWriteLockNode {
	node.operationLock.Lock()
	return operationWriteLockNode{node}
}

func (nodeH *operationWriteLockNode) lock(userId UserId) {
	nodeH.node.lockedBy = userId
	nodeH.node.isLocked = true
}

func (node *node) unlock() UserId {
	oldUser := node.lockedBy
	node.isLocked = false
	return oldUser
}

func (nodeH *operationWriteLockNode) unlock() UserId {
	return nodeH.node.unlock()
}

func (nodeH *operationWriteLockNode) releaseOperationLock() {
	nodeH.node.operationLock.Unlock()
}

func (nodeH *operationReadLockNode) releaseOperationLock() {
	nodeH.node.operationLock.RUnlock()
}

func (nodeH *operationWriteLockNode) isLocked() bool {
	return nodeH.node.isLocked
}

func (nodeH *operationReadLockNode) isLocked() bool {
	return nodeH.node.isLocked
}

func (nodeH *operationWriteLockNode) lockingUser() UserId {
	return nodeH.node.lockedBy
}

func (nodeH *operationReadLockNode) lockingUser() UserId {
	return nodeH.node.lockedBy
}

func (nodeH *operationWriteLockNode) lockedDescendants(userId UserId) []nodeIdx {
	nodeH.node.descendantLockingUserMappingLock.RLock()
	defer nodeH.node.descendantLockingUserMappingLock.RUnlock()
	return nodeH.node.descendantLockingUserMapping.Lookup(userId)
}

func (nodeH *operationReadLockNode) lockedDescendants(userId UserId) []nodeIdx {
	nodeH.node.descendantLockingUserMappingLock.RLock()
	defer nodeH.node.descendantLockingUserMappingLock.RUnlock()
	return nodeH.node.descendantLockingUserMapping.Lookup(userId)
}

func (nodeH *operationWriteLockNode) anyLockedDescendants() bool {
	nodeH.node.descendantLockingUserMappingLock.RLock()
	defer nodeH.node.descendantLockingUserMappingLock.RUnlock()
	return nodeH.node.descendantLockingUserMapping.Size() > 0
}

func (nodeH *operationReadLockNode) anyLockedDescendants() bool {
	nodeH.node.descendantLockingUserMappingLock.RLock()
	defer nodeH.node.descendantLockingUserMappingLock.RUnlock()
	return nodeH.node.descendantLockingUserMapping.Size() > 0
}

func (nodeH *operationWriteLockNode) isUpgradable(userId UserId) bool {
	nodeH.node.descendantLockingUserMappingLock.RLock()
	defer nodeH.node.descendantLockingUserMappingLock.RUnlock()
	return !nodeH.isLocked() &&
		nodeH.node.descendantLockingUserMapping.NumTags() == 1 &&
		nodeH.node.descendantLockingUserMapping.Contains(userId)
}

func (nodeH *operationReadLockNode) isUpgradable(userId UserId) bool {
	nodeH.node.descendantLockingUserMappingLock.RLock()
	defer nodeH.node.descendantLockingUserMappingLock.RUnlock()
	return !nodeH.isLocked() &&
		nodeH.node.descendantLockingUserMapping.NumTags() == 1 &&
		nodeH.node.descendantLockingUserMapping.Contains(userId)
}

func (node *node) addDescedantLockingUser(userId UserId, nodeIdx nodeIdx) {
	node.descendantLockingUserMappingLock.Lock()
	node.descendantLockingUserMapping.AddEntry(userId, nodeIdx)
	node.descendantLockingUserMappingLock.Unlock()
}

func (nodeH *operationWriteLockNode) addDescedantLockingUser(userId UserId, nodeIdx nodeIdx) {
	nodeH.node.descendantLockingUserMappingLock.Lock()
	nodeH.node.descendantLockingUserMapping.AddEntry(userId, nodeIdx)
	nodeH.node.descendantLockingUserMappingLock.Unlock()
}

func (nodeH *operationReadLockNode) addDescedantLockingUser(userId UserId, nodeIdx nodeIdx) {
	nodeH.node.descendantLockingUserMappingLock.Lock()
	nodeH.node.descendantLockingUserMapping.AddEntry(userId, nodeIdx)
	nodeH.node.descendantLockingUserMappingLock.Unlock()
}

func (node *node) removeDescedantLockingUser(userId UserId, nodeIdx nodeIdx) {
	node.descendantLockingUserMappingLock.Lock()
	node.descendantLockingUserMapping.RemoveEntry(userId, nodeIdx)
	node.descendantLockingUserMappingLock.Unlock()
}

func (nodeH *operationWriteLockNode) removeDescedantLockingUser(userId UserId, nodeIdx nodeIdx) {
	nodeH.node.descendantLockingUserMappingLock.Lock()
	nodeH.node.descendantLockingUserMapping.RemoveEntry(userId, nodeIdx)
	nodeH.node.descendantLockingUserMappingLock.Unlock()
}

func (nodeH *operationReadLockNode) removeDescedantLockingUser(userId UserId, nodeIdx nodeIdx) {
	nodeH.node.descendantLockingUserMappingLock.Lock()
	nodeH.node.descendantLockingUserMapping.RemoveEntry(userId, nodeIdx)
	nodeH.node.descendantLockingUserMappingLock.Unlock()
}
