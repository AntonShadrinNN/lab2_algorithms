package visualization

import (
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"os"
)

func getLineItems(data []int64) []opts.LineData {
	items := make([]opts.LineData, 0, len(data))
	for _, p := range data {
		//val := math.Log(float64(p))
		//if val == math.Inf(-1) {
		//	val = 0
		//}
		items = append(items, opts.LineData{
			Value: p,
		})
		//fmt.Println(val)
	}
	return items
}

func Draw(minDeg, maxDeg int, title string, filename string, coords ...[]int64) {

	lineItems := make([][]opts.LineData, len(coords))
	for _, c := range coords {
		lineItems = append(lineItems, getLineItems(c))
	}

	line := charts.NewLine()

	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			PageTitle: title,
			Theme:     types.ThemeInfographic,
		}),
		charts.WithTitleOpts(opts.Title{
			Title: title,
		}),
	)

	var xAxis []string
	for i := minDeg; i < maxDeg; i++ {
		//xAxis = append(xAxis, fmt.Sprintf("%f", math.Log(float64(i))))
		xAxis = append(xAxis, fmt.Sprintf("%d", i))
	}
	line.SetXAxis(xAxis).
		AddSeries("map", getLineItems(coords[0])).
		AddSeries("brute", getLineItems(coords[1])).
		AddSeries("tree", getLineItems(coords[2])).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))
	f, _ := os.Create(filename + ".html")
	_ = line.Render(f)
}
