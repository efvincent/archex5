package MemoryEventStore

import (
	"errors"
	"fmt"
	"sync"

	"github.com/efvincent/archex5/eventStore"
	es "github.com/efvincent/archex5/eventStore"
	"github.com/efvincent/archex5/eventStore/esErrors.go"
)

type MemoryEventStore struct {
	nss   map[string]map[string][]es.EventEnvelope
	mutex *sync.Mutex
}

// Create an initialized memory event store
// Using a top level mutex to sync access to the map that is the memory event store.
// a better option would be use more granular mutexes to synchronize the top level
// map of streams and then a mutex to sync each stream. Left as an exercise.
func makeMemoryEventStore() es.EventStore {
	return MemoryEventStore{
		nss:   map[string]map[string][]es.EventEnvelope{},
		mutex: &sync.Mutex{},
	}
}

// Make a singleton event store available for all processes to use
var SingletonMemoryEventStore = makeMemoryEventStore()

// Gets namespaces
func (ms MemoryEventStore) GetNamespaces() ([]string, error) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	keys := make([]string, len(ms.nss))
	i := 0
	for k := range ms.nss {
		keys[i] = k
		i = i + 1
	}
	return keys, nil
}

func (ms MemoryEventStore) GetStreams(ns string) ([]string, error) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	if n, ok := ms.nss[ns]; ok {
		streams := make([]string, len(n))
		i := 0
		for k := range n {
			streams[i] = k
			i = i + 1
		}
		return streams, nil
	}
	return []string{}, nil
}

// Checks if a namespace exists
func (ms MemoryEventStore) NamespaceExists(ns string) (bool, error) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	_, ok := ms.nss[ns]
	return ok, nil
}

// Check if a Stream exists in a namespace
func (ms MemoryEventStore) StreamExists(ns string, streamId string) (bool, error) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	if nspace, ok := ms.nss[ns]; ok {
		_, ok := nspace[streamId]
		return ok, nil
	}
	return false, nil
}

// Writes a single event into a stream in a namespace using the consistency mode. If the write
// fails a custom error is returned
func (ms MemoryEventStore) WriteEvent(ns string, streamId string,
	cMode es.ConcurrencyMode, expected int64, e *es.EventEnvelope) (int64, error) {
	events := []es.EventEnvelope{*e}
	lastId, err := ms.WriteBatch(ns, streamId, cMode, expected, events)
	return lastId, err
}

// Write several events into a stream as a single operation
func (ms MemoryEventStore) WriteBatch(ns string, streamId string,
	cMode es.ConcurrencyMode, expected int64, events []es.EventEnvelope) (int64, error) {

	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	// The memory based event store creates the namespace when it
	// hasn't been seen before. This might not be possible in all
	// implementations, and might make sense to have an explicit
	// create namespace api on the event store interface
	nspace, ok := ms.nss[ns]
	if !ok {
		nspace = map[string][]es.EventEnvelope{}
		ms.nss[ns] = nspace
	}

	// I'm assuming that the pass by value here does a deep copy of the event so
	// that the event store is storing a copy that cannot be mutated later
	insertEvent := func(e es.EventEnvelope) (int64, error) {
		if strm, ok := nspace[streamId]; ok {
			switch cMode {
			case eventStore.ANY:
			case eventStore.EXISTING_STREAM:
				e.SeqNum = int64(len(strm))
				nspace[streamId] = append(strm, e)
				return e.SeqNum, nil
			case eventStore.EXPECTING_SEQ_NUM:
				if len(strm) == 0 {
					return 0, esErrors.NewSeqExpectedErr(streamId, expected, -1)
				}
				last := strm[len(strm)-1].SeqNum
				if last != expected {
					return 0, esErrors.NewSeqExpectedErr(streamId, expected, last)
				}
				e.SeqNum = last + 1
				nspace[streamId] = append(strm, e)
				return e.SeqNum, nil
			case eventStore.NEW_STREAM:
				return 0, esErrors.NewStreamExists(streamId)
			}
		} else {
			switch cMode {
			case eventStore.ANY:
			case eventStore.NEW_STREAM:
				e.SeqNum = int64(len(strm))
				nspace[streamId] = append(strm, e)
				return e.SeqNum, nil
			default:
				return 0, esErrors.NewStreamDoesNotExist(streamId)
			}
		}
		panic("unreachable code reached")
	}

	var lastId int64
	for i, e := range events {
		last, err := insertEvent(e)
		if err != nil {
			// remove elements from the stream that were part of the batch and written (stored
			// in the map) before this error
			s := nspace[streamId]
			l := len(s)
			nspace[streamId] = s[:l-i]

			return 0, err
		}
		lastId = last
	}
	return lastId, nil
}

func (ms MemoryEventStore) GetEvent(ns string, streamId string,
	seqNum int64) (*es.EventEnvelope, error) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	if nspace, ok := ms.nss[ns]; ok {
		if stream, ok := nspace[streamId]; ok {
			for _, e := range stream {
				if e.SeqNum == seqNum {
					return &e, nil
				}
			}
			return nil, errors.New(fmt.Sprintf("Event sequence %v not found in stream %s, namespace %s",
				seqNum, streamId, ns))
		}
		return nil, errors.New(fmt.Sprintf("Stream %s not found in namespace %s", streamId, ns))
	}
	return nil, errors.New(fmt.Sprintf("Namespace %s not found", ns))
}

func (ms MemoryEventStore) GetEventRange(ns string, streamId string,
	starting int64, ending int64) ([]es.EventEnvelope, error) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	if nspace, ok := ms.nss[ns]; ok {
		if stream, ok := nspace[streamId]; ok {
			if starting < 0 {
				starting = 0
			}
			if ending >= int64(len(stream)) || ending < 0 || ending < starting {
				return stream[starting:], nil
			}
			return stream[starting : ending+1], nil
		}
		return nil, errors.New(fmt.Sprintf("Stream %s not found in namespace %s", streamId, ns))
	}
	return nil, errors.New(fmt.Sprintf("Namespace %s not found", ns))
}
