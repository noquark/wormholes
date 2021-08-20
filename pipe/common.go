package pipe

type (
	Task  chan interface{}
	Queue chan chan interface{}
)
