package naryTree

import (
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
		var ret = false
		var retPtr = &ret
		switch request.Operation() {
		case Lock:
			tree.Lock(request.NodeId(), request.UserId(), func(b bool) { *retPtr = b })
		case Unlock:
			tree.Unlock(request.NodeId(), request.UserId(), func(b bool) { *retPtr = b })
		case Upgrade:
			tree.Upgrade(request.NodeId(), request.UserId(), func(b bool) { *retPtr = b })
		default:
			break
		}
		return handler(request, ret)
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
	processor := func(request Req, clb func(Res)) {
		clb2 := func(b bool) {
			clb(handler(request, b))
		}
		switch request.Operation() {
		case Lock:
			tree.Lock(request.NodeId(), request.UserId(), clb2)
		case Unlock:
			tree.Unlock(request.NodeId(), request.UserId(), clb2)
		case Upgrade:
			tree.Upgrade(request.NodeId(), request.UserId(), clb2)
		default:
			clb(handler(request, false))
		}
	}
	return chanUtils.MapPar(requests, processor)
}
