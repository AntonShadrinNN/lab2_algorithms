package src

import (
	"lab2_algorithms/internal"
)

//rects, p, err := internal.ReadData()
//if err != nil {
//	fmt.Println(err)
//	return
//}

//for _, v := range res {
//	fmt.Print(v, " ")
//}
//fmt.Println()

func CompressedMap(rects internal.Rectangles, p internal.Points) {
	res := make([]int32, 0, len(p))
	compX, compY, rectMap := rects.BuildMap()

	for i := range p {
		pX := internal.BinarySearch(compX, p[i].X)
		pY := internal.BinarySearch(compY, p[i].Y)

		if pX == -1 || pY == -1 {
			res = append(res, 0)
			continue
		}
		res = append(res, rectMap[pY][pX])
	}
}
