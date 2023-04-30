package v2

import "fmt"

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

func (r rectangle) belongs(p point) bool {
	return (p.x >= r.lb.x && p.x <= r.ru.x) && (p.y >= r.lb.y && p.y <= r.ru.y)
}

type rectangles []rectangle
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

func checkSinglePoint(p point, rects *rectangles) int32 {
	var count int32
	for _, r := range *rects {
		if r.belongs(p) {
			count++
		}
	}
	return count
}

func Solve(r rectangles, ps points) {
	//r, ps, err := readData()
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	var res []int32
	for _, p := range ps {
		res = append(res, checkSinglePoint(p, &r))
	}
	//for _, val := range res {
	//	fmt.Print(val)
	//	fmt.Print(" ")
	//}
}
