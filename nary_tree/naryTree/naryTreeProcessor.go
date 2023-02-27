package naryTree

import (
	"sync"
	"treeOfSpace/chanUtils"
)

type Operation = uint

const (
	Lock    Operation = 1
	Unlock            = 2
	Upgrade           = 3
)

type Request interface {
	Operation() Operation
	UserId() UserId
	NodeId() NodeId
}

func ProcessSeq[Req Request, Res any](
	nodeIds []NodeId,
	brachingFactor uint,
	handler func(Req, bool) Res,
	requests <-chan Req,
) <-chan Res {
	tree := New(nodeIds, brachingFactor)
	processor := func(request Req) Res {
		switch request.Operation() {
		case Lock:
			return handler(request, tree.Lock(request.NodeId(), request.UserId()))
		case Unlock:
			return handler(request, tree.Unlock(request.NodeId(), request.UserId()))
		case Upgrade:
			return handler(request, tree.Upgrade(request.NodeId(), request.UserId()))
		default:
			return handler(request, false)
		}
	}
	return chanUtils.Map(requests, processor)
}

func ProcessPar[Req Request, Res any](
	nodeIds []NodeId,
	brachingFactor uint,
	handler func(Req, bool) Res,
	requests <-chan Req,
) <-chan Res {
	tree := New(nodeIds, brachingFactor)
	lock := sync.Mutex{}
	processor := func(request Req, clb func(Res)) {
		switch request.Operation() {
		case Lock:
			lock.Lock()
			defer lock.Unlock()
			clb(handler(request, tree.Lock(request.NodeId(), request.UserId())))
		case Unlock:
			lock.Lock()
			defer lock.Unlock()
			clb(handler(request, tree.Unlock(request.NodeId(), request.UserId())))
		case Upgrade:
			lock.Lock()
			defer lock.Unlock()
			clb(handler(request, tree.Upgrade(request.NodeId(), request.UserId())))
		default:
			clb(handler(request, false))
		}
	}
	return chanUtils.MapPar(requests, processor)
}
