package twitchirc

import (
	"sync"
)

// MsgObj -
type MsgObj struct {
	Command string
	Params  []string
}

// QueueList -
type QueueList struct {
	Messages []*MsgObj
	Lock     *sync.Mutex
}

// NewQueue -
func NewQueue() *QueueList {
	return &QueueList{}
}

// Add - add element
func (q *QueueList) Add(m *MsgObj) {
	q.Messages = append(q.Messages, m)
}

// Get - get element
func (q *QueueList) Get() (m *MsgObj) {
	if q.IsEmpty() {
		return nil
	}
	m = q.Messages[0]
	q.Messages = q.Messages[1:]
	return
}

// IsEmpty -
func (q *QueueList) IsEmpty() bool {
	if len(q.Messages) == 0 {
		return true
	}
	return false
}

// Size -
func (q *QueueList) Size() int {
	return len(q.Messages)
}

// Clear -
func (q *QueueList) Clear() {
	if q.IsEmpty() {
		return
	}

	for i := 0; i < len(q.Messages); i++ {
		q.Messages[i] = nil
	}
	q.Messages = nil
	return
}
