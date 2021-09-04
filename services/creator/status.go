package main

import "sync/atomic"

type Status int32

func NewStatus() *Status {
	return new(Status)
}

func (r *Status) SetBusy() {
	atomic.StoreInt32((*int32)(r), 1)
}

func (r *Status) SetIdle() {
	atomic.StoreInt32((*int32)(r), 0)
}

func (r *Status) IsBusy() bool {
	return atomic.LoadInt32((*int32)(r))&1 == 1
}

func (r *Status) IsIdle() bool {
	return !r.IsBusy()
}
