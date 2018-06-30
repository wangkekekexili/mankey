package object

import "testing"

func TestHashKey(t *testing.T) {
	s1 := &String{Value: "42"}
	s2 := &String{Value: "42"}
	n1 := &Integer{Value: 42}
	n2 := &Integer{Value: 42}

	if s1.HashKey() != s2.HashKey() {
		t.Fatalf("expected s1 and s2 to have the same hash key")
	}
	if n1.HashKey() != n2.HashKey() {
		t.Fatalf("expected n1 and n2 to have the same hash key")
	}
	if s1.HashKey() == n1.HashKey() {
		t.Fatalf("expected s1 and n1 to have different hash keys")
	}
}
