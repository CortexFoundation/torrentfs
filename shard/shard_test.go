package shard

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

type keyT = string
type valueT = any

func k(key int) keyT {
	return strconv.FormatInt(int64(key), 10)
}

func add(x keyT, delta int) int {
	i, err := strconv.ParseInt(x, 10, 64)
	if err != nil {
		panic(err)
	}
	return int(i + int64(delta))
}

// /////////////////////////
func random(N int, perm bool) []keyT {
	nums := make([]keyT, N)
	if perm {
		for i, x := range rand.Perm(N) {
			nums[i] = k(x)
		}
	} else {
		m := make(map[keyT]bool)
		for len(m) < N {
			m[k(int(rand.Uint64()))] = true
		}
		var i int
		for k := range m {
			nums[i] = k
			i++
		}
	}
	return nums
}

func shuffle(nums []keyT) {
	for i := range nums {
		j := rand.Intn(i + 1)
		nums[i], nums[j] = nums[j], nums[i]
	}
}

func init() {
	//var seed int64 = 1519776033517775607
	seed := (time.Now().UnixNano())
	println("seed:", seed)
	rand.Seed(seed)
}

func TestRandomData(t *testing.T) {
	t.Parallel()
	N := 10000
	start := time.Now()
	for time.Since(start) < time.Second*2 {
		nums := random(N, true)
		var m = New[int](0)
		switch rand.Int() % 5 {
		default:
			m = New[int](N / ((rand.Int() % 3) + 1))
		case 1:
			m = New[int](0) //new(Map[int])
		case 2:
			m = New[int](0)
		}
		v, ok := m.Get(k(999))
		if ok {
			t.Fatalf("expected %v, got %v", nil, v)
		}
		v, ok = m.Delete(k(999))
		if ok {
			t.Fatalf("expected %v, got %v", nil, v)
		}
		if m.Len() != 0 {
			t.Fatalf("expected %v, got %v", 0, m.Len())
		}
		// set a bunch of items
		for i := 0; i < len(nums); i++ {
			x, _ := strconv.Atoi(nums[i])
			v, ok := m.Set(nums[i], x)
			if ok {
				t.Fatalf("expected %v, got %v", nil, v)
			}
		}
		if m.Len() != N {
			t.Fatalf("expected %v, got %v", N, m.Len())
		}
		// retrieve all the items
		shuffle(nums)
		for i := 0; i < len(nums); i++ {
			v, ok := m.Get(nums[i])
			x, _ := strconv.Atoi(nums[i])
			if !ok || v != x {
				t.Fatalf("expected %v, got %v", nums[i], v)
			}
		}
		// replace all the items
		shuffle(nums)
		for i := 0; i < len(nums); i++ {
			v, ok := m.Set(nums[i], add(nums[i], 1))
			x, _ := strconv.Atoi(nums[i])
			if !ok || v != x {
				t.Fatalf("expected %v, got %v, %v", nums[i], v, ok)
			}
		}
		if m.Len() != N {
			t.Fatalf("expected %v, got %v", N, m.Len())
		}
		// retrieve all the items
		shuffle(nums)
		for i := 0; i < len(nums); i++ {
			v, ok := m.Get(nums[i])
			if !ok || v != add(nums[i], 1) {
				t.Fatalf("expected %v, got %v", add(nums[i], 1), v)
			}
		}
		// remove half the items
		shuffle(nums)
		for i := 0; i < len(nums)/2; i++ {
			v, ok := m.Delete(nums[i])
			if !ok || v != add(nums[i], 1) {
				t.Fatalf("expected %v, got %v", add(nums[i], 1), v)
			}
		}
		if m.Len() != N/2 {
			t.Fatalf("expected %v, got %v", N/2, m.Len())
		}
		// check to make sure that the items have been removed
		for i := 0; i < len(nums)/2; i++ {
			v, ok := m.Get(nums[i])
			if ok {
				t.Fatalf("expected %v, got %v", nil, v)
			}
		}
		// check the second half of the items
		for i := len(nums) / 2; i < len(nums); i++ {
			v, ok := m.Get(nums[i])
			if !ok || v != add(nums[i], 1) {
				t.Fatalf("expected %v, got %v", add(nums[i], 1), v)
			}
		}
		// try to delete again, make sure they don't exist
		for i := 0; i < len(nums)/2; i++ {
			v, ok := m.Delete(nums[i])
			if ok {
				t.Fatalf("expected %v, got %v", nil, v)
			}
		}
		if m.Len() != N/2 {
			t.Fatalf("expected %v, got %v", N/2, m.Len())
		}
		m.Range(func(key keyT, value int) bool {
			if value != add(key, 1) {
				t.Fatalf("expected %v, got %v", add(key, 1), value)
			}
			return true
		})
		var n int
		m.Range(func(key keyT, value int) bool {
			n++
			return false
		})
		if n != 1 {
			t.Fatalf("expected %v, got %v", 1, n)
		}
		for i := len(nums) / 2; i < len(nums); i++ {
			v, ok := m.Delete(nums[i])
			if !ok || v != add(nums[i], 1) {
				t.Fatalf("expected %v, got %v", add(nums[i], 1), v)
			}
		}
	}
}

