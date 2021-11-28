package main

import (
	"os"

	"github.com/wcharczuk/go-chart"
)

type Point struct {
	Value float64
	Label int
}

type PointList []Point

var Values PointList

func Hist(Range map[int]float64) {
	graph := chart.BarChart{
		Title: "Benford",
		Background: chart.Style{
			Padding: chart.Box{
				Top: 40,
			},
		},
		Height:   512,
		BarWidth: 60,
		Bars: []chart.Value{
			{Value: Range[1], Label: "1"},
			{Value: Range[2], Label: "2"},
			{Value: Range[3], Label: "3"},
			{Value: Range[4], Label: "4"},
			{Value: Range[5], Label: "5"},
			{Value: Range[6], Label: "6"},
			{Value: Range[7], Label: "7"},
			{Value: Range[8], Label: "8"},
			{Value: Range[9], Label: "9"},
		},
	}

	f, _ := os.Create("output.png")
	defer f.Close()
	graph.Render(chart.PNG, f)
}
