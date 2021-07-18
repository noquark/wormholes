package pipe

type Task chan interface{}
type Queue chan chan interface{}
