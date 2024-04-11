package utility

func IsSliceEqual[Element comparable](slice1, slice2 []Element) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	for i := 0; i < len(slice1); i++ {
		if slice1[i] != slice2[i] {
			return false
		}
	}
	return true
}

// FindInStrings would return first index of `slice` if found `item`, or -1 for not.
func FindInSlice[Element comparable](slice []Element, item Element) int {
	return FindFirstInSliceWhere(slice, func(e Element) bool {
		return e == item
	})
}

func FindFirstRefInSliceWhere[Element any](slice []Element, fn func(e Element) bool) *Element {
	var index = FindFirstInSliceWhere(slice, fn)
	if index < 0 {
		return nil
	}
	return &slice[index]
}

func FindFirstInSliceWhere[Element any](slice []Element, fn func(e Element) bool) int {
	count := len(slice)
	if count > 0 {
		for i := 0; i < count; i++ {
			if fn(slice[i]) {
				return i
			}
		}
	}
	return -1
}

func FilterSlice[Element any](slice []Element, fn func(e Element) bool) (ret []Element) {
	count := len(slice)
	if count > 0 {
		for i := 0; i < count; i++ {
			if fn(slice[i]) {
				ret = append(ret, slice[i])
			}
		}
	}
	return
}

func MapSlice[Element any, Return any](slice []Element, fn func(Element) Return) []Return {
	if slice == nil {
		return nil
	}
	count := len(slice)
	ret := make([]Return, count)
	for i := 0; i < count; i++ {
		ret[i] = fn(slice[i])
	}
	return ret
}

func ReversedSlice[Element any](slice []Element) []Element {
	if slice == nil {
		return nil
	}
	count := len(slice)
	ret := make([]Element, count)
	ri := 0
	for i := count - 1; i >= 0; i-- {
		ret[ri] = slice[i]
		ri += 1
	}
	return ret
}
