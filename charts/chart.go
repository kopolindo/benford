package charts

import (
	"benford/structure"
	"benford/utilities"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"

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
type LineData struct{}

// CreateScatter creates a scatter plot
// Input: structure.Result
func (ScatterData) CreateScatter(r structure.Result) {
	outputPath := "output"
	utilities.MkdirINE(outputPath)
	name := strconv.Itoa(r.Sample)
	data := r.SSDs
	fname := path.Join(outputPath, fmt.Sprintf("%s_scatter.html", name))
	page := components.NewPage()
	page.AddCharts(scatterGenerate(data))
	f, err := os.Create(fname)
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(f))
}

// ---------------------- LINE PLOT --------------------------//
func generateLineItems(data []float64) []opts.LineData {
	items := make([]opts.LineData, 0)
	for i := 0; i < len(data); i++ {
		items = append(items, opts.LineData{Value: data[i]})
	}
	return items
}

func lineSmoothArea(plot structure.LinePlot) *charts.Line {
	plotName := plot.PlotName
	categories := plot.Categories
	lineSeries := plot.LineSeries
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: plotName}),
	)
	for _, serie := range lineSeries {
		// serie_n -> serie_n.Categories, serie_n.Values
		line.SetXAxis(categories).AddSeries(
			serie.Name,
			generateLineItems(serie.Values)).
			SetSeriesOptions(
				charts.WithLabelOpts(opts.Label{
					Show: true,
				}),
				charts.WithAreaStyleOpts(opts.AreaStyle{
					Opacity: 0.2,
				}),
				charts.WithLineChartOpts(opts.LineChart{
					Smooth: true,
				}),
			)
	}
	return line
}

// CreateLine creates a line plot
// Input: structure.LinePlot
func (LineData) CreateLine(plot structure.LinePlot) {
	outputPath := "output"
	utilities.MkdirINE(outputPath)
	page := components.NewPage()
	page.AddCharts(lineSmoothArea(plot))
	fname := path.Join("output", fmt.Sprintf("%s_line.html", plot.PlotName))
	f, err := os.Create(fname)
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(f))
}
