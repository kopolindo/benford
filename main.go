package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path"
	"strconv"
	"sync"
	"time"

	charts "benford/charts"
	"benford/sampleGeneration"
	"benford/structure"
	"benford/utilities"

	"github.com/schollz/progressbar/v3"
)

const (
	DEFAULT_MIN = 100
	DEFAULT_MAX = 100000
)

var (
	sample           int
	minSample        int
	maxSample        int
	iterations       int
	Version          string
	BuildCommitShort string
	version          bool
	verbose          bool
	humanReadable    bool
	chart            bool
	csvOutput        string
	csvFile          *os.File
	outputFolder     string
	w                *csv.Writer
	iterationsWG     sync.WaitGroup
	sampleWG         sync.WaitGroup
	mainMutex        sync.Mutex
	mainWg           sync.WaitGroup
	linePlot         structure.LinePlot
	minSerie         structure.LineSerie
	averageSerie     structure.LineSerie
	maxSerie         structure.LineSerie
	samplesBar       *progressbar.ProgressBar
	iterBar          *progressbar.ProgressBar
)

func checkConditions() {
	// Print version
	if version {
		fmt.Println("Version:\t", Version)
		fmt.Println("Build:\t\t", BuildCommitShort)
		os.Exit(0)
	}
	// If no sample flag is provided
	if !utilities.IsFlagPassed("sample") &&
		!utilities.IsFlagPassed("min-sample") &&
		!utilities.IsFlagPassed("max-sample") {
		fmt.Println("You need to specify at least one sample ;)\n")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if utilities.IsFlagPassed("sample") &&
		(utilities.IsFlagPassed("min-sample") || utilities.IsFlagPassed("max-sample")) {
		fmt.Println("The following flags are incompatible one with each other:")
		fmt.Println("\tsample with min-sample and max-sample\n")
		fmt.Println("You use sample for a one-shot execution")
		fmt.Println("You use min-sample and max-sample to range samples\n")
		fmt.Println("eg: benford -iterations 1000 -min-sample 100 -max-sample 1000\n")
		fmt.Println("This would execute the program 900 times, using increasing samples")
		os.Exit(1)
	}
	// If minSample and maxSample are set, then sample is not needed
	if utilities.IsFlagPassed("sample") {
		minSample = sample
		maxSample = sample
	}
	// If csvOutput is matched, then create csv file and header
	if utilities.IsFlagPassed("csv") || utilities.IsFlagPassed("chart") {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println(err.Error())
		}
		outputFolder = path.Join(cwd, "output")
		utilities.MkdirINE(outputFolder)
		if utilities.IsFlagPassed("csv") {
			csvFile, err = os.Create(path.Join(outputFolder, csvOutput))
			if err != nil {
				fmt.Println(err.Error())
			}
			w = csv.NewWriter(csvFile)
		} else {
			w = csv.NewWriter(os.Stdout)
		}
		if e := w.Write([]string{
			"sample",
			"min",
			"max",
			"average",
			"devstd",
		}); e != nil {
			fmt.Println(e.Error())
		}
		w.Flush()
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
	flag.IntVar(&iterations, "iterations", 1, "Number of iterations")
	flag.IntVar(&sample, "sample", 0, "Size of the sample to be generated")
	flag.IntVar(&minSample, "min-sample", -1, "Start from this sample size")
	flag.IntVar(&maxSample, "max-sample", -1, "Finish with this sample size")
	flag.StringVar(&csvOutput, "csv", "", "CSV Output filename")
	flag.BoolVar(&version, "version", false, "Print version")
	flag.BoolVar(&verbose, "verbose", false, "Verbose, print compliancy")
	flag.BoolVar(&humanReadable, "human", false, "Human readable vs CSV readable")
	flag.BoolVar(&chart, "chart", false, "Create a scattered chart in output folder")
	flag.Parse()
	checkConditions()
}

// Worker function (for sample)
func worker(sample int, iter int, localWg *sync.WaitGroup, resChan chan structure.Result) {
	defer localWg.Done()
	//defer samplesBar.Add(1)
	var res structure.Result
	var ssdResults []float64
	var wg sync.WaitGroup
	wg.Add(iter)
	SSDsResultsChan := make(chan float64, 1)
	go func() {
		// Close channel
		wg.Wait()
		close(SSDsResultsChan)
	}()
	mainMutex.Lock()
	for it := 0; it < iter; it++ {
		// Create channel to make goroutine and main routine communicate
		go func(samp int, innerWg *sync.WaitGroup, result chan float64) {
			defer innerWg.Done()
			defer iterBar.Add(1)
			// Generate CVSS scores, normalize them (Exp) and take the first digit
			fdCVSSScores := sampleGeneration.GenerateFirstDigitCVSSScores(samp)
			// count ces of first left digits
			occurrences := utilities.CalcOccurrences(fdCVSSScores)
			// Communicate with main routine
			result <- utilities.SSD(occurrences, sample)
		}(sample, &wg, SSDsResultsChan)
	}
	// Fetch value from gorouting
	for r := range SSDsResultsChan {
		ssdResults = append(ssdResults, r)
	}
	mainMutex.Unlock()
	res.SSDs = ssdResults
	res.Sample = sample
	res.Average = utilities.Average(ssdResults)
	res.Max = utilities.Max(ssdResults)
	res.Min = utilities.Min(ssdResults)
	res.DevStd = utilities.DevStd(ssdResults)
	resChan <- res
}

func main() {
	resultChannel := make(chan structure.Result, 1)
	sampleSetSize := maxSample - minSample + 1
	totalIterations := sampleSetSize * iterations
	iterBar = progressbar.NewOptions(
		totalIterations,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetDescription("[green]Total Iterations[reset]"),
	)
	/*samplesBar = progressbar.NewOptions(
		sampleSetSize,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetDescription("[cyan]Samples[reset]"),
		progressbar.OptionShowCount(),
	)*/
	mainWg.Add(sampleSetSize)
	if verbose {
		fmt.Printf("minSample: %d\nmaxSample: %d\nsampleSetSize: %d\n",
			minSample, maxSample, sampleSetSize)
	}
	go func() {
		mainWg.Wait()
		close(resultChannel)
	}()
	for sample = minSample; sample <= maxSample; sample++ {
		go worker(sample, iterations, &mainWg, resultChannel)
	}
	// ----- OUTPUT ----- //
	for workerResult := range resultChannel {
		sample := workerResult.Sample
		min := workerResult.Min
		max := workerResult.Max
		average := workerResult.Average
		devstd := workerResult.DevStd
		if chart {
			var scatterChart charts.ScatterData
			scatterChart.CreateScatter(workerResult)
			// Create LinePlot
			linePlot.PlotName = "SSDs result distribution vs samples"
			linePlot.Categories = append(linePlot.Categories, strconv.Itoa(sample))
			//  Create minSerie
			minSerie.Name = "MIN"
			minSerie.Color = "blue"
			minSerie.Values = append(minSerie.Values, min)
			//  Create averageSerie
			averageSerie.Name = "AVERAGE"
			averageSerie.Color = "green"
			averageSerie.Values = append(averageSerie.Values, average)
			//  Create maxSerie
			maxSerie.Name = "MAX"
			maxSerie.Color = "red"
			maxSerie.Values = append(maxSerie.Values, max)
			linePlot.LineSeries = append(linePlot.LineSeries, maxSerie)
			linePlot.LineSeries = append(linePlot.LineSeries, averageSerie)
			linePlot.LineSeries = append(linePlot.LineSeries, minSerie)
			var lineChart charts.LineData
			lineChart.CreateLine(linePlot)
		}
		if humanReadable {
			fmt.Println("Min:", min)
			fmt.Println("Max:", max)
			fmt.Println("Average:", average)
			fmt.Println("DevStd", devstd)
		} else {
			if e := w.Write([]string{
				strconv.Itoa(sample),
				fmt.Sprintf("%.2f", min),
				fmt.Sprintf("%.2f", max),
				fmt.Sprintf("%.2f", average),
				fmt.Sprintf("%.2f", devstd),
			}); e != nil {
				fmt.Println(e.Error())
			}
			w.Flush()
		}
	}
	// If csvOutput is matched, then create csv file and header
	if utilities.IsFlagPassed("csvOutput") {
		csvFile.Close()
	}
}
