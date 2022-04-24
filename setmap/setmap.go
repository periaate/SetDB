package setmap

import (
	"bytes"
	"sync"
)

const increase = 1.3

type capacity struct {
	Max     uint64
	Current uint64
}

// Setmap is a data structure with both hashmap and linked list functionality
type Setmap struct {
	Name      string
	Fp        string
	Sets      []*Setmap
	Capacity  capacity
	writeLock sync.Mutex
}

// hash is a modified DJB2 that runs roughly 50% faster for SHA256, while giving identical results.
func (sm *Setmap) hash(Name string) uint64 {
	var (
		buf bytes.Buffer
		n1  uint64
		n2  uint64
	)

	buf.WriteString(Name)
	b := buf.Bytes()

	var h uint64 = 5381
	for i := 0; i < len(b)/4; i += 2 {
		n1 = uint64(b[i])
		n2 = uint64(b[i+2])
		h += (h*33 + n1)
		h += (h*33 + n2)
	}

	buf.Reset()
	return h
}

// Init initializes setmap instance
func (sm *Setmap) Init(size uint64) {
	sm.Capacity.Max = size
	sm.Sets = make([]*Setmap, size)
}

// Get finds a Setmap by name
func (sm *Setmap) Get(name string) *Setmap {
	keyOr := sm.hash(name)
	key := keyOr % sm.Capacity.Max
	var i uint64 = 1
	for sm.Sets[key] != nil {
		if sm.Sets[key].Name == name {
			return sm.Sets[key]
		}
		key = (keyOr + i) % sm.Capacity.Max
		i = i + 1
	}
	return nil
}

// New entry
func (sm *Setmap) New(sr *Setmap) {
	keyOr := sm.hash(sr.Name)
	key := keyOr % sm.Capacity.Max

	var i uint64 = 1
	for sm.Sets[key] != nil {

		if sm.Sets[key].Name == sr.Name {
			return
		}
		key = (keyOr + i) % sm.Capacity.Max
		i = i + 1
	}

	sm.writeLock.Lock()
	sm.Sets[key] = sr
	sm.writeLock.Unlock()
	sm.runcap(true)
}

func (sm *Setmap) runcap(b bool) {
	if b {
		sm.Capacity.Current++
	} else {
		sm.Capacity.Current--
	}
	if (float64(sm.Capacity.Current) / float64(sm.Capacity.Max)) > 0.65 {
		sm.writeLock.Lock()
		sm.resize()
	}
}

func (sm *Setmap) resize() {
	newcapacity := uint64(float64(sm.Capacity.Max+1) * increase)
	// TODO: verbose flag
	// fmt.Printf("capacity: %v.Currentcapacity: %v newcapacity: %v\n", sm.Capacity.Max, sm.Capacity.Current, newcapacity)

	oldSets := sm.Sets
	sm.Capacity.Max = newcapacity
	sm.Capacity.Current = 0
	sm.Sets = make([]*Setmap, newcapacity)
	sm.writeLock.Unlock()
	for i := 0; i < len(oldSets); i++ {
		if oldSets[i] != nil {
			sm.New(oldSets[i])
		}
	}
}

// Stringify returns the names of keys of a map
func (sm *Setmap) Stringify(name bool) []string {
	sar := []string{}
	for _, v := range sm.Sets {
		if v == nil {
			continue
		}
		if name {
			sar = append(sar, v.Name)
			continue
		}
		sar = append(sar, v.Fp)
	}
	return sar
}
