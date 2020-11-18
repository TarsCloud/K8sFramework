package main

type ReconcileResult uint

const (
	AllOk      ReconcileResult = 0
	RateLimit  ReconcileResult = 1
	FatalError ReconcileResult = 2
)

type TReconcile interface {
	EnqueueObj(interface{})
	Start(chan struct{})
}
