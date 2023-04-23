package v3

import (
	"errors"
	"fmt"
	"sort"
)

// actions to be done on nodes of a segment tree
const (
	isBeg = 1
	isEnd = -1
)

// Keys returns keys of a map (in this project it is used to get the elements of the set)
func Keys[M ~map[K]V, K comparable, V any](m M) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}

// ordered represents types comparable with <, <=, >, >=, ==, !=
type ordered interface {
	~int | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64 | ~string
}

// binarySearch is the function of searching for an element in a slice of the ordered type,
// implemented through the upper-bound algorithm
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

// event is a struct for rectangle that contains
type event struct {
	x, b, e, stat int32
}

// newEvent is a constructor for event
func newEvent(x, lRange, rRange, status int32) event {
	return event{
		x:    x,
		b:    lRange,
		e:    rRange,
		stat: status,
	}
}

// point represents point in 2D-space
type point struct {
	x, y int32
}

func (p point) String() string {
	return fmt.Sprintf("(%d, %d)", p.x, p.y)
}

// rectangle represents two angles - left bottom and right upper angles
type rectangle struct {
	lb, ru point
}

func (r rectangle) String() string {
	return fmt.Sprintf("Left bottom: %d --- Right upper: %d", r.lb, r.ru)
}

// compressCoordinates replace original coordinates presented by points with their indices in sorted set
// returns compressed X-components and compressed Y-components
func (r *rectangles) compressCoordinates() ([]int32, []int32) {
	compXset := make(map[int32]struct{}, len(*r)*2)
	compYset := make(map[int32]struct{}, len(*r)*2)

	// creating a go-style set of coordinates and adding 1 to coordinates of right-upper angle of a rectangle
	// what is comfortable to use with upper-bound binary search
	for _, rect := range *r {
		compXset[rect.lb.x] = struct{}{}
		compXset[rect.ru.x+1] = struct{}{}
		compYset[rect.lb.y] = struct{}{}
		compYset[rect.ru.y+1] = struct{}{}
	}

	// convert set to slice
	compX := Keys(compXset)
	compY := Keys(compYset)
	// sort by value
	sort.Slice(compX, func(i, j int) bool {
		return compX[i] < compX[j]
	})
	sort.Slice(compY, func(i, j int) bool {
		return compY[i] < compY[j]
	})

	return compX, compY
}

// Node represents a single node in segment tree
type Node struct {
	left, right    *Node // descendants
	lRange, rRange int32 // show the segment which the node is responsible for
	sum            int32 // contains sum of range given
}

// NewNode is a constructor for Node
func NewNode(left, right *Node, li, ri, s int32) *Node {
	return &Node{
		left:   left,
		right:  right,
		lRange: li,
		rRange: ri,
		sum:    s,
	}
}

// addNode creates new condition of persistent segment tree and returns root of a new tree
func addNode(root *Node, li, ri, val int32) *Node {
	// node range is inside root
	// simply create new node
	if li <= root.lRange && ri >= root.rRange {
		return NewNode(root.left, root.right, root.lRange, root.rRange, root.sum+val)
	}

	// node range is outside the root
	if root.lRange >= ri || root.rRange <= li {
		return root
	}

	// node range is greater than root
	// create new root
	newRoot := NewNode(root.left, root.right, root.lRange, root.rRange, root.sum)

	// hang the descendants to a new root
	newRoot.left = addNode(newRoot.left, li, ri, val)
	newRoot.right = addNode(newRoot.right, li, ri, val)

	return newRoot
}

// buildEmpty builds empty tree on a given sample
// the construction takes place from top to bottom
func buildEmpty(sample []int32, li, ri int32) *Node {
	// leaf reached
	if ri-li == 1 {
		return NewNode(nil, nil, li, ri, sample[li])
	}

	var mid = (li + ri) / 2
	// left part is responsible for [li; mid)
	left := buildEmpty(sample, li, mid)
	// right part is responsible for [mid; ri)
	right := buildEmpty(sample, mid, ri)

	return NewNode(left, right, left.lRange, right.rRange, left.sum+right.sum)
}

// buildPersistentTree builds a persistent tree for rectangles type given
// returns []*Node, where every single element reflects the state of the tree at the i-th change
func buildPersistentTree(rects *rectangles, compX, compY []int32) []*Node {
	events := make([]event, 0, len(*rects)*2)
	persistentNodes := make([]*Node, 0, len(*rects)*2)

	for _, r := range *rects {
		// add begin of a rectangle to event
		events = append(events, newEvent(
			binarySearch(compX, r.lb.x),
			binarySearch(compY, r.lb.y),
			binarySearch(compY, r.ru.y+1),
			isBeg,
		))
		// add end of rectangle to event
		events = append(events, newEvent(
			binarySearch(compX, r.ru.x+1),
			binarySearch(compY, r.lb.y),
			binarySearch(compY, r.ru.y+1),
			isEnd,
		))
	}

	// sort by X-component
	sort.Slice(events, func(i, j int) bool {
		return events[i].x < events[j].x
	})

	sample := make([]int32, len(compY), len(compY))
	empty := buildEmpty(sample, 0, int32(len(compY)))

	x := events[0].x
	// process every x coordinate
	for _, e := range events {
		// processing finished - add new condition of a persistent tree
		if x != e.x {
			persistentNodes = append(persistentNodes, empty)
			// change x
			x = e.x
		}

		empty = addNode(empty, e.b, e.e, e.stat)
	}

	return persistentNodes
}

// searchInTree is a depth-first search (DFS) which finds sum of a given segment
func searchInTree(root *Node, val int32) int32 {
	// if is not a leaf
	if root != nil {
		mid := (root.lRange + root.rRange) / 2

		// since descendants are sorted from top to bottom and from left to right
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

// readData reads data from stdin in the following format:
// n - number of rectangles
// n lines with coordinates of rectangles
// x1 y1 x2 y2
// ...
// m - number of points to check
// m lines with coordinates of points
// x y
// ...
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

// Solve runs algorithm and prints result to stdout
func Solve() {
	rects, ps, err := readData()

	// no rectangles in input data
	if errors.Is(err, emptyRectangles) {
		for range ps {
			fmt.Print(0, " ")
		}
		return
	}
	// unexpected error
	if err != nil {
		fmt.Print(err)
		return
	}

	compX, compY := rects.compressCoordinates()
	tree := buildPersistentTree(&rects, compX, compY)

	var res []int32
	for _, p := range ps {
		// cX is a condition of a tree
		cX := binarySearch(compX, p.x)
		// is a value in a tree (point in compressed coordinates)
		cY := binarySearch(compY, p.y)

		// coordinates out of bound
		if cX == -1 || cY == -1 {
			res = append(res, 0)
			continue
		}
		res = append(res, searchInTree(tree[cX], cY))
	}

	for _, v := range res {
		fmt.Print(v, " ")
	}
}
