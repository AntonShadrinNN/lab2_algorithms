package src

import (
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

func BruteForce(r internal.Rectangles, ps internal.Points) {
	//r, ps, err := internal.ReadData()
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	res := make([]int32, 0, len(ps))
	for _, p := range ps {
		res = append(res, checkSinglePoint(p, &r))
	}
	//for _, val := range res {
	//	fmt.Print(val)
	//	fmt.Print(" ")
	//}
}
