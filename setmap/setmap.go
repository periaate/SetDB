package setmap

import (
	"bytes"
	"sync"
)

var (
	increase float64 = 1.3
)

type capacity struct {
	Max     uint64
	Current uint64
}

// Setmap is a data structure with both hashmap and linked list functionality
type Setmap struct {
	Name      string
	Fp        string
	Sets      []*Setmap
	next      *Setmap
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
	return h % sm.Capacity.Max
}

// Newmap creates a new Setmap
func Newmap(size uint64) *Setmap {
	sm := new(Setmap)
	sm.Init(size)
	return sm
}

// Init initializes setmap instance
func (sm *Setmap) Init(size uint64) {
	sm.Capacity.Max = size
	sm.Sets = make([]*Setmap, size)
}

// Get finds a Setmap by name
func (sm *Setmap) Get(Name string) (*Setmap, bool) {
	key := sm.hash(Name)
	set := sm.Sets[key]
	if set != nil {
		return find(Name, set)
	}
	return nil, false
}

func find(Name string, sm *Setmap) (*Setmap, bool) {
	if sm.Name == Name || sm.next == nil {
		return sm, sm.Name == Name
	}
	return find(Name, sm.next)
}

func unset(Name string, smC *Setmap, smL *Setmap) bool {
	if smC == nil {
		return false
	}

	if smC.Name == Name {
		smL.next = smC.next
		return true
	}
	return unset(Name, smL.next, smC.next)
}

// Unset removes a value from the map, linked list
func (sm *Setmap) Unset(Name string) bool {
	key := sm.hash(Name)
	set := sm.Sets[key]
	if set == nil {
		return false
	}

	res := unset(set.Name, set, set)
	if res {
		sm.runcap(false)
	}
	return res
}

// New entry
func (sm *Setmap) New(sr *Setmap) {
	key := sm.hash(sr.Name)
	set := sm.Sets[key]

	if set == nil {
		sm.writeLock.Lock()
		sm.Sets[key] = sr
		sm.writeLock.Unlock()
		sm.runcap(true)
	} else {
		set.writeLock.Lock()
		set.add(sr)
		set.writeLock.Unlock()
		set.runcap(true)
	}
}

func (sm *Setmap) add(st *Setmap) {
	if sm.next == nil {
		sm.next = st
	} else {
		sm.next.add(st)
	}
}

func (sm *Setmap) runcap(b bool) {
	if b {
		sm.Capacity.Current++
	} else {
		sm.Capacity.Current--
	}
	if (float64(sm.Capacity.Current) / float64(sm.Capacity.Max)) > 0.8 {
		sm.writeLock.Lock()
		sm.resize()
	}
}

func (sm *Setmap) rehash(fSet *Setmap) bool {
	if fSet == nil {
		return false
	}
	if fSet.next != nil {
		if sm.rehash(fSet.next) {
			fSet.next = nil
		}
	}
	sm.New(fSet)
	return true
}

func (sm *Setmap) resize() {
	newcapacity := uint64((float64(sm.Capacity.Max) + 1) * increase)
	// TODO: verbose flag
	// fmt.Printf("capacity: %v.Currentcapacity: %v newcapacity: %v\n", sm.Capacity.Max, sm.Capacity.Current, newcapacity)

	oldSets := sm.Sets
	sm.Capacity.Max = newcapacity
	sm.Capacity.Current = 0
	sm.Sets = make([]*Setmap, newcapacity)
	sm.writeLock.Unlock()
	for i := 0; i < len(oldSets); i++ {
		sm.rehash(oldSets[i])
	}
}

// Stringify returns the names of keys of a map
func (sm *Setmap) Stringify() []string {
	sar := []string{}
	for _, v := range sm.Sets {
		if v == nil {
			continue
		}
		sar = append(sar, v.llAsString(v)...)
	}
	return sar
}

func (sm *Setmap) llAsString(s *Setmap) []string {
	if s.next == nil {
		return []string{s.Fp}
	}
	return append([]string{s.Fp}, s.llAsString(s.next)...)
}
