package internal

import (
	"fmt"
	"sort"
	"sync"
)

// Ordered represents types comparable with <, <=, >, >=, ==, !=
type Ordered interface {
	~int | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64 | ~string
}

// Point represents Point in 2D-space
type Point struct {
	X, Y int32
}

func (p Point) String() string {
	return fmt.Sprintf("(%d, %d)", p.X, p.Y)
}

// Rectangle represents two angles - left bottom and right upper angles
type Rectangle struct {
	Lb, Ru Point
}

// Belongs checks if Point belongs to a Rectangle
func (r Rectangle) Belongs(p Point) bool {
	return (p.X >= r.Lb.X && p.X <= r.Ru.X) && (p.Y >= r.Lb.Y && p.Y <= r.Ru.Y)
}

func (r Rectangle) String() string {
	return fmt.Sprintf("%v - %v\n", r.Lb, r.Ru)
}

type Rectangles []Rectangle
type Points []Point

func (r *Rectangles) BuildMap() ([]int32, []int32, [][]int32) {
	compX, compY := r.CompressCoordinates()
	rectMap := make([][]int32, len(compY))
	for i := range rectMap {
		rectMap[i] = make([]int32, len(compX))
	}

	for _, rect := range *r {
		lX := BinarySearch(compX, rect.Lb.X)
		lY := BinarySearch(compY, rect.Lb.Y)
		rX := BinarySearch(compX, rect.Ru.X+1)
		rY := BinarySearch(compY, rect.Ru.Y+1)

		for i := lY; i < rY; i++ {
			for j := lX; j < rX; j++ {
				rectMap[i][j]++
			}
		}
	}
	return compX, compY, rectMap
}

// CompressCoordinates replace original coordinates presented by Points with their indices in sorted set
// returns compressed X-components and compressed Y-components
func (r *Rectangles) CompressCoordinates() ([]int32, []int32) {
	compXset := make(map[int32]struct{}, len(*r)*2)
	compYset := make(map[int32]struct{}, len(*r)*2)

	for _, rect := range *r {
		compXset[rect.Lb.X] = struct{}{}
		compXset[rect.Ru.X+1] = struct{}{}
		compYset[rect.Lb.Y] = struct{}{}
		compYset[rect.Ru.Y+1] = struct{}{}
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

type Event struct {
	X, B, E, Stat int32
}

// NewEvent is a constructor for Event
func NewEvent(x, lRange, rRange, status int32) Event {
	return Event{
		X:    x,
		B:    lRange,
		E:    rRange,
		Stat: status,
	}
}

type TestPair struct {
	m sync.Mutex
	N int32
	T int64
}
