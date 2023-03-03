package naryTree

import "sync"

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
	out := make(chan Res)
	go func() {
		for request := range requests {
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
			out <- handler(request, ret)
		}
		close(out)
	}()
	return out
}

func ProcessPar[Req Request, Res any](
	nodeIds []NodeId,
	brachingFactor uint,
	handler func(Req, bool) Res,
	requests <-chan Req,
) <-chan Res {
	tree := New(nodeIds, brachingFactor)
	out := make(chan Res)
	go func() {
		wg := sync.WaitGroup{}
		for request := range requests {
			wg.Add(1)
			go func(request Req) {
				clb := func(b bool) {
					out <- handler(request, b)
					wg.Done()
				}
				switch request.Operation() {
				case Lock:
					tree.Lock(request.NodeId(), request.UserId(), clb)
				case Unlock:
					tree.Unlock(request.NodeId(), request.UserId(), clb)
				case Upgrade:
					tree.Upgrade(request.NodeId(), request.UserId(), clb)
				default:
					clb(false)
				}
			}(request)
		}
		wg.Wait()
		close(out)
	}()
	return out
}
