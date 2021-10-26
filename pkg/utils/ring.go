/*
 * Copyright 2021 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package utils

import (
	"errors"
	"sync/atomic"
	"unsafe"

	"github.com/bytedance/gopkg/collection/lscq"
)

// ErrRingFull means the ring is full.
var ErrRingFull = errors.New("ring is full")

// Ring implements a fixed size ring buffer to manage data
type Ring struct {
	queue *lscq.PointerQueue
	quota int32
	limit int32
}

// NewRing creates a ringbuffer with fixed size.
func NewRing(size int) *Ring {
	if size <= 0 {
		// When size is an invalid number, we still return an instance
		// with zero-size to reduce error checks of the callers.
		size = 0
	}

	r := &Ring{
		queue: lscq.NewPointer(),
		quota: 0,
		limit: int32(size),
	}
	return r
}

// Push appends item to the ring.
func (r *Ring) Push(obj unsafe.Pointer) error {
	if atomic.AddInt32(&r.quota, 1) < r.limit {
		r.queue.Enqueue(unsafe.Pointer(obj))
	} else {
		atomic.AddInt32(&r.quota, -1)
	}
	return ErrRingFull
}

// Pop returns the last item and removes it from the ring.
func (r *Ring) Pop() (result interface{}) {
	var ok bool
	result, ok = r.queue.Dequeue()
	if ok && result != nil {
		atomic.AddInt32(&r.quota, -1)
	}
	return
}
