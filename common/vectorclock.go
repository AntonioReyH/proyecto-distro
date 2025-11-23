package common

import (
	pb "proyecto-sd/proto"
)

// MergeClocks combina dos relojes tomando el máximo de cada entrada.
func MergeClocks(c1, c2 *pb.VectorClock) *pb.VectorClock {
	merged := &pb.VectorClock{Versions: make(map[string]int64)}
	if c1 != nil {
		for k, v := range c1.Versions {
			merged.Versions[k] = v
		}
	}
	if c2 != nil {
		for k, v := range c2.Versions {
			if v > merged.Versions[k] {
				merged.Versions[k] = v
			}
		}
	}
	return merged
}

// IsAfter verifica si c1 es "posterior" o igual a c2 (Monotonic Check).
// Retorna true si todos los elementos de c1 >= c2.
func IsAfter(c1, c2 *pb.VectorClock) bool {
	if c2 == nil {
		return true
	}
	if c1 == nil {
		return false
	}
	for k, v2 := range c2.Versions {
		if c1.Versions[k] < v2 {
			return false
		}
	}
	return true
}