func TestSetAccept(t *testing.T) {
	t.Parallel()
	var m = New[string](0)
	m.Set("hello", "world")
	prev, replaced := m.SetAccept("hello", "planet", nil)
	if !replaced {
		t.Fatal("expected true")
	}
	if prev != "world" {
		t.Fatalf("expected '%v', got '%v'", "world", prev)
	}
	if v, _ := m.Get("hello"); v != "planet" {
		t.Fatalf("expected '%v', got '%v'", "planet", v)
	}
	prev, replaced = m.SetAccept("hello", "world", func(prev string, replaced bool) bool {
		if !replaced {
			t.Fatal("expected true")
		}
		if prev != "planet" {
			t.Fatalf("expected '%v', got '%v'", "planet", prev)
		}
		return true
	})
	if !replaced {
		t.Fatal("expected true")
	}
	if prev != "planet" {
		t.Fatalf("expected '%v', got '%v'", "planet", prev)
	}
	prev, replaced = m.SetAccept("hello", "planet", func(prev string, replaced bool) bool {
		if !replaced {
			t.Fatal("expected true")
		}
		if prev != "world" {
			t.Fatalf("expected '%v', got '%v'", "world", prev)
		}
		return false
	})
	if replaced {
		t.Fatal("expected false")
	}
	//if prev != nil {
	//	t.Fatalf("expected '%v', got '%v'", nil, prev)
	//}
	if v, _ := m.Get("hello"); v != "world" {
		t.Fatalf("expected '%v', got '%v'", "world", v)
	}

	prev, replaced = m.SetAccept("hi", "world", func(prev string, replaced bool) bool {
		if replaced {
			t.Fatal("expected false")
		}
		//if prev != nil {
		//	t.Fatalf("expected '%v', got '%v'", nil, prev)
		//}
		return false
	})
	if replaced {
		t.Fatal("expected false")
	}
	//if prev != nil {
	//	t.Fatalf("expected '%v', got '%v'", nil, prev)
	//}
}

func TestDeleteAccept(t *testing.T) {
	t.Parallel()
	var m = New[string](0)
	m.Set("hello", "world")
	prev, deleted := m.DeleteAccept("hello", nil)
	if !deleted {
		t.Fatal("expected true")
	}
	if prev != "world" {
		t.Fatalf("expected '%v', got '%v'", "world", prev)
	}
	m.Set("hello", "world")
	prev, deleted = m.DeleteAccept("hello", func(prev string, deleted bool) bool {
		if !deleted {
			t.Fatal("expected true")
		}
		if prev != "world" {
			t.Fatalf("expected '%v', got '%v'", "world", prev)
		}
		return true
	})
	if !deleted {
		t.Fatal("expected true")
	}
	if prev != "world" {
		t.Fatalf("expected '%v', got '%v'", "world", prev)
	}
	m.Set("hello", "world")
	prev, deleted = m.DeleteAccept("hello", func(prev string, deleted bool) bool {
		if !deleted {
			t.Fatal("expected true")
		}
		if prev != "world" {
			t.Fatalf("expected '%v', got '%v'", "world", prev)
		}
		return false
	})
	if deleted {
		t.Fatal("expected false")
	}
	//if prev != nil {
	//	t.Fatalf("expected '%v', got '%v'", nil, prev)
	//}
	prev, ok := m.Get("hello")
	if !ok {
		t.Fatal("expected true")
	}
	if prev != "world" {
		t.Fatalf("expected '%v', got '%v'", "world", prev)
	}

}

func TestClear(t *testing.T) {
	t.Parallel()
	var m = New[int](0)
	for i := 0; i < 1000; i++ {
		m.Set(fmt.Sprintf("%d", i), i)
	}
	if m.Len() != 1000 {
		t.Fatalf("expected '%v', got '%v'", 1000, m.Len())
	}
	m.Clear()
	if m.Len() != 0 {
		t.Fatalf("expected '%v', got '%v'", 0, m.Len())
	}

}
