// Package slices provides implementations of some built-in slice functions.
package slices

// CopyBytes copies bytes from src to dst. src and dst may overlap. CopyBytes
// returns the number of elements copied, which will be the minimum of len(dst)
// and len(src).
func CopyBytes(dst, src []byte) int {
	n := len(dst)
	if len(src) < n {
		n = len(src)
	}
	for i := 0; i < n; i++ {
		dst[i] = src[i]
	}
	return n
}

// AppendBytes append b to s and returns the resulting slice. If s has
// sufficient capacity, it is resliced to accommodate b. If it doesn't, a new
// slice is allocated.
func AppendBytes(s []byte, b byte) []byte {
	if len(s) == cap(s) {
		newCap := 2 * cap(s)
		if newCap == 0 {
			newCap = 1
		}
		newSlice := make([]byte, len(s), newCap)
		CopyBytes(newSlice, s)
		s = newSlice
	}
	s = s[:len(s)+1]
	s[len(s)-1] = b
	return s
}
