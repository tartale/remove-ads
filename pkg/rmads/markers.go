package rmads

import (
	"time"
)

type Markers struct {
	Segments   Segments   `json:"segments,omitempty"`
	Timestamps Timestamps `json:"timestamps,omitempty"`
}

type Segment struct {
	Description string        `json:"description,omitempty"`
	StartOffset time.Duration `json:"startOffset,omitempty"`
	EndOffset   time.Duration `json:"endOffset,omitempty"`
}

type Segments []Segment

type Timestamp struct {
	Description string        `json:"description,omitempty"`
	Timestamp   time.Duration `json:"timestamp,omitempty"`
}

type Timestamps []Timestamp
