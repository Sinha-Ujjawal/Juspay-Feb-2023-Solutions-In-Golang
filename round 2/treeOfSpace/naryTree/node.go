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

type readLock struct {
	node *node
}

func (node *node) acquireReadLock() readLock {
	node.operationLock.RLock()
	return readLock{node}
}

type writeLock struct {
	node *node
}

func (node *node) acquireWriteLock() writeLock {
	node.operationLock.Lock()
	return writeLock{node}
}

func (nodeH *writeLock) lock(userId UserId) {
	nodeH.node.lockedBy = userId
	nodeH.node.isLocked = true
}

func (node *node) unlock() UserId {
	oldUser := node.lockedBy
	node.isLocked = false
	return oldUser
}

func (nodeH *writeLock) unlock() UserId {
	return nodeH.node.unlock()
}

func (nodeH *writeLock) releaseLock() {
	nodeH.node.operationLock.Unlock()
}

func (nodeH *readLock) releaseLock() {
	nodeH.node.operationLock.RUnlock()
}

func (nodeH *writeLock) isLocked() bool {
	return nodeH.node.isLocked
}

func (nodeH *readLock) isLocked() bool {
	return nodeH.node.isLocked
}

func (nodeH *writeLock) lockingUser() UserId {
	return nodeH.node.lockedBy
}

func (nodeH *readLock) lockingUser() UserId {
	return nodeH.node.lockedBy
}

func (nodeH *writeLock) lockedDescendants(userId UserId) []nodeIdx {
	nodeH.node.descendantLockingUserMappingLock.RLock()
	defer nodeH.node.descendantLockingUserMappingLock.RUnlock()
	return nodeH.node.descendantLockingUserMapping.Lookup(userId)
}

func (nodeH *readLock) lockedDescendants(userId UserId) []nodeIdx {
	nodeH.node.descendantLockingUserMappingLock.RLock()
	defer nodeH.node.descendantLockingUserMappingLock.RUnlock()
	return nodeH.node.descendantLockingUserMapping.Lookup(userId)
}

func (nodeH *writeLock) anyLockedDescendants() bool {
	nodeH.node.descendantLockingUserMappingLock.RLock()
	defer nodeH.node.descendantLockingUserMappingLock.RUnlock()
	return nodeH.node.descendantLockingUserMapping.Size() > 0
}

func (nodeH *readLock) anyLockedDescendants() bool {
	nodeH.node.descendantLockingUserMappingLock.RLock()
	defer nodeH.node.descendantLockingUserMappingLock.RUnlock()
	return nodeH.node.descendantLockingUserMapping.Size() > 0
}

func (nodeH *writeLock) isUpgradable(userId UserId) bool {
	nodeH.node.descendantLockingUserMappingLock.RLock()
	defer nodeH.node.descendantLockingUserMappingLock.RUnlock()
	return !nodeH.isLocked() &&
		nodeH.node.descendantLockingUserMapping.NumTags() == 1 &&
		nodeH.node.descendantLockingUserMapping.Contains(userId)
}

func (nodeH *readLock) isUpgradable(userId UserId) bool {
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

func (nodeH *writeLock) addDescedantLockingUser(userId UserId, nodeIdx nodeIdx) {
	nodeH.node.descendantLockingUserMappingLock.Lock()
	nodeH.node.descendantLockingUserMapping.AddEntry(userId, nodeIdx)
	nodeH.node.descendantLockingUserMappingLock.Unlock()
}

func (nodeH *readLock) addDescedantLockingUser(userId UserId, nodeIdx nodeIdx) {
	nodeH.node.descendantLockingUserMappingLock.Lock()
	nodeH.node.descendantLockingUserMapping.AddEntry(userId, nodeIdx)
	nodeH.node.descendantLockingUserMappingLock.Unlock()
}

func (node *node) removeDescedantLockingUser(userId UserId, nodeIdx nodeIdx) {
	node.descendantLockingUserMappingLock.Lock()
	node.descendantLockingUserMapping.RemoveEntry(userId, nodeIdx)
	node.descendantLockingUserMappingLock.Unlock()
}

func (nodeH *writeLock) removeDescedantLockingUser(userId UserId, nodeIdx nodeIdx) {
	nodeH.node.descendantLockingUserMappingLock.Lock()
	nodeH.node.descendantLockingUserMapping.RemoveEntry(userId, nodeIdx)
	nodeH.node.descendantLockingUserMappingLock.Unlock()
}

func (nodeH *readLock) removeDescedantLockingUser(userId UserId, nodeIdx nodeIdx) {
	nodeH.node.descendantLockingUserMappingLock.Lock()
	nodeH.node.descendantLockingUserMapping.RemoveEntry(userId, nodeIdx)
	nodeH.node.descendantLockingUserMappingLock.Unlock()
}
