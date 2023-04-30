package v3

import (
	"fmt"
	"sort"
)

const (
	isBeg = 1
	isEnd = -1
)

func Keys[M ~map[K]V, K comparable, V any](m M) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}

type ordered interface {
	~int | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64 | ~string
}

func binarySearch[S ~[]T, T ordered](s S, v T) int32 {
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

type event struct {
	x, b, e, stat int32
}

func newEvent(x, lRange, rRange, status int32) event {
	return event{
		x:    x,
		b:    lRange,
		e:    rRange,
		stat: status,
	}
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
	return fmt.Sprintf("Left bottom: %d --- Right upper: %d", r.lb, r.ru)
}

func (r *rectangles) compressCoordinates() ([]int32, []int32) {
	compXset := make(map[int32]struct{}, len(*r)*2)
	compYset := make(map[int32]struct{}, len(*r)*2)

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

type Node struct {
	left, right    *Node
	lRange, rRange int32
	sum            int32
}

func NewNode(left, right *Node, li, ri, s int32) *Node {
	return &Node{
		left:   left,
		right:  right,
		lRange: li,
		rRange: ri,
		sum:    s,
	}
}

func addNode(root *Node, li, ri, val int32) *Node {
	if li <= root.lRange && ri >= root.rRange {
		return NewNode(root.left, root.right, root.lRange, root.rRange, root.sum+val)
	}

	if root.lRange >= ri || root.rRange <= li {
		return root
	}

	toAdd := NewNode(root.left, root.right, root.lRange, root.rRange, root.sum)

	toAdd.left = addNode(toAdd.left, li, ri, val)
	toAdd.right = addNode(toAdd.right, li, ri, val)

	return toAdd
}

func buildEmpty(sample []int32, li, ri int32) *Node {
	if ri-li == 1 {
		return NewNode(nil, nil, li, ri, sample[li])
	}

	var mid = (li + ri) / 2
	left := buildEmpty(sample, li, mid)
	right := buildEmpty(sample, mid, ri)

	return NewNode(left, right, left.lRange, right.rRange, left.sum+right.sum)
}

func buildPersistentTree(rects *rectangles, compX, compY []int32) []*Node {
	events := make([]event, 0, len(*rects)*2)
	persistentNodes := make([]*Node, 0, len(*rects)*2)

	for _, r := range *rects {
		events = append(events, newEvent(
			binarySearch(compX, r.lb.x),
			binarySearch(compY, r.lb.y),
			binarySearch(compY, r.ru.y+1),
			isBeg,
		))
		events = append(events, newEvent(
			binarySearch(compX, r.ru.x+1),
			binarySearch(compY, r.lb.y),
			binarySearch(compY, r.ru.y+1),
			isEnd,
		))
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].x < events[j].x
	})

	sample := make([]int32, len(compY), len(compY))
	empty := buildEmpty(sample, 0, int32(len(compY)))

	x := events[0].x

	for _, e := range events {
		if x != e.x {
			persistentNodes = append(persistentNodes, empty)
			x = e.x
		}

		empty = addNode(empty, e.b, e.e, e.stat)
	}

	return persistentNodes
}

func searchInTree(root *Node, val int32) int32 {
	if root != nil {
		mid := (root.lRange + root.rRange) / 2

		if val < mid {
			return root.sum + searchInTree(root.left, val)
		} else {
			return root.sum + searchInTree(root.right, val)
		}
	}
	return 0
}

type rectangles []rectangle
type points []point

var wrongDataFormat = fmt.Errorf("wrong data format received during data processing")
var emptyRectangles = fmt.Errorf("no rectangles were given")

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
	if len(r) == 0 {
		return nil, p, emptyRectangles
	}
	return r, p, nil
}

func Solve(rects rectangles, ps points) {
	//rects, ps, err := readData()
	//
	//if errors.Is(err, emptyRectangles) {
	//	for _ = range ps {
	//		fmt.Print(0, " ")
	//	}
	//	return
	//}
	//if err != nil {
	//	fmt.Print(err)
	//	return
	//}
	compX, compY := rects.compressCoordinates()
	tree := buildPersistentTree(&rects, compX, compY)
	var res []int32
	for _, p := range ps {
		cX := binarySearch(compX, p.x)
		cY := binarySearch(compY, p.y)

		if cX == -1 || cY == -1 {
			res = append(res, 0)
			continue
		}
		res = append(res, searchInTree(tree[cX], cY))
	}
	//for _, v := range res {
	//	fmt.Print(v, " ")
	//}
}
