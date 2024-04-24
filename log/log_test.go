package log

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	l0 := NewLogger(
		newZapLogger("I", LocalFormat{}, LevelInfo, nil),
		newZapLogger("D", LocalFormat{}, LevelDebug, nil),
	).WithTrace(NewTrace("A", time.Now(), Any("a", 1)))
	// 新logger的trace不会影响旧logger的trace
	l1 := l0.WithTraceLogs(Any("b", 2))
	len0 := len(l0.trace.pairs)
	len1 := len(l1.trace.pairs)
	assert.Equal(t, 2, len0)
	assert.Equal(t, 3, len1)

	assert.Equal(t, 1, len(l0.marchLevelOutputs(LevelDebug)))
	assert.Equal(t, len1+1, len(l1.producePairs(nil)))
	assert.Equal(t, len1+2, len(l1.producePairs([]LogPair{Any("C", 1)})))

	assert.Equal(t, len0, len(l0.trace.pairs))
}
