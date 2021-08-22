package graph

import (
	"fmt"
)

type Edge struct {
	V int
	W int
	Value float32
}

func (e *Edge) String() string { 
	return fmt.Sprintf("(%d, %d, %f)", e.V, e.W, e.Value) 
}

func (e *Edge) Less(other *Edge) bool {
	return e.Value < other.Value
}