package internal

import (
	"fmt"
	internalerrs "lab2_algorithms/internal/errors"
)

// Keys returns Keys of a map (in this project it is used to get the elements of the set)
func Keys[M ~map[K]V, K comparable, V any](m M) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}

// BinarySearch is the function of searching for an element in a slice of the Ordered type,
// implemented through the upper-bound algorithm
func BinarySearch[S ~[]T, T Ordered](s S, v T) int32 {
	var b int32
	var e = int32(len(s))
	for b < e {
		mid := b + (e-b)/2
		if v >= s[mid] {
			b = mid + 1
		} else {
			e = mid
		}
	}
	return b - 1
}

// ReadData reads data from stdin in the following format:
// n - number of Rectangles
// n lines with coordinates of Rectangles
// x1 y1 x2 y2
// ...
// m - number of Points to check
// m lines with coordinates of Points
// x y
// ...
func ReadData() (Rectangles, Points, error) {
	var n, m int32
	var r Rectangles
	var p Points
	if _, err := fmt.Scanf("%d\n", &n); err != nil {
		return nil, nil, internalerrs.WrongDataFormat
	}
	for i := int32(0); i < n; i++ {
		var x1, y1, x2, y2 int32
		if _, err := fmt.Scanf("%d %d %d %d\n", &x1, &y1, &x2, &y2); err != nil {
			return nil, nil, internalerrs.WrongDataFormat
		}
		lb := Point{
			X: x1,
			Y: y1,
		}
		ru := Point{
			X: x2,
			Y: y2,
		}
		r = append(r,
			Rectangle{Lb: lb, Ru: ru})
	}

	if _, err := fmt.Scanf("%d\n", &m); err != nil {
		return nil, nil, internalerrs.WrongDataFormat
	}

	for i := int32(0); i < m; i++ {
		var x, y int32
		if _, err := fmt.Scanf("%d %d\n", &x, &y); err != nil {
			return nil, nil, internalerrs.WrongDataFormat
		}
		p = append(p, Point{x, y})
	}

	return r, p, nil
}
