package compatible


type NodeDetailWrapper struct {
	Detail []NodeRecordDetail
	By func(e1, e2 *NodeRecordDetail) bool
}

func (e NodeDetailWrapper) Len() int {
	return len(e.Detail)
}

func (e NodeDetailWrapper) Swap(i, j int)  {
	e.Detail[i], e.Detail[j] = e.Detail[j], e.Detail[i]
}

func (e NodeDetailWrapper) Less(i, j int) bool  {
	return e.By(&e.Detail[i], &e.Detail[j])
}
