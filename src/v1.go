package src

import (
	"fmt"
	"lab2_algorithms/internal"
)

func checkSinglePoint(p internal.Point, rects *internal.Rectangles) int32 {
	var count int32
	for _, r := range *rects {
		if r.Belongs(p) {
			count++
		}
	}
	return count
}

func BruteForce() {
	r, ps, err := internal.ReadData()
	if err != nil {
		fmt.Println(err)
		return
	}
	var res []int32
	for _, p := range ps {
		res = append(res, checkSinglePoint(p, &r))
	}
	for _, val := range res {
		fmt.Print(val)
		fmt.Print(" ")
	}
}
