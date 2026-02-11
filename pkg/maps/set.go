// Package maps provides map-based data structures and utilities.
//
// This file contains a Set implementation built on top of the Map type,
// providing a simple set interface for comparable key types.
package maps

type Set[K comparable] interface {
	Has(k K) bool
}

func SetOf[K comparable](a ...K) Set[K] {
	s := Map[K, interface{}]{}
	for _, v := range a {
		s.Put(v, nil)
	}
	return s
}
