package setmap

import (
	"crypto"
	"encoding/hex"
	"fmt"
	"math"
	"testing"
	"time"
)

const (
	depth     float64 = 1.3
	tablesize         = 1299709 // 104395303 // 15485863 // 1299709
	l         float64 = tablesize * depth
)

var (
	length uint64 = uint64(math.Floor(l))
	sm     *Setmap
	Names  []string
)

func TestNodemap(t *testing.T) {
	Names = generate(int(length))
	sm = Newmap(length)
}

func TestHasher(t *testing.T) {
	var key uint64
	collision := 0
	percent := 0.0
	hashes := make([]uint64, length)
	start := time.Now()
	for _, Name := range Names {
		key = sm.hash(Name)
		if hashes[key] != 0 {
			collision++
			continue
		}
		hashes[key] = key
	}
	elapsed := time.Since(start)
	percent = collisionPercentage(uint64(collision), length)

	fmt.Printf("Percent collisions: %v%% || Time elapsed: %v\n", percent, elapsed)
}

func TestSet(t *testing.T) {
	for _, Name := range Names {
		set := Newmap(1)
		set.Name = Name
		sm.New(set)
	}
}

func TestGet(t *testing.T) {
	for _, Name := range Names {
		res, found := sm.Get(Name)
		if res == nil {
			t.Error("Name not found")
			continue
		}
		if found && res.Name != Name {
			t.Errorf("Names do not match received: %v  expected: %v", res.Name, Name)
		}
	}
}

func TestUnset(t *testing.T) {
	for _, Name := range Names {
		sm.Unset(Name)
	}
}

// Benchmarks depend on test setup

func BenchmarkSet(b *testing.B) {
	for _, Name := range Names {
		set := Newmap(0)
		set.Name = Name
		sm.New(set)
	}
}

func BenchmarkGet(b *testing.B) {
	for _, Name := range Names {
		sm.Get(Name)
	}
}

func generate(length int) []string {
	Names := make([]string, length)
	for i := range Names {
		h := crypto.SHA256.New()
		h.Write([]byte(fmt.Sprint(i)))
		Names[i] = hex.EncodeToString(h.Sum(nil))
	}
	return Names
}

func collisionPercentage(collisions uint64, length uint64) float64 {
	return float64(int(float64(collisions)/float64(length)*1000)) / 10
}
