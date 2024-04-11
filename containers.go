package utility

import (
	"sync"
)

type Number interface {
	~int8 | ~int16 | ~int32 | ~int64 | ~int |
		~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uint |
		~float32 | ~float64
}

// Optional is a type represents either a wrapped value or `nil`, the absence of a value.
type Optional[Wrapped any] struct {
	some  Wrapped
	valid bool
}

// Get value if `valid` or zero value of Wrapped.
func (this Optional[Wrapped]) Get() Wrapped {
	return this.some
}

// IsValid of the value.
func (this Optional[Wrapped]) IsValid() bool {
	return this.valid
}

// Wrap nullable pointer as Optional, nil with NoneOptional, non-nil with SomeOptional.
func OptionalWithPointer[Wrapped any](v *Wrapped) Optional[Wrapped] {
	if v == nil {
		return NoneOptional[Wrapped]()
	}
	return SomeOptional(*v)
}

func SomeOptional[Wrapped any](v Wrapped) Optional[Wrapped] {
	return Optional[Wrapped]{some: v, valid: true}
}

func NoneOptional[Wrapped any]() Optional[Wrapped] {
	return Optional[Wrapped]{}
}

func MapOptional[Wrapped any, Return any](opt Optional[Wrapped], fn func(Wrapped) Return) Optional[Return] {
	if opt.valid && fn != nil {
		return SomeOptional(fn(opt.some))
	}
	return NoneOptional[Return]()
}

func FlatMapOptional[Wrapped any, Return any](opt Optional[Wrapped], fn func(Wrapped) Optional[Return]) Optional[Return] {
	if opt.valid && fn != nil {
		return fn(opt.some)
	}
	return NoneOptional[Return]()
}

// Range of any number `Bound`.
type Range[Bound Number] struct {
	lowerBound  Bound
	upperBound  Bound
	lowerClosed bool
	upperClosed bool
}

// Range with [lower, upper).
func NewUpperOpenedRange[Bound Number](lower, upper Bound) Range[Bound] {
	return Range[Bound]{
		lowerBound: lower, upperBound: upper,
		lowerClosed: true, upperClosed: false,
	}
}

// Range with [lower, upper].
func NewClosedRange[Bound Number](lower, upper Bound) Range[Bound] {
	return Range[Bound]{
		lowerBound: lower, upperBound: upper,
		lowerClosed: true, upperClosed: true,
	}
}

func (this Range[Bound]) WithLowerOpened() Range[Bound] {
	this.lowerClosed = false
	return this
}

func (this Range[Bound]) WithUpperOpened() Range[Bound] {
	this.upperClosed = false
	return this
}

func (this Range[Bound]) IsEmpty() bool {
	return this.lowerBound > this.upperBound ||
		(this.lowerBound == this.upperBound && (!this.lowerClosed || !this.upperClosed))
}

func (this Range[Bound]) Contains(v Bound) bool {
	if this.lowerClosed {
		if v < this.lowerBound {
			return false
		}
	} else if v <= this.lowerBound {
		return false
	}
	if this.upperClosed {
		if this.upperBound < v {
			return false
		}
	} else if this.upperBound <= v {
		return false
	}
	return true
}

// A value that represents either a success or a failure, including an associated value in each case.
type Result[Success any, Failure error] struct {
	success Optional[Success]
	failure Optional[Failure]

	isSuccess bool
}

func ResultWithSuccess[Success any, Failure error](v Success) Result[Success, Failure] {
	return Result[Success, Failure]{success: SomeOptional(v), failure: NoneOptional[Failure](), isSuccess: true}
}

func ResultWithFailure[Success any, Failure error](e Failure) Result[Success, Failure] {
	return Result[Success, Failure]{success: NoneOptional[Success](), failure: SomeOptional(e), isSuccess: false}
}

func (this Result[Success, Failure]) Success() Optional[Success] {
	return this.success
}

func (this Result[Success, Failure]) Failure() Optional[Failure] {
	return this.failure
}

func (this Result[Success, Failure]) IsSuccess() bool {
	return this.isSuccess
}

func MapResult[Success any, Failure error, NewSuccess any](
	result Result[Success, Failure],
	fn func(Success) NewSuccess,
) Result[NewSuccess, Failure] {
	if result.isSuccess {
		r := fn(result.success.some)
		return ResultWithSuccess[NewSuccess, Failure](r)
	}
	return ResultWithFailure[NewSuccess](result.failure.some)
}

func FlatMapResult[Success any, Failure error, NewSuccess any](
	result Result[Success, Failure],
	fn func(Success) Result[NewSuccess, Failure],
) Result[NewSuccess, Failure] {
	if result.isSuccess {
		return fn(result.success.some)
	}
	return ResultWithFailure[NewSuccess](result.failure.some)
}

// An unordered collection of unique elements.
type Set[Element comparable] struct {
	container map[Element]struct{}
}

func NewSet[Element comparable](values ...Element) Set[Element] {
	inst := Set[Element]{}
	inst.Add(values...)
	return inst
}

func (this *Set[Element]) Add(values ...Element) {
	if this.container == nil {
		this.container = make(map[Element]struct{}, len(values))
	}
	for _, value := range values {
		this.container[value] = struct{}{}
	}
}

func (this *Set[Element]) Remove(values ...Element) {
	if len(this.container) == 0 {
		return
	}
	for _, value := range values {
		delete(this.container, value)
	}
}

func (this Set[Element]) Contains(value Element) bool {
	_, exists := this.container[value]
	return exists
}

func (this Set[Element]) All() []Element {
	all := make([]Element, len(this.container))
	i := 0
	for key := range this.container {
		all[i] = key
		i++
	}
	return all
}

func (this *Set[Element]) Count() int {
	return len(this.container)
}

type RWLock[Value any] struct {
	value Value
	lock  sync.RWMutex
}

func NewRWLock[Value any](v Value) RWLock[Value] {
	return RWLock[Value]{value: v}
}

func (self *RWLock[Value]) Read() Value {
	self.lock.RLock()
	defer self.lock.RUnlock()
	return self.value
}

func (self *RWLock[Value]) ReadFn(op func(v *Value)) {
	self.lock.RLock()
	defer self.lock.RUnlock()
	op(&self.value)
}

func (self *RWLock[Value]) Write(v Value) {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.value = v
}

func (self *RWLock[Value]) WriteFn(op func(v *Value)) {
	self.lock.Lock()
	defer self.lock.Unlock()
	op(&self.value)
}
