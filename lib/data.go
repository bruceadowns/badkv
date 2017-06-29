package lib

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// Key type definition
type Key string

// Value type definition
type Value []byte

// TimeStampValue manages a value
type TimeStampValue struct {
	data      Value
	timestamp time.Time
	tombstone bool
}

// Data is the in-memory map
type Data map[Key]*TimeStampValue

// Store defines the protected in-memory storage type
type Store struct {
	sync.RWMutex
	data Data
}

var store Store

func init() {
	store.data = make(Data)
}

func (v *TimeStampValue) String() string {
	return fmt.Sprintf(
		"data: %s timestamp: %d tombstone: %t",
		v.data, v.timestamp.UnixNano(), v.tombstone)
}

func (store *Store) get(k Key) ([]byte, error) {
	store.RLock()
	defer store.RUnlock()

	if v, ok := store.data[k]; ok {
		log.Printf("Found %s '%s'", k, v)

		if v.tombstone {
			return nil, fmt.Errorf("%s tombstoned", k)
		}

		return v.data, nil
	}

	return nil, fmt.Errorf("%s not found", k)
}

func (store *Store) put(k Key, tsv *TimeStampValue) (err error) {
	store.Lock()
	defer store.Unlock()

	if _, ok := store.data[k]; ok {
		log.Printf("Replace existing %s", k)
	} else {
		log.Printf("Add new %s", k)
	}

	store.data[k] = tsv
	return
}

func (store *Store) delete(k Key) (err error) {
	store.Lock()
	defer store.Unlock()

	if v, ok := store.data[k]; ok {
		log.Printf("Found %s '%s'", k, v)

		if v.tombstone {
			err = fmt.Errorf("%s tombstoned", k)
		} else {
			log.Printf("Tombstone %s", k)

			v.data = nil
			v.tombstone = true
		}
	} else {
		err = fmt.Errorf("%s not found", k)
	}

	return
}
