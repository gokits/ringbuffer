package ringstore

import (
	"testing"
)

func TestMemRingStore(t *testing.T) {
	strs := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8"}
	s1 := NewMemoryRingStore(7)
	for i := 0; i < 2; i++ {
		s1.Append(strs[i])
	}

	c1 := NewCursor()
	for i := 0; i < 2; i++ {
		d, st, err := s1.Get(c1)
		if err != nil {
			t.Fatalf("Get failed with err %v", err)
		}
		if st != Match {
			t.Fatalf("GetState unexpected %s", st)
		}
		if d != strs[i] {
			t.Fatalf("value not expected, expect: %s, actual: %s", strs[i], d)
		}
	}
	for i := 0; i < 2; i++ {
		_, st, err := s1.Get(c1)
		if err != nil {
			t.Fatalf("Get failed with err %v", err)
		}
		if st != Beyond {
			t.Fatalf("GetState unexpected. expect: Beyond, actual: %s", st)
		}
	}
	for i := 2; i < len(strs); i++ {
		if err := s1.Append(strs[i]); err != nil {
			t.Fatalf("append failed: %v", err)
		}
	}
	for i := 2; i < len(strs); i++ {
		d, st, err := s1.Get(c1)
		if err != nil {
			t.Fatalf("[%d] Get failed with err %v", i, err)
		}
		if st != Match {
			t.Fatalf("[%d] GetState unexpected. expect: Match, actual: %s", i, st)
		}
		if d != strs[i] {
			t.Fatalf("[%d] value not expected. expect: %s, actual: %s", i, strs[i], d)
		}
	}
}
