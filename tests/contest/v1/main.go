package v1

import (
	"fmt"

	"sort"
)

type ordered interface {
	~int | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64 | ~string
}

func Keys[M ~map[K]V, K comparable, V any](m M) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}

func binarySearch[S ~[]T, T ordered](s S, v T) int32 {
	var b int32
	var e = int32(len(s))
	for e > 0 {
		cur := b
		mid := e / 2
		cur += mid
		if v >= s[cur] {
			b = cur + 1
			e -= mid + 1
		} else {
			e = mid
		}
	}
	return b - 1
}

type point struct {
	x, y int32
}

func (p point) String() string {
	return fmt.Sprintf("(%d, %d)", p.x, p.y)
}

type rectangle struct {
	lb, ru point
}

func (r rectangle) String() string {
	return fmt.Sprintf("%v - %v\n", r.lb, r.ru)
}

type rectangles []rectangle

func (r *rectangles) compressCoordinates() ([]int32, []int32) {
	compXset := make(map[int32]struct{}, len(*r)*2)
	compYset := make(map[int32]struct{}, len(*r)*2)

	//compXset := make(map[int32]struct{})
	//compYset := make(map[int32]struct{})

	for _, rect := range *r {
		compXset[rect.lb.x] = struct{}{}
		compXset[rect.ru.x+1] = struct{}{}
		compYset[rect.lb.y] = struct{}{}
		compYset[rect.ru.y+1] = struct{}{}
	}

	compX := Keys(compXset)
	compY := Keys(compYset)
	sort.Slice(compX, func(i, j int) bool {
		return compX[i] < compX[j]
	})
	sort.Slice(compY, func(i, j int) bool {
		return compY[i] < compY[j]
	})

	return compX, compY
}

func (r *rectangles) buildMap() ([]int32, []int32, [][]int32) {
	compX, compY := r.compressCoordinates()
	rectMap := make([][]int32, len(compY))
	for i := range rectMap {
		rectMap[i] = make([]int32, len(compX))
	}

	for _, rect := range *r {
		lX := binarySearch(compX, rect.lb.x)
		lY := binarySearch(compY, rect.lb.y)
		rX := binarySearch(compX, rect.ru.x+1)
		rY := binarySearch(compY, rect.ru.y+1)

		for i := lY; i < rY; i++ {
			for j := lX; j < rX; j++ {
				rectMap[i][j]++
			}
		}
	}
	return compX, compY, rectMap
}

type points []point

var wrongDataFormat = fmt.Errorf("wrong data format received during data processing")

func readData() (rectangles, points, error) {
	var n, m int32
	var r rectangles
	var p points
	if _, err := fmt.Scanf("%d\n", &n); err != nil {
		return nil, nil, wrongDataFormat
	}
	for i := int32(0); i < n; i++ {
		var x1, y1, x2, y2 int32
		if _, err := fmt.Scanf("%d %d %d %d\n", &x1, &y1, &x2, &y2); err != nil {
			return nil, nil, wrongDataFormat
		}
		lb := point{
			x: x1,
			y: y1,
		}
		ru := point{
			x: x2,
			y: y2,
		}
		r = append(r,
			rectangle{lb: lb, ru: ru})
	}

	if _, err := fmt.Scanf("%d\n", &m); err != nil {
		return nil, nil, wrongDataFormat
	}

	for i := int32(0); i < m; i++ {
		var x, y int32
		if _, err := fmt.Scanf("%d %d\n", &x, &y); err != nil {
			return nil, nil, wrongDataFormat
		}
		p = append(p, point{x, y})
	}

	return r, p, nil
}

func Solve(rects rectangles, p points) {
	//rects, p, err := readData()
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	var res []int32
	compX, compY, rectMap := rects.buildMap()
	for i := range p {
		pX := binarySearch(compX, p[i].x)
		pY := binarySearch(compY, p[i].y)

		if pX == -1 || pY == -1 {
			res = append(res, 0)
			continue
		}
		res = append(res, rectMap[pY][pX])
	}

	//for _, v := range res {
	//	fmt.Print(v, " ")
	//}
	//fmt.Println()
}
