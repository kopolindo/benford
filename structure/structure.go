package structure

import "sync"

type Result struct {
	sync.RWMutex
	Sample  int
	Average float64
	Min     float64
	Max     float64
	DevStd  float64
	SSDs    []float64
}

type LineSerie struct {
	Name   string
	Values []float64
}

type LinePlot struct {
	PlotName   string
	Categories []string
	LineSeries []LineSerie
}
