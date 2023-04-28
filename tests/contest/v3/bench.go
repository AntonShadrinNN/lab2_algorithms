package v3

import (
	"fmt"
	"math"
	"math/big"
	"time"
)

func getTestPoints(n int32) points {
	var sample points
	//rand.Seed(time.Now().UnixNano())
	//largePrime := rand.Intn(6000) + 5000

	//primesSieve := primes.Sieve(largePrime)
	//p1 := int64(primesSieve[len(primesSieve)-1])
	//p2 := int64(primesSieve[len(primesSieve)-2])
	p1 := 101
	p2 := 103
	for i := int32(0); i < n; i++ {
		xBig := big.NewInt(int64(p1))
		x := xBig.Exp(xBig, big.NewInt(31), big.NewInt(int64(20*i))).Int64()
		yBig := big.NewInt(int64(p2))
		y := yBig.Exp(yBig, big.NewInt(31), big.NewInt(int64(20*i))).Int64()

		p := point{
			x: int32(x),
			y: int32(y),
		}
		//fmt.Println(p)
		sample = append(sample, p)
	}
	return sample
}

func getTestData(n int32) (rectangles, points) {
	return getRectangles(n), getTestPoints(n)
}

func getRectangles(n int32) rectangles {

	var sample rectangles

	for i := int32(0); i < n; i++ {
		leftBottomEdge := point{
			x: 10 * i,
			y: 10 * i,
		}
		rightUpperEdge := point{
			x: 10 * (2*n - i),
			y: 10 * (2*n - i),
		}

		rect := rectangle{
			lb: leftBottomEdge,
			ru: rightUpperEdge,
		}

		//fmt.Println(rect)

		sample = append(sample, rect)
	}

	return sample
}

func Bench() ([]int64, []int64, []int64) {

	testRects := make([]rectangles, 18)
	testPoints := make([]points, 18)
	for i := int32(1); i <= 18; i++ {
		r, p := getTestData(int32(math.Pow(2, float64(i))))
		testRects[i-1] = r
		testPoints[i-1] = p
	}

	prepTime := make([]int64, 0, 18)
	//wg := &sync.WaitGroup{}
	for i := range testRects {
		fmt.Printf("Testing 2^%d\n", i+1)
		start := time.Now()
		time.Sleep(10 * time.Nanosecond)
		compX, compY := testRects[i].compressCoordinates()
		_ = buildPersistentTree(&testRects[i], compX, compY)
		prepTime = append(prepTime, time.Since(start).Nanoseconds())
	}

	//os.Exit(0)
	executionTime := make([]int64, 0, 18)
	for j := range testRects {
		fmt.Printf("Testing 2^%d\n", j+1)
		compX, compY := testRects[j].compressCoordinates()
		tree := buildPersistentTree(&testRects[j], compX, compY)
		start := time.Now()
		time.Sleep(10 * time.Nanosecond)
		for _, p := range testPoints[j] {
			cX := binarySearch(compX, p.x)
			cY := binarySearch(compY, p.y)
			if cX == -1 || cY == -1 {
				continue
			}
			searchInTree(tree[cX], cY)
		}
		executionTime = append(executionTime, time.Since(start).Nanoseconds())
	}

	totalTime := make([]int64, 0, 18)
	for i := range testRects {
		fmt.Printf("Testing 2^%d\n", i+1)
		start := time.Now()
		time.Sleep(10 * time.Nanosecond)
		Solve(testRects[i], testPoints[i])
		//os.Exit(0)
		totalTime = append(totalTime, time.Since(start).Nanoseconds())
	}

	return totalTime, prepTime, executionTime
}
