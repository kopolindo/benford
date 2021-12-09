package charts

import (
	"fmt"
	"io"
	"os"
	"path"

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
	categories = []int{}
	for i := 1; i <= N; i++ {
		categories = append(categories, i)
	}
	scatter := charts.NewScatter()
	scatter.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Scatter plot for generated SSDs"}),
	)

	scatter.SetXAxis(categories).
		AddSeries("Series1", generateScatterItems(data))
	return scatter
}

type ScatterData struct{}

func (ScatterData) Create(name string, data []float64) {
	fname := path.Join("output", fmt.Sprintf("%s_scatter.html", name))
	page := components.NewPage()
	page.AddCharts(scatterGenerate(data))
	f, err := os.Create(fname)
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(f))
}
