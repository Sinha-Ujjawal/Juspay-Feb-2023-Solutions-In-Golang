package naryTree

import "sync"

type Operation = string

const (
	Lock    Operation = "Lock"
	Unlock            = "Unlock"
	Upgrade           = "Upgrade"
)

type Request struct {
	Operation Operation
	UserId    UserId
	NodeId    NodeId
}

type Response struct {
	Request Request
	Result  bool
}

func ProcessSeq(
	nodeIds []NodeId,
	brachingFactor uint,
	requests <-chan Request,
) <-chan Response {
	tree := New(nodeIds, brachingFactor)
	out := make(chan Response)
	go func() {
		for request := range requests {
			clb := func(result bool) {
				out <- Response{request, result}
			}
			switch request.Operation {
			case Lock:
				tree.Lock(request.NodeId, request.UserId, clb)
			case Unlock:
				tree.Unlock(request.NodeId, request.UserId, clb)
			case Upgrade:
				tree.Upgrade(request.NodeId, request.UserId, clb)
			default:
				clb(false)
			}
		}
		close(out)
	}()
	return out
}

func ProcessPar(
	nodeIds []NodeId,
	brachingFactor uint,
	requests <-chan Request,
) <-chan Response {
	tree := New(nodeIds, brachingFactor)
	out := make(chan Response)
	go func() {
		wg := sync.WaitGroup{}
		for request := range requests {
			wg.Add(1)
			go func(request Request) {
				clb := func(result bool) {
					out <- Response{request, result}
				}
				switch request.Operation {
				case Lock:
					tree.Lock(request.NodeId, request.UserId, clb)
				case Unlock:
					tree.Unlock(request.NodeId, request.UserId, clb)
				case Upgrade:
					tree.Upgrade(request.NodeId, request.UserId, clb)
				default:
					clb(false)
				}
				wg.Done()
			}(request)
		}
		wg.Wait()
		close(out)
	}()
	return out
}
