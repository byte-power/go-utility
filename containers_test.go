package utility

import (
	"errors"
	"sort"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptional(t *testing.T) {
	type Foo struct {
		name string
	}
	var opt = SomeOptional(Foo{name: "F"})
	assert.True(t, opt.IsValid())
	assert.Equal(t, opt.Get().name, "F")
	assert.Equal(t, MapOptional(opt, func(w Foo) int {
		return len(w.name)
	}).Get(), 1)
	assert.Equal(t, FlatMapOptional(opt, func(w Foo) Optional[int] {
		return SomeOptional(len(w.name))
	}).Get(), 1)

	opt = NoneOptional[Foo]()
	assert.False(t, opt.IsValid())
	assert.Equal(t, opt.Get().name, "")
	assert.Equal(t, MapOptional(opt, func(w Foo) int {
		return len(w.name)
	}).Get(), 0)
	assert.Equal(t, FlatMapOptional(opt, func(w Foo) Optional[int] {
		return SomeOptional(len(w.name))
	}).Get(), 0)
}

func TestResult(t *testing.T) {
	var result = ResultWithSuccess[int, error](1)
	assert.True(t, result.IsSuccess())
	assert.Equal(t, MapResult(result, func(s int) string {
		return strconv.Itoa(s)
	}), ResultWithSuccess[string, error]("1"))

	var err = errors.New("1")
	result = ResultWithFailure[int](err)
	assert.False(t, result.IsSuccess())
	assert.Equal(t, MapResult(result, func(s int) string {
		return strconv.Itoa(s)
	}), ResultWithFailure[string](err))
}

func TestSet(t *testing.T) {
	var ints = NewSet[int]()
	ints.Add(1, 2, 1)
	assert.Equal(t, ints, NewSet(1, 2, 1))
	assert.Equal(t, ints.Count(), 2)

	var arr = ints.All()
	sort.Ints(arr)
	assert.Equal(t, arr, []int{1, 2})

	ints.Remove(10)

	arr = ints.All()
	sort.Ints(arr)
	assert.Equal(t, arr, []int{1, 2})

	ints.Remove(1)

	arr = ints.All()
	sort.Ints(arr)
	assert.Equal(t, arr, []int{2})
}
