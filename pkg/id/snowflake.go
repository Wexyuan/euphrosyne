package id

import (
	"fmt"

	"github.com/bwmarrin/snowflake"
)

// Snowflake generates unique IDs from a single worker node.
type Snowflake struct {
	node *snowflake.Node
}

// NewSnowflake creates a generator for the worker node number (0-1023).
func NewSnowflake(node int64) (*Snowflake, error) {
	if node < 0 || node > 1023 {
		return nil, fmt.Errorf("[id] node %d out of range: must be between 0 and 1023", node)
	}
	n, err := snowflake.NewNode(node)
	if err != nil {
		return nil, fmt.Errorf("[id] create snowflake node %d error: %w", node, err)
	}
	return &Snowflake{node: n}, nil
}

// NextID returns a new unique ID as int64.
func (s *Snowflake) NextID() int64 {
	return s.node.Generate().Int64()
}

// NextString returns a new unique ID as string.
func (s *Snowflake) NextString() string {
	return s.node.Generate().String()
}
