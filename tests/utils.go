package tests

import (
	"lab2_algorithms/internal"
	"math/big"
)

func getRectangles(n int32) internal.Rectangles {

	var sample internal.Rectangles

	for i := int32(0); i < n; i++ {
		leftBottomEdge := internal.Point{
			X: 10 * i,
			Y: 10 * i,
		}
		rightUpperEdge := internal.Point{
			X: 10 * (2*n - i),
			Y: 10 * (2*n - i),
		}

		rect := internal.Rectangle{
			Lb: leftBottomEdge,
			Ru: rightUpperEdge,
		}

		sample = append(sample, rect)
	}

	return sample
}

func getTestPoints(n int32) internal.Points {
	var sample internal.Points
	//rand.Seed(time.Now().UnixNano())
	//largePrime := rand.Intn(6000) + 5000
	//
	//primesSieve := primes.Sieve(largePrime)
	//p1 := int64(primesSieve[len(primesSieve)-1])
	//p2 := int64(primesSieve[len(primesSieve)-2])
	p1 := 101
	p2 := 103
	for i := int32(0); i < n; i++ {
		xBig := big.NewInt(int64(p1))
		x := xBig.Exp(xBig, big.NewInt(31), big.NewInt(int64(20*n))).Int64()
		yBig := big.NewInt(int64(p2))
		y := yBig.Exp(yBig, big.NewInt(31), big.NewInt(int64(20*n))).Int64()

		p := internal.Point{
			X: int32(x),
			Y: int32(y),
		}
		sample = append(sample, p)
	}
	return sample
}

func getTestData(n int32) (internal.Rectangles, internal.Points) {
	return getRectangles(n), getTestPoints(n)
}
