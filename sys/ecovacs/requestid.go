package ecovacs

import (
	"fmt"
	"sync"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type RequestId struct {
	reqId uint
	sync.Mutex
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Returns next requestId as a string
func (this *RequestId) Next() string {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	this.reqId = this.reqId + 1
	return fmt.Sprint(this.reqId)
}
