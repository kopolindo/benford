package charts

import (
	"io"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

var (
	categories []int
)

func generateScatterItems(data []float64) []opts.ScatterData {
	items := make([]opts.ScatterData, 0)
	for i := 0; i < len(data); i++ {
		items = append(items, opts.ScatterData{
			Value:      data[i],
			SymbolSize: 10,
		})
	}
	return items
}
func scatterGenerate(data []float64) *charts.Scatter {
	N := len(data)
	for i := 1; i <= N; i++ {
		categories = append(categories, i)
	}
	scatter := charts.NewScatter()
	scatter.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Scatter plot for generated SSDs"}),
	)

	scatter.YAxis.Scale = true
	scatter.SetXAxis(categories).
		AddSeries("Series1", generateScatterItems(data))
		//AddSeries("Category B", generateScatterItems()).

	return scatter
}

type ScatterData struct{}

func (ScatterData) Create(data []float64) {
	page := components.NewPage()
	page.AddCharts(scatterGenerate(data))
	f, err := os.Create("output/scatter.html")
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(f))
}
