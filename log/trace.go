package log

import (
	"encoding/hex"
	"time"

	"github.com/google/uuid"
)

type Trace struct {
	startTime time.Time
	pairs     []LogPair
}

func (t Trace) length() int { return len(t.pairs) }

func (t Trace) isEmpty() bool {
	return t.startTime.IsZero() && len(t.pairs) == 0
}

func (t Trace) merge(other Trace) Trace {
	if t.startTime.IsZero() {
		t.startTime = other.startTime
	}
	if t.pairs == nil {
		t.pairs = other.pairs
	} else {
		t.pairs = append(t.pairs, other.pairs...)
	}
	return t
}

func NewTrace(traceID string, t time.Time, pairs ...LogPair) Trace {
	if traceID != "" {
		pairs = append(pairs, Any(fieldTraceID, traceID))
	}
	return Trace{startTime: t, pairs: pairs}
}

type TraceID struct {
	v [16]byte
}

func NewTraceID() TraceID {
	t := TraceID{v: uuid.New()}
	return t
}

func (t TraceID) String() string {
	return hex.EncodeToString(t.v[:])
}
