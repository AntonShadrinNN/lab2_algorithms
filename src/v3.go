package src

import (
	"sort"

	"lab2_algorithms/internal"
)

// actions to be done on nodes of a segment tree
const (
	isBeg = 1
	isEnd = -1
)

// node represents a single node in segment tree
type node struct {
	left, right    *node // descendants
	lRange, rRange int32 // show the segment which the node is responsible for
	sum            int32 // contains sum of range given
}

// newNode is a constructor for node
func newNode(left, right *node, li, ri, s int32) *node {
	return &node{
		left:   left,
		right:  right,
		lRange: li,
		rRange: ri,
		sum:    s,
	}
}

// addNode creates new condition of persistent segment tree and returns root of a new tree
func addNode(root *node, li, ri, val int32) *node {
	// node range is inside root
	// simply create new node
	if li <= root.lRange && ri >= root.rRange {
		return newNode(root.left, root.right, root.lRange, root.rRange, root.sum+val)
	}

	// node range is outside the root
	if root.lRange >= ri || root.rRange <= li {
		return root
	}

	// node range is greater than root
	// create new root
	newRoot := newNode(root.left, root.right, root.lRange, root.rRange, root.sum)

	// hang the descendants to a new root
	newRoot.left = addNode(newRoot.left, li, ri, val)
	newRoot.right = addNode(newRoot.right, li, ri, val)

	return newRoot
}

// buildEmpty builds empty tree on a given sample
// the construction takes place from top to bottom
func buildEmpty(sample []int32, li, ri int32) *node {
	// leaf reached
	if ri-li == 1 {
		return newNode(nil, nil, li, ri, sample[li])
	}

	var mid = (li + ri) / 2
	// left part is responsible for [li; mid)
	left := buildEmpty(sample, li, mid)
	// right part is responsible for [mid; ri)
	right := buildEmpty(sample, mid, ri)

	return newNode(left, right, left.lRange, right.rRange, left.sum+right.sum)
}

// BuildPersistentTree builds a persistent tree for rectangles type given
// returns []*node, where every single element reflects the state of the tree at the i-th change
func BuildPersistentTree(rects *internal.Rectangles, compX, compY []int32) []*node {
	events := make([]internal.Event, 0, len(*rects)*2)
	persistentNodes := make([]*node, 0, len(*rects)*2)

	for _, r := range *rects {
		// add begin of a rectangle to event
		events = append(events, internal.NewEvent(
			internal.BinarySearch(compX, r.Lb.X),
			internal.BinarySearch(compY, r.Lb.Y),
			internal.BinarySearch(compY, r.Ru.Y+1),
			isBeg,
		))
		// add end of rectangle to event
		events = append(events, internal.NewEvent(
			internal.BinarySearch(compX, r.Ru.X+1),
			internal.BinarySearch(compY, r.Lb.Y),
			internal.BinarySearch(compY, r.Ru.Y+1),
			isEnd,
		))
	}

	// sort by X-component
	sort.Slice(events, func(i, j int) bool {
		return events[i].X < events[j].X
	})

	sample := make([]int32, len(compY), len(compY))
	empty := buildEmpty(sample, 0, int32(len(compY)))

	x := events[0].X
	// process every x coordinate
	for _, e := range events {
		// processing finished - add new condition of a persistent tree
		if x != e.X {
			persistentNodes = append(persistentNodes, empty)
			// change x
			x = e.X
		}

		empty = addNode(empty, e.B, e.E, e.Stat)
	}

	return persistentNodes
}

// searchInTree is a depth-first search (DFS) which finds sum of a given segment
func searchInTree(root *node, val int32) int32 {
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

// PersistenceSegmentTree runs algorithm and prints result to stdout
func PersistenceSegmentTree(rects internal.Rectangles, ps internal.Points) {
	//rects, ps, err := internal.ReadData()

	// no rectangles in input data
	//if errors.Is(err, internalerrs.EmptyRectangles) {
	//	for range ps {
	//		fmt.Print(0, " ")
	//	}
	//	return
	//}
	//// unexpected error
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	compX, compY := rects.CompressCoordinates()
	tree := BuildPersistentTree(&rects, compX, compY)
	var res []int32
	for _, p := range ps {
		// cX is a condition of a tree
		cX := internal.BinarySearch(compX, p.X)
		// is a value in a tree (point in compressed coordinates)
		cY := internal.BinarySearch(compY, p.Y)

		// coordinates out of bound
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
