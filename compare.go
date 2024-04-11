package utility

func MaxInt64(a, b int64) int64 {
	return Max(a, b)
}

func MinInt64(a, b int64) int64 {
	return Min(a, b)
}

func Max[Value Number](l, r Value) Value {
	if l > r {
		return l
	}
	return r
}

func Min[Value Number](l, r Value) Value {
	if l < r {
		return l
	}
	return r
}
