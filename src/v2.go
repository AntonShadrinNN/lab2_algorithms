package src

import (
	"fmt"
	"lab2_algorithms/internal"
)

func CompressedMap(rects internal.Rectangles, p internal.Points) {
	rects, p, err := internal.ReadData()
	if err != nil {
		fmt.Println(err)
		return
	}
	var res []int32
	compX, compY, rectMap := rects.BuildMap()
	fmt.Println(rectMap)
	for i := range p {
		pX := internal.BinarySearch(compX, p[i].X)
		pY := internal.BinarySearch(compY, p[i].Y)

		if pX == -1 || pY == -1 {
			res = append(res, 0)
			continue
		}
		res = append(res, rectMap[pY][pX])
	}

	for _, v := range res {
		fmt.Print(v, " ")
	}
	fmt.Println()
}
