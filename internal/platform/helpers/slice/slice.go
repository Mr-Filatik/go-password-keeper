package slicehelpers

func ApplyIf[T any](slice []T, pred func(T) bool, f func(T)) {
	for _, v := range slice {
		if pred(v) {
			f(v)
		}
	}
}

func Reverse[T any](s []T) []T {
	reversed := make([]T, len(s))
	for i, j := 0, len(s)-1; i <= j; i, j = i+1, j-1 {
		reversed[i], reversed[j] = s[j], s[i]
	}
	return reversed
}